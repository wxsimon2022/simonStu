package nacos

import (
	"fmt"
	"log"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/wxsimon2022/nacoswrap"

	"github.com/wxsimon2022/simonStu/config"
)

// client is the nacoswrap client shared across the app, initialized by InitNacos.
var client *nacoswrap.Client

// InitNacos 初始化 Nacos 客户端。
func InitNacos(cfg *config.Config) error {
	c, err := nacoswrap.NewClient(nacoswrap.Config{
		Host:      cfg.NacosHost,
		Port:      cfg.NacosPort,
		Namespace: cfg.NacosNamespace,
		Username:  cfg.NacosUsername,
		Password:  cfg.NacosPassword,
		LogDir:    cfg.LogDir + "/nacos",
		AppName:   "simonStu-go-service",
	})
	if err != nil {
		return fmt.Errorf("Nacos 客户端创建失败: %w", err)
	}

	client = c
	log.Printf("Nacos 客户端初始化成功 (%s:%d/%s)", cfg.NacosHost, cfg.NacosPort, cfg.NacosNamespace)
	if cfg.NacosUsername != "" {
		log.Printf("Nacos 已启用用户认证: %s", cfg.NacosUsername)
	}
	return nil
}

// GetServiceInstances 查询指定服务的可用实例列表。
func GetServiceInstances(serviceName string) ([]model.Instance, error) {
	if client == nil {
		return nil, fmt.Errorf("Nacos 客户端未初始化")
	}
	return client.GetInstances(serviceName)
}

// WatchService 监听服务实例变化。
func WatchService(serviceName string, onChange func(instances []model.Instance)) error {
	if client == nil {
		return fmt.Errorf("Nacos 客户端未初始化")
	}
	return client.Watch(serviceName, onChange)
}
