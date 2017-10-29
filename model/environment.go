package model

import (
	"time"

	"encoding/json"

	"gopkg.in/go-playground/validator.v9"
)

type Environment struct {
	Config    string    `json:"config" validate:"json"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type EnvironmentService interface {
	GetEnvironment() (*Environment, error)
	GetEnvironmentConfig() (map[string]interface{}, error)
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

func (c *EnvironmentServiceImpl) GetEnvironmentConfig() (map[string]interface{}, error) {
	env, err := c.Repository.Get()
	if err != nil {
		return nil, err
	}
	return decodeConfig(env.Config)
}

func (c *EnvironmentServiceImpl) Update(newItem *Environment) (*Environment, error) {
	return c.Repository.Save(newItem)
}

func decodeConfig(config string) (map[string]interface{}, error) {
	var item map[string]interface{}
	err := json.Unmarshal([]byte(config), &item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *EnvironmentServiceImpl) validateJSON(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := decodeConfig(value)
	if err != nil {
		return false
	}
	return true
}
