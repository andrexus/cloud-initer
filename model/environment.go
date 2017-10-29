package model

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
	"encoding/json"
)

type Environment struct {
	Config string    `json:"config" validate:"json"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type EnvironmentService interface {
	GetEnvironment() (*Environment, error)
	Update(newItem *Environment) (*Environment, error)
}

type EnvironmentServiceImpl struct {
	Repository EnvironmentRepository
}

func NewEnvironmentService(repository EnvironmentRepository, validator *validator.Validate) *EnvironmentServiceImpl {
	service := &EnvironmentServiceImpl{
		Repository: repository,
	}
	validator.RegisterValidation("json", service.validateJSON)
	return service
}

func (c *EnvironmentServiceImpl) GetEnvironment() (*Environment, error) {
	return c.Repository.Get()
}

func (c *EnvironmentServiceImpl) Update(newItem *Environment) (*Environment, error) {
	return c.Repository.Save(newItem)
}

func (c *EnvironmentServiceImpl) validateJSON(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	var item map[string]interface{}
	err := json.Unmarshal([]byte(value), &item)
	if err != nil {
		return false
	}
	return true
}
