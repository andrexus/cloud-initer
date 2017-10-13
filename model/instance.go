package model

import (
	"time"

	"fmt"

	"gopkg.in/mgo.v2/bson"
)

type Instance struct {
	ID         bson.ObjectId `json:"id"`
	Name       string        `json:"name"`
	IPAddress  string        `json:"ipAddress"`
	MACAddress string        `json:"macAddress"`
	MetaData   string        `json:"metaData"`
	UserData   string        `json:"userData"`
	CreatedAt  time.Time     `json:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt"`
}

type IPAddressValidationError struct {
	IPAddress string
}

func (e IPAddressValidationError) Error() string {
	return fmt.Sprintf("Instance with IP address %s already exists", e.IPAddress)
}

type MACAddressValidationError struct {
	MACAddress string
}

func (e MACAddressValidationError) Error() string {
	return fmt.Sprintf("Instance with MAC address %s already exists", e.MACAddress)
}

type FieldRequiredValidationError struct {
	FieldName string
}

func (e FieldRequiredValidationError) Error() string {
	return fmt.Sprintf("%s is required", e.FieldName)
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
	Repository InstanceRepository
}

func NewInstanceService(repository InstanceRepository) *InstanceServiceImpl {
	return &InstanceServiceImpl{
		Repository: repository,
	}
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
	if err := c.validatePayload(item); err != nil {
		return nil, err
	}

	item.ID = ""
	return c.Repository.Save(item)
}

func (c *InstanceServiceImpl) Update(item *Instance, newItem *Instance) (*Instance, error) {
	if err := c.validatePayload(newItem); err != nil {
		return nil, err
	}
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

func (c *InstanceServiceImpl) validatePayload(item *Instance) error {
	if item.IPAddress == "" {
		return &FieldRequiredValidationError{"ipAddress"}
	}
	if item.MACAddress == "" {
		return &FieldRequiredValidationError{"macAddress"}
	}
	existingItem, err := c.Repository.FindByIPAddress(item.IPAddress)
	if err != nil {
		return err
	}
	if existingItem != nil {
		return &IPAddressValidationError{item.IPAddress}
	}
	existingItem, err = c.Repository.FindByMACAddress(item.MACAddress)
	if err != nil {
		return err
	}
	if existingItem != nil {
		return &MACAddressValidationError{item.MACAddress}
	}
	return nil
}
