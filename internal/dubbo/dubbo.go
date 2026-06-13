package dubbo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/wxsimon2022/dubboconn"

	"github.com/wxsimon2022/simonStu/config"
)

// DemoService 对应 Java 端 org.apache.dubbo.simon.quickstart.dubbo.api.DemoService。
type DemoService struct {
	SayHello   func(ctx context.Context, name string) (string, error)
	SayGoodBye func(ctx context.Context, name string) (string, error)
}

var DemoSvc *DemoService

// NacosServiceName Java 端注册到 Nacos 的服务名。
const NacosServiceName = "providers:org.apache.dubbo.simon.quickstart.dubbo.api.DemoService::"

// ready 在 DubboInit 完成后关闭（无论成功还是失败）。
var ready = make(chan struct{})

// Ready 返回一个 channel，DubboInit 返回后关闭。
func Ready() <-chan struct{} { return ready }

// WaitReady 等待 Dubbo 就绪，最多等 timeout。支持通过 ctx 提前取消。
func WaitReady(ctx context.Context, timeout time.Duration) error {
	select {
	case <-ready:
		if DemoSvc == nil {
			return errors.New("dubbo: 初始化失败（DemoSvc 为 nil）")
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("dubbo: %v 内未就绪", timeout)
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DubboInit 初始化 Dubbo 消费者。
// 通过 dubboconn 一站式完成 Nacos 服务发现 + dubbo-go 代理创建。
func DubboInit(cfg *config.Config) (err error) {
	defer func() {
		close(ready)
		if err != nil {
			log.Printf("Dubbo 初始化失败: %v", err)
		} else {
			log.Println("Dubbo 消费者初始化完成，已直连 provider")
		}
	}()

	log.Println("Dubbo 消费者初始化中...")

	// 不能直接传 &DemoSvc（DemoSvc 是 *DemoService，&DemoSvc 是 **DemoService），
	// dubbo-go 的 ref.Refer/Implement 需要 *struct，分配本地变量确保非 nil。
	svc := &DemoService{}
	conn, err := dubboconn.Connect(dubboconn.Config{
		NacosHost:      cfg.NacosHost,
		NacosPort:      cfg.NacosPort,
		NacosNamespace: cfg.NacosNamespace,
		NacosUsername:  cfg.NacosUsername,
		NacosPassword:  cfg.NacosPassword,
		ServiceName:    NacosServiceName,
		InterfaceName:  "org.apache.dubbo.simon.quickstart.dubbo.api.DemoService",
		NacosAppName:   cfg.NacosAppName,
	}, svc)
	if err != nil {
		return fmt.Errorf("dubboconn: %w", err)
	}
	DemoSvc = svc

	// 注册 Nacos 订阅，让数据能在 Nacos 控制台的"订阅者列表"中显示
	_ = conn.Watch(func(url string) {
		log.Printf("Nacos 服务实例变化，新地址: %s", url)
	})

	return nil
}
