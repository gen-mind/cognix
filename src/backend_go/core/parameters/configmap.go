package parameters

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ConfigMapValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (v ConfigMapValue) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Key, validation.Required),
	)
}
