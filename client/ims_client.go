package client

import (
	"os"
	"sync"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ims "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ims/v20201229"
)

var (
	imsClient     *ims.Client
	imsClientOnce sync.Once
)

// GetIMSClient 获取腾讯云图片内容安全客户端单例
func GetIMSClient() *ims.Client {
	imsClientOnce.Do(func() {
		// 从环境变量读取密钥（与 elysia-llm-tool 混元服务共用同一组密钥）
		credential := common.NewCredential(
			os.Getenv("TENCENTCLOUD_SECRET_ID"),
			os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		)

		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "ims.tencentcloudapi.com"

		client, err := ims.NewClient(credential, "ap-guangzhou", cpf)
		if err != nil {
			panic("初始化腾讯云 IMS 客户端失败: " + err.Error())
		}

		imsClient = client
	})

	return imsClient
}
