package model

import (
	"errors"

	"github.com/aymerick/raymond"
	"time"
)

type CloudInitData struct {
	UserData string `json:"userData"`
	MetaData string `json:"metaData"`
}

type CloudInitService interface {
	PreviewCloudInitData(userDataTemplate, metaDataTemplate string) (*CloudInitData, error)
	GetCloudInitDataForClient(ipAddress, userAgent string) (*CloudInitData, error)
}

type CloudInitServiceImpl struct {
	InstanceService    InstanceService
	EnvironmentService EnvironmentService
}

func NewCloudInitService(instanceService InstanceService, environmentService EnvironmentService) *CloudInitServiceImpl {
	service := &CloudInitServiceImpl{
		InstanceService:    instanceService,
		EnvironmentService: environmentService,
	}
	return service
}

func (c *CloudInitServiceImpl) PreviewCloudInitData(userDataTemplate, metaDataTemplate string) (*CloudInitData, error) {
	return c.newCloudInitDataFromTemplate(userDataTemplate, metaDataTemplate)
}

func (c *CloudInitServiceImpl) GetCloudInitDataForClient(ipAddress, userAgent string) (*CloudInitData, error) {
	item, err := c.InstanceService.FindByIP(ipAddress)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("no instance")
	}
	cloudInitData, err := c.newCloudInitDataFromTemplate(item.UserData, item.MetaData)
	if err != nil {
		return nil, err
	}
	item.RequestedAt = time.Now()
	item.RequestedBy = userAgent
	c.InstanceService.Update(item, item)

	return cloudInitData, nil
}

func (c *CloudInitServiceImpl) newCloudInitDataFromTemplate(userDataTemplate, metaDataTemplate string) (*CloudInitData, error) {
	config, err := c.EnvironmentService.GetEnvironmentConfig()
	if err != nil {
		return nil, err
	}
	cloudInitData := new(CloudInitData)

	userData, err := renderTemplate(userDataTemplate, config)
	if err != nil {
		return nil, err
	}
	cloudInitData.UserData = userData

	metaData, err := renderTemplate(metaDataTemplate, config)
	if err != nil {
		return nil, err
	}
	cloudInitData.MetaData = metaData

	return cloudInitData, nil
}

func renderTemplate(template string, ctx interface{}) (string, error) {
	tpl, err := raymond.Parse(template)
	if err != nil {
		return "", err
	}
	return tpl.Exec(ctx)
}
