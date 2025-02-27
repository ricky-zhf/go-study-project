package felo

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
	"log"
	"os"
)

const (
	packageName = "com.felolive.felo"        // 替换为你的应用包名
	orderId     = "GPA.3366-9359-3098-00247" // 替换为真实订单号
)

// 初始化 Google Play 服务客户端
func initAndroidPublisherService(credentialsJSON []byte) (*androidpublisher.Service, error) {
	// 从 JSON 密钥加载配置
	config, err := google.JWTConfigFromJSON(credentialsJSON, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, fmt.Errorf("加载密钥失败: %v", err)
	}

	// 创建 HTTP 客户端
	ctx := context.Background()
	client := config.Client(ctx)

	// 初始化 Android Publisher 服务
	service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("创建服务失败: %v", err)
	}

	return service, nil
}

// 发起退款请求
func refundOrder(service *androidpublisher.Service, packageName, orderId string, revoke bool) error {
	// 调用退款接口
	err := service.Orders.Refund(packageName, orderId).Revoke(revoke).Do()
	if err != nil {
		return fmt.Errorf("退款请求失败: %v", err)
	}
	return nil
}

func googleRefund() {
	// 1. 加载服务账号密钥文件（替换为实际路径）
	credentialsJSON, err := os.ReadFile("felo-live-google-credentials-file.json")
	if err != nil {
		log.Fatalf("读取密钥文件失败: %v", err)
	}

	// 2. 初始化服务客户端
	service, err := initAndroidPublisherService(credentialsJSON)
	if err != nil {
		log.Fatalf("初始化服务失败: %v", err)
	}

	// 4. 发起退款
	if err = refundOrder(service, packageName, orderId, true); err != nil {
		log.Fatalf("退款失败: %v", err)
	}

	fmt.Println("退款请求成功！")
}
