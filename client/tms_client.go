package client

import (
	"os"
	"sync"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tms/v20201229"
)

var (
	tmsClient     *tms.Client
	tmsClientOnce sync.Once
)

// GetTMSClient 获取腾讯云文本内容安全客户端单例
func GetTMSClient() *tms.Client {
	tmsClientOnce.Do(func() {
		// 从环境变量读取密钥（与 elysia-llm-tool 混元服务共用同一组密钥）
		credential := common.NewCredential(
			os.Getenv("TENCENTCLOUD_SECRET_ID"),
			os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		)

		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "tms.tencentcloudapi.com"

		client, err := tms.NewClient(credential, "ap-guangzhou", cpf)
		if err != nil {
			panic("初始化腾讯云 TMS 客户端失败: " + err.Error())
		}

		tmsClient = client
	})

	return tmsClient
}
