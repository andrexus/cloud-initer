package model

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

type Instance struct {
	ID          bson.ObjectId `json:"id"`
	Name        string        `json:"name" validate:"required"`
	IPAddress   string        `json:"ipAddress" validate:"required,ip,uniqueIP"`
	MACAddress  string        `json:"macAddress" validate:"required,mac,uniqueMAC"`
	UserData    string        `json:"userData"`
	MetaData    string        `json:"metaData"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	RequestedAt time.Time     `json:"requestedAt"`
	RequestedBy string        `json:"requestedBy"`
}

type InstanceService interface {
	FindAll() ([]Instance, error)
	FindOne(id string) (*Instance, error)
	FindByIP(ipAddress string) (*Instance, error)
	Create(item *Instance) (*Instance, error)
	Update(item *Instance, newItem *Instance) (*Instance, error)
	Delete(id string) error
}

type InstanceServiceImpl struct {
	Repository         InstanceRepository
	EnvironmentService EnvironmentService
}

func NewInstanceService(repository InstanceRepository, validator *validator.Validate) *InstanceServiceImpl {
	service := &InstanceServiceImpl{
		Repository:         repository,
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
	item.UserData = newItem.UserData
	item.MetaData = newItem.MetaData
	return c.Repository.Save(item)
}

func (c *InstanceServiceImpl) Delete(id string) error {
	return c.Repository.Delete(id)
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
