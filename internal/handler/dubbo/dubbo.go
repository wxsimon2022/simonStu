// Package dubbo Dubbo 服务调用 handler。HTTP 入口 → Nacos 发现 → Dubbo Triple 调用 Java。
package dubbo

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dubboSvc "github.com/wxsimon2022/simonStu/internal/dubbo"
	"github.com/wxsimon2022/simonStu/internal/logger"
	"github.com/wxsimon2022/simonStu/internal/nacos"
	"github.com/wxsimon2022/simonStu/internal/response"
)

// =========================================================================
// 请求链路：
//
//	HTTP Client → Gin Handler → dubbo.DemoSvc.SayHello() → Nacos 寻址 → Java Dubbo Server
// =========================================================================

// NacosServices 查询 Nacos 上已注册的服务实例（用于调试和监控）。
//
// GET /commonApi/dubbo/services?service=providers:org.apache.dubbo.simon.quickstart.dubbo.api.DemoService:1.0.0
func NacosServices(c *gin.Context) {
	serviceName := c.DefaultQuery("service",
		"providers:org.apache.dubbo.simon.quickstart.dubbo.api.DemoService::")

	instances, err := nacos.GetServiceInstances(serviceName)
	if err != nil {
		logger.Errorf(c, "NacosServices 查询失败 service=%s err=%v", serviceName, err)
		response.Error(c, http.StatusBadGateway, "查询 Nacos 服务失败: "+err.Error())
		return
	}

	type instanceInfo struct {
		IP       string            `json:"ip"`
		Port     uint64            `json:"port"`
		Metadata map[string]string `json:"metadata"`
		Healthy  bool              `json:"healthy"`
	}

	var list []instanceInfo
	for _, inst := range instances {
		meta := make(map[string]string)
		for k, v := range inst.Metadata {
			meta[k] = v
		}
		list = append(list, instanceInfo{
			IP:       inst.Ip,
			Port:     inst.Port,
			Metadata: meta,
			Healthy:  inst.Healthy,
		})
	}

	logger.Infof(c, "NacosServices 查询 service=%s instances=%d", serviceName, len(list))
	response.Success(c, gin.H{
		"service":   serviceName,
		"instances": list,
	})
}

// SayHello 调用 Java Dubbo 服务的 sayHello 方法。
//
// GET /commonApi/dubbo/hello?name=simon
//
// Java 端方法签名：
//
//	public interface DemoService {
//	    String sayHello(String name);
//	}
func SayHello(c *gin.Context) {
	name := c.DefaultQuery("name", "World")

	if dubboSvc.DemoSvc == nil {
		logger.Errorf(c, "SayHello Dubbo 服务未初始化")
		response.Error(c, http.StatusServiceUnavailable, "Dubbo 服务未初始化")
		return
	}

	// 通过 dubbo-go 代理调用 Java 端的 sayHello 方法
	msg, err := dubboSvc.DemoSvc.SayHello(c.Request.Context(), name)
	if err != nil {
		logger.Errorf(c, "SayHello Dubbo 调用失败 err=%v", err)
		response.Error(c, http.StatusBadGateway, "调用 Dubbo 服务失败: "+err.Error())
		return
	}

	logger.Infof(c, "SayHello 成功 name=%s msg=%s", name, msg)
	response.Success(c, gin.H{
		"message": msg,
	})
}

func SayGoodBye(c *gin.Context) {

	name := c.DefaultQuery("name", "World")

	if dubboSvc.DemoSvc == nil {
		logger.Errorf(c, "SayHello Dubbo 服务未初始化")
		response.Error(c, http.StatusServiceUnavailable, "Dubbo 服务未初始化")
		return
	}

	// 通过 dubbo-go 代理调用 Java 端的 sayHello 方法
	msg, err := dubboSvc.DemoSvc.SayGoodBye(c.Request.Context(), name)
	if err != nil {
		logger.Errorf(c, "SayHello Dubbo 调用失败 err=%v", err)
		response.Error(c, http.StatusBadGateway, "调用 Dubbo 服务失败: "+err.Error())
		return
	}

	logger.Infof(c, "SayGoodBye 成功 name=%s msg=%s", name, msg)
	response.Success(c, gin.H{
		"message": msg,
	})
}
