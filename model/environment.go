package model

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

type Environment struct {
	Config    string    `json:"config" validate:"yaml"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	validator.RegisterValidation("yaml", service.validateYAML)
	return service
}

func (c *EnvironmentServiceImpl) GetEnvironment() (*Environment, error) {
	return c.Repository.Get()
}

func (c *EnvironmentServiceImpl) Update(newItem *Environment) (*Environment, error) {
	return c.Repository.Save(newItem)
}

func (e *Environment) decodeConfig() (interface{}, error) {
	item := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(e.Config), &item); err != nil {
		return nil, err
	}
	return item, nil
}

func (c *EnvironmentServiceImpl) validateYAML(fl validator.FieldLevel) bool {
	item := fl.Parent().Interface().(*Environment)
	_, err := item.decodeConfig()
	if err != nil {
		return false
	}
	return true
}
