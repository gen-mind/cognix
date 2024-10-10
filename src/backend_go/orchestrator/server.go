package main

import (
	"cognix.ch/api/v2/core/messaging"
	"cognix.ch/api/v2/core/repository"
	"context"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
	"time"
)

// Server represents a server that handles various tasks related to connectors and documents.
type Server struct {
	renewInterval time.Duration
	connectorRepo repository.ConnectorRepository
	docRepo       repository.DocumentRepository
	messenger     messaging.Client
	scheduler     gocron.Scheduler
	streamCfg     *messaging.StreamConfig
	cfg           *Config
}

// NewServer creates a new instance of Server.
// It takes a pointer to Config, ConnectorRepository, DocumentRepository,
// messaging.Client, and messaging.Config as input parameters.
// It returns a pointer to Server and an error.
func NewServer(
	cfg *Config,
	connectorRepo repository.ConnectorRepository,
	docRepo repository.DocumentRepository,
	messenger messaging.Client,
	messagingCfg *messaging.Config) (*Server, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Server{connectorRepo: connectorRepo,
		docRepo:       docRepo,
		renewInterval: time.Duration(cfg.RenewInterval) * time.Second,
		cfg:           cfg,
		messenger:     messenger,
		streamCfg:     messagingCfg.Stream,
		scheduler:     s,
	}, nil
}

// run executes the main logic of the Server.
// It schedules a reload task, starts a listener, and loads connectors from the database.
// It runs concurrently with other goroutines.
//
// The method does not return any errors.
func (s *Server) run(ctx context.Context) error {
	zap.S().Infof("Schedule reload task")
	go s.schedule()
	zap.S().Infof("Start listener ...")
	go s.loadFromDatabase()
	return nil
}

// loadFromDatabase loads connectors from the database and triggers the execution for each connector.
// If the messenger is offline, no action is taken and the method returns nil.
// If an error occurs while getting active connectors from the database,
// the method logs the error and returns the error.
// For each connector, it creates a new trigger and executes the Do method.
// If an error occurs during the execution of the trigger, it is logged.
// The method returns nil if all operations are successful. If any error occurs,
// it is returned as the result of the method.
func (s *Server) loadFromDatabase() error {
	ctx := context.Background()
	if !s.messenger.IsOnline() {
		zap.S().Infof("Messenger is offline.")
		return nil
	}
	zap.S().Infof("Loading connectors from db")
	connectors, err := s.connectorRepo.GetActive(ctx)
	if err != nil {
		zap.S().Errorf("Load connectors failed: %v", err)
		return err
	}
	for _, connector := range connectors {
		if err = NewTrigger(s.messenger, s.connectorRepo, s.docRepo, connector, s.cfg.FileSizeLimit, s.cfg.OAuthURL).Do(ctx); err != nil {
			zap.S().Errorf("run connector %d failed: %v", connector.ID, err)
		}
	}
	return nil
}

// schedule sets up a job to reload data from the database at a specified interval and starts the scheduler.
// It returns an error if there was an issue setting up the job.
func (s *Server) schedule() error {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(s.renewInterval),
		gocron.NewTask(s.loadFromDatabase),
		gocron.WithName("reload from database"),
	)
	if err != nil {
		return err
	}
	s.scheduler.Start()
	return nil

}
