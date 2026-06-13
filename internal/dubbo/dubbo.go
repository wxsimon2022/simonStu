package dubbo

import (
	"context"
	"fmt"
	"log"
	"time"

	dubboCfg "dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/nacos-group/nacos-sdk-go/v2/model"

	"github.com/wxsimon2022/simonStu/internal/nacos"
)

// =========================================================================
// 对应 Java 端：org.apache.dubbo.simon.quickstart.dubbo.api.DemoService
//   public interface DemoService {
//       String sayHello(String name);
//   }
// =========================================================================

type DemoService struct {
	SayHello   func(ctx context.Context, name string) (string, error)
	SayGoodBye func(ctx context.Context, name string) (string, error)
}

var DemoSvc *DemoService

// NacosServiceName Java 端注册到 Nacos 的服务名。
const NacosServiceName = "providers:org.apache.dubbo.simon.quickstart.dubbo.api.DemoService::"

// DubboInit 初始化 Dubbo 消费者。
// 先用 Nacos 查 provider 地址，再用 dubbo-go Config API + direct URL 直连。
func DubboInit() error {
	log.Println("Dubbo 消费者初始化中...")

	// 1. 通过 Nacos 查询 provider 实例
	instances, err := nacos.GetServiceInstances(NacosServiceName)
	if err != nil {
		return fmt.Errorf("查询 Nacos 服务失败: %w", err)
	}
	if len(instances) == 0 {
		return fmt.Errorf("Nacos 上未找到可用的 DemoService 实例")
	}

	inst := instances[0]
	providerURL := fmt.Sprintf("tri://%s:%d", inst.Ip, inst.Port)
	log.Printf("Nacos 发现 provider: %s (healthy=%v)", providerURL, inst.Healthy)

	// 2. 用 Config API 注册服务（直接 URL，不配 registry）
	ref := &dubboCfg.ReferenceConfig{
		InterfaceName:  "org.apache.dubbo.simon.quickstart.dubbo.api.DemoService",
		Protocol:       "tri",
		URL:            providerURL,
		Retries:        "2",
		RequestTimeout: "30s",
		Serialization:  "hessian2",
	}

	// 3. 创建服务代理
	svc := &DemoService{}
	rootCfg := dubboCfg.NewRootConfigBuilder().Build()
	ref.Init(rootCfg)

	// 4. 建立连接、创建代理（必须在 Implement 之前）
	ref.Refer(svc)

	// 5. 注册服务代理
	ref.Implement(svc)

	// 等待连接就绪
	time.Sleep(1 * time.Second)

	// 6. 注册 Nacos 订阅，让数据能在 Nacos 控制台的"订阅者列表"中显示
	if err := nacos.WatchService(NacosServiceName, func(instances []model.Instance) {
		log.Printf("Nacos 服务实例变化，当前 %d 个实例", len(instances))
	}); err != nil {
		log.Printf("Nacos 订阅失败（非关键）: %v", err)
	}

	DemoSvc = svc
	log.Println("Dubbo 消费者初始化完成，已直连 provider")
	return nil
}
