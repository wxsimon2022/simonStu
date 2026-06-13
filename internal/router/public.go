package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon2022/simonStu/internal/handler"
	dubboHandler "github.com/wxsimon2022/simonStu/internal/handler/dubbo"
	"github.com/wxsimon2022/simonStu/internal/handler/httpcall"
	"github.com/wxsimon2022/simonStu/internal/handler/redis"
	"github.com/wxsimon2022/simonStu/internal/handler/upload"
)

func RegisterPublicRoutes(api *gin.RouterGroup) {

	api.GET("/hello", handler.Hello)
	api.GET("/ping", handler.Ping)

	api.POST("/redis/set", handler.RedisSet)
	api.GET("/redis/get", handler.RedisGet)
	api.POST("/redis/stock/deduct", handler.StockDeduct)

	api.POST("/concurrent/process", handler.ConcurrentProcess)
	api.GET("/concurrent/fetch", handler.ConcurrentMultiFetch)

	api.Any("/http/call", httpcall.HttpCall)
	api.GET("/http/call/local", httpcall.HttpCallLocal)

	api.POST("/upload", upload.Upload)
	api.POST("/upload/multiple", upload.UploadMultiple)

	api.POST("/redis/mget", redis.RedisMultiGet)
	api.POST("/redis/mget/serial", redis.RedisMultiGetOneByOne)
	api.POST("/redis/mset", redis.RedisMultiSet)

	// ———— Dubbo 服务调用（通过 Nacos 发现 Java 服务）————
	api.GET("/dubbo/services", dubboHandler.NacosServices)
	api.GET("/dubbo/hello", dubboHandler.SayHello)
	api.GET("/dubbo/sayGoodBye", dubboHandler.SayGoodBye)

	api.GET("/test", handler.Test)
}
