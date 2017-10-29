package model

import (
	"time"

	"github.com/aymerick/raymond"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

type Instance struct {
	ID          bson.ObjectId `json:"id"`
	Name        string        `json:"name" validate:"required"`
	IPAddress   string        `json:"ipAddress" validate:"required,ip,uniqueIP"`
	MACAddress  string        `json:"macAddress" validate:"required,mac,uniqueMAC"`
	MetaData    string        `json:"metaData"`
	UserData    string        `json:"userData"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	RequestedAt time.Time     `json:"requestedAt"`
	RequestedBy string        `json:"requestedBy"`
}

type Preview struct {
	MetaData string `json:"metaData"`
	UserData string `json:"userData"`
}

type InstanceService interface {
	FindAll() ([]Instance, error)
	FindOne(id string) (*Instance, error)
	FindByIP(ipAddress string) (*Instance, error)
	Create(item *Instance) (*Instance, error)
	Update(item *Instance, newItem *Instance) (*Instance, error)
	Delete(id string) error

	// Returns rendered template for meta-data and user-data
	Preview(id string) (*Preview, error)
}

type InstanceServiceImpl struct {
	Repository         InstanceRepository
	EnvironmentService EnvironmentService
}

func NewInstanceService(repository InstanceRepository, environmentService EnvironmentService, validator *validator.Validate) *InstanceServiceImpl {
	service := &InstanceServiceImpl{
		Repository:         repository,
		EnvironmentService: environmentService,
	}
	validator.RegisterValidation("uniqueIP", service.validateUniqueIP)
	validator.RegisterValidation("uniqueMAC", service.validateUniqueMAC)
	return service
}

func (c *InstanceServiceImpl) FindAll() ([]Instance, error) {
	return c.Repository.FindAll()
}

func (c *InstanceServiceImpl) FindOne(id string) (*Instance, error) {
	return c.Repository.FindOne(id)
}

func (c *InstanceServiceImpl) FindByIP(ipAddress string) (*Instance, error) {
	return c.Repository.FindByIPAddress(ipAddress)
}

func (c *InstanceServiceImpl) Create(item *Instance) (*Instance, error) {
	item.ID = ""
	return c.Repository.Save(item)
}

func (c *InstanceServiceImpl) Update(item *Instance, newItem *Instance) (*Instance, error) {
	item.Name = newItem.Name
	item.IPAddress = newItem.IPAddress
	item.MACAddress = newItem.MACAddress
	item.MetaData = newItem.MetaData
	item.UserData = newItem.UserData
	return c.Repository.Save(item)
}

func (c *InstanceServiceImpl) Delete(id string) error {
	return c.Repository.Delete(id)
}

func (c *InstanceServiceImpl) Preview(id string) (*Preview, error) {
	item, err := c.Repository.FindOne(id)
	if err != nil {
		return nil, err
	}
	config, err := c.EnvironmentService.GetEnvironmentConfig()
	if err != nil {
		return nil, err
	}

	preview := new(Preview)

	metaData, err := renderTemplate(item.MetaData, config)
	if err != nil {
		return nil, err
	}
	preview.MetaData = metaData
	userData, err := renderTemplate(item.UserData, config)
	if err != nil {
		return nil, err
	}
	preview.UserData = userData

	return preview, nil
}

func (c *InstanceServiceImpl) validateUniqueIP(fl validator.FieldLevel) bool {
	item := fl.Parent().Interface().(*Instance)
	existingItem, err := c.Repository.FindByIPAddress(item.IPAddress)
	if err != nil {
		return false
	}
	if existingItem != nil && existingItem.ID != item.ID {
		return false
	}
	return true
}

func (c *InstanceServiceImpl) validateUniqueMAC(fl validator.FieldLevel) bool {
	item := fl.Parent().Interface().(*Instance)
	existingItem, err := c.Repository.FindByMACAddress(item.MACAddress)
	if err != nil {
		return false
	}
	if existingItem != nil && existingItem.ID != item.ID {
		return false
	}
	return true
}

func renderTemplate(template string, ctx interface{}) (string, error) {
	tpl, err := raymond.Parse(template)
	if err != nil {
		return "", err
	}
	return tpl.Exec(ctx)
}
