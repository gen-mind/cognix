package logic

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/parameters"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"time"
)

type (

	// PersonaBL represents the business logic layer for personas.
	PersonaBL interface {
		GetAll(ctx context.Context, user *model.User, archived bool) ([]*model.Persona, error)
		GetByID(ctx context.Context, user *model.User, id int64) (*model.Persona, error)
		Create(ctx context.Context, user *model.User, param *parameters.PersonaParam) (*model.Persona, error)
		Update(ctx context.Context, id int64, user *model.User, param *parameters.PersonaParam) (*model.Persona, error)
		Archive(ctx context.Context, user *model.User, id int64, restore bool) (*model.Persona, error)
	}

	// personaBL represents the business logic layer for personas.
	personaBL struct {
		personaRepo repository.PersonaRepository
		chatRepo    repository.ChatRepository
	}
)

// Archive is a method of the personaBL struct that archives or restores a persona.
// It takes a context, user, persona ID, and a restore flag as input parameters.
// If the user does not have the SuperAdmin or Admin roles, it returns an error indicating insufficient permissions.
// It retrieves the persona from the persona repository by ID and tenant ID, and masks the persona's LLM.ApiKey if it is not nil.
// If an error occurs while retrieving the persona, it returns the error.
// If the restore flag is true, it sets the persona's DeletedDate to pg.NullTime{}, otherwise it sets it to the current UTC time.
// It also sets the persona's LastUpdate to the current UTC time.
// If an error occurs while archiving the persona, it returns the error.
// Finally, it returns the archived or restored persona and a nil error.
func (b *personaBL) Archive(ctx context.Context, user *model.User, id int64, restore bool) (*model.Persona, error) {
	if !user.HasRoles(model.RoleSuperAdmin, model.RoleAdmin) {
		return nil, utils.ErrorPermission.New("do not have permission")
	}

	persona, err := b.personaRepo.GetByID(ctx, id, user.TenantID)
	if persona.LLM != nil {
		persona.LLM.ApiKey = persona.LLM.MaskApiKey()
	}

	if err != nil {
		return nil, err
	}
	if restore {
		persona.DeletedDate = pg.NullTime{}
	} else {
		persona.DeletedDate = pg.NullTime{time.Now().UTC()}
	}
	persona.LastUpdate = pg.NullTime{time.Now().UTC()}

	if err = b.personaRepo.Archive(ctx, persona); err != nil {
		return nil, err
	}
	return persona, nil
}

// Create is a method of the personaBL struct that creates a new persona.
// It takes a context, user, and persona parameter as input parameters.
// It marshals the starter messages from the parameter into JSON.
// It creates a new persona with the name, description, tenant ID, visibility, starter messages, creation date, LLM, and prompt values from the parameter.
// It sets the LLM name by combining the user's first name and the parameter model ID.
// It sets the LLM creation date and URL from the parameter.
// If the persona's LLM API key is different from the parameter API key, it sets the LLM API key to the parameter API key.
// It sets the LLM last update time.
// It sets the prompt name, description, system prompt, task prompt, and creation date from the parameter.
// If an error occurs while creating the persona, it returns the error.
// Finally, it returns the created persona and a nil error.
func (b *personaBL) Create(ctx context.Context, user *model.User, param *parameters.PersonaParam) (*model.Persona, error) {

	starterMessages, err := json.Marshal(param.StarterMessages)
	if err != nil {
		return nil, utils.ErrorBadRequest.Wrap(err, "fail to marshal starter messages")
	}
	persona := model.Persona{
		Name:            param.Name,
		DefaultPersona:  true,
		Description:     param.Description,
		TenantID:        user.TenantID,
		IsVisible:       true,
		StarterMessages: starterMessages,
		CreationDate:    time.Now().UTC(),
		LLM: &model.LLM{
			Name:         fmt.Sprintf("%s %s", user.FirstName, param.ModelID),
			ModelID:      param.ModelID,
			TenantID:     user.TenantID,
			CreationDate: time.Now().UTC(),
			Url:          param.URL,
			ApiKey:       param.APIKey,
			Endpoint:     param.Endpoint,
		},
		Prompt: &model.Prompt{
			UserID:       user.ID,
			Name:         param.Name,
			Description:  param.Description,
			SystemPrompt: param.SystemPrompt,
			TaskPrompt:   param.TaskPrompt,
			CreationDate: time.Now().UTC(),
		},
	}
	if err := b.personaRepo.Create(ctx, &persona); err != nil {
		return nil, err
	}
	return &persona, nil
}

// Update is a method of the personaBL struct that updates a persona based on the given ID, user, and parameter.
// It takes a context, ID, user, and parameter as input parameters.
// It retrieves the persona from the persona repository by ID and tenant ID.
// If an error occurs while retrieving the persona, it returns the error.
// It marshals the starter messages from the parameter into JSON.
// It updates the persona's name, description, last update time, and starter messages with the values from the parameter.
// It updates the persona's LLM endpoint and model ID with the values from the parameter.
// If the persona's LLM.ApiKey is different from the parameter API key, it updates the persona's LLM API key with the parameter API key.
// It also updates the persona's LLM last update time.
// It updates the persona's prompt name, description, system prompt,
func (b *personaBL) Update(ctx context.Context, id int64, user *model.User, param *parameters.PersonaParam) (*model.Persona, error) {
	persona, err := b.personaRepo.GetByID(ctx, id, user.TenantID)
	if err != nil {
		return nil, err
	}
	starterMessages, err := json.Marshal(param.StarterMessages)
	if err != nil {
		return nil, utils.ErrorBadRequest.Wrap(err, "fail to marshal starter messages")
	}
	persona.Name = param.Name
	persona.Description = param.Description
	persona.LastUpdate = pg.NullTime{time.Now().UTC()}
	persona.StarterMessages = starterMessages
	persona.LLM.Endpoint = param.Endpoint
	persona.LLM.ModelID = param.ModelID
	// update api key if user updates it.
	if persona.LLM.MaskApiKey() != param.APIKey {
		persona.LLM.ApiKey = param.APIKey
	}
	persona.LLM.LastUpdate = pg.NullTime{time.Now().UTC()}
	persona.Prompt.Name = param.Name
	persona.Prompt.Description = param.Description
	persona.Prompt.SystemPrompt = param.SystemPrompt
	persona.Prompt.TaskPrompt = param.TaskPrompt
	persona.Prompt.LastUpdate = pg.NullTime{time.Now().UTC()}

	if err = b.personaRepo.Update(ctx, persona); err != nil {
		return nil, err
	}
	return persona, nil
}

// NewPersonaBL is a function that creates a new instance of the PersonaBL interface.
// It takes a PersonaRepository and a ChatRepository as input parameters and returns a PersonaBL.
// The PersonaBL implementation uses the provided PersonaRepository and ChatRepository to interact with the data layer.
// The PersonaBL interface provides methods for managing personas, such as retrieving all personas, creating a new persona, updating a persona, retrieving a persona by ID, and archiving or restoring a persona.
func NewPersonaBL(personaRepo repository.PersonaRepository,
	chatRepo repository.ChatRepository) PersonaBL {
	return &personaBL{
		personaRepo: personaRepo,
		chatRepo:    chatRepo,
	}
}

// GetAll is a method of the personaBL struct that retrieves all personas based on the given archived flag.
// It takes a context, user, and archived flag as input parameters.
// It delegates the retrieval operation to the persona repository by calling the GetAll method with the given context, user's tenant ID, and archived flag.
// Finally, it returns the retrieved personas and a nil error.
func (b *personaBL) GetAll(ctx context.Context, user *model.User, archived bool) ([]*model.Persona, error) {
	return b.personaRepo.GetAll(ctx, user.TenantID, archived)
}

// GetByID is a method of the personaBL struct that retrieves a persona by ID.
// It takes a context, user, and persona ID as input parameters.
// It retrieves the persona from the persona repository by ID and tenant ID.
// If an error occurs while retrieving the persona, it returns the error.
// If the persona's LLM is not nil, it masks the persona's LLM.ApiKey.
// Finally, it returns the retrieved persona and a nil error.
func (b *personaBL) GetByID(ctx context.Context, user *model.User, id int64) (*model.Persona, error) {
	persona, err := b.personaRepo.GetByID(ctx, id, user.TenantID)
	if err != nil {
		return nil, err
	}
	if persona.LLM != nil {
		persona.LLM.ApiKey = persona.LLM.MaskApiKey()
	}
	return persona, nil
}
