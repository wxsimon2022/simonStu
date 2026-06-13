package nacos

import (
	"fmt"
	"log"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/wxsimon8888/simonStu/config"
)

// NamingClient 服务发现客户端，用于查询 Nacos 上的服务实例。
var NamingClient naming_client.INamingClient

// InitNacos 初始化 Nacos 客户端。
func InitNacos(cfg *config.Config) error {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      cfg.NacosHost,
			Port:        cfg.NacosPort,
			ContextPath: "/nacos",
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         cfg.NacosNamespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              cfg.LogDir + "/nacos",
		CacheDir:            cfg.LogDir + "/nacos/cache",
		LogLevel:            "info",
		Username:            cfg.NacosUsername,
		Password:            cfg.NacosPassword,
		AppName:             "simonStu-go-service",
	}

	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return fmt.Errorf("Nacos 客户端创建失败: %w", err)
	}

	NamingClient = client
	log.Printf("Nacos 客户端初始化成功 (%s:%d/%s)", cfg.NacosHost, cfg.NacosPort, cfg.NacosNamespace)
	if cfg.NacosUsername != "" {
		log.Printf("Nacos 已启用用户认证: %s", cfg.NacosUsername)
	}
	return nil
}

// GetServiceInstances 查询指定服务的可用实例列表。
func GetServiceInstances(serviceName string) ([]model.Instance, error) {
	if NamingClient == nil {
		return nil, fmt.Errorf("Nacos 客户端未初始化")
	}

	instances, err := NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		HealthyOnly: true,
	})
	if err != nil {
		return nil, fmt.Errorf("查询服务 %s 失败: %w", serviceName, err)
	}
	return instances, nil
}

// WatchService 监听服务实例变化。
func WatchService(serviceName string, onChange func(instances []model.Instance)) error {
	if NamingClient == nil {
		return fmt.Errorf("Nacos 客户端未初始化")
	}

	return NamingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: serviceName,
		SubscribeCallback: func(instances []model.Instance, err error) {
			if err != nil {
				log.Printf("Nacos 服务 %s 变更通知出错: %v", serviceName, err)
				return
			}
			log.Printf("Nacos 服务 %s 实例变化，当前 %d 个实例", serviceName, len(instances))
			onChange(instances)
		},
	})
}
