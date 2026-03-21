package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ims "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ims/v20201229"
	tms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tms/v20201229"
	"github.com/yzf120/elysia-access/client"
)

// ContentSecurityService 内容安全审核服务（核心业务逻辑层）
type ContentSecurityService struct{}

// NewContentSecurityService 创建内容安全审核服务实例
func NewContentSecurityService() *ContentSecurityService {
	return &ContentSecurityService{}
}

// TextModerationResult 文本审核结果
type TextModerationResult struct {
	Suggestion    string                  // 建议操作：Pass-通过 Block-违规 Review-疑似
	Label         string                  // 恶意标签
	Score         int32                   // 恶意置信度 0-100
	DetailResults []*TextModerationDetail // 详细结果
	RequestId     string                  // 腾讯云请求ID
}

// TextModerationDetail 文本审核详细结果
type TextModerationDetail struct {
	Label      string
	Suggestion string
	Keywords   []string
	Score      int32
}

// ImageModerationResult 图片审核结果
type ImageModerationResult struct {
	Suggestion string // 建议操作：Pass-通过 Block-违规 Review-疑似
	Label      string // 恶意标签
	Score      int32  // 恶意置信度 0-100
	RequestId  string // 腾讯云请求ID
}

// TextModeration 文本内容安全审核
// 调用腾讯云 TMS（文本内容安全）接口进行文本审核
func (s *ContentSecurityService) TextModeration(ctx context.Context, content, bizType, userId, dataId string) (*TextModerationResult, error) {
	log.Printf("[ContentSecurity] 文本审核请求，userId: %s, dataId: %s", userId, dataId)

	tmsClient := client.GetTMSClient()

	// 构建腾讯云 TMS 请求
	tmsReq := tms.NewTextModerationRequest()

	// 文本内容需要 base64 编码
	// 检查内容是否已经是 base64 编码，如果不是则进行编码
	if _, err := base64.StdEncoding.DecodeString(content); err != nil {
		content = base64.StdEncoding.EncodeToString([]byte(content))
	}
	tmsReq.Content = common.StringPtr(content)

	// 设置业务类型
	if bizType == "" {
		bizType = "default"
	}
	tmsReq.BizType = common.StringPtr(bizType)

	// 设置用户ID和数据ID（用于审核日志追踪）
	if userId != "" {
		tmsReq.User = &tms.User{
			UserId: common.StringPtr(userId),
		}
	}
	if dataId != "" {
		tmsReq.DataId = common.StringPtr(dataId)
	}

	// 调用腾讯云 TMS 接口
	tmsResp, err := tmsClient.TextModeration(tmsReq)
	if err != nil {
		log.Printf("[ContentSecurity] 文本审核调用失败，userId: %s, err: %v", userId, err)
		return nil, fmt.Errorf("文本审核调用失败: %w", err)
	}

	// 构建结果
	result := &TextModerationResult{}

	if tmsResp.Response != nil {
		if tmsResp.Response.Suggestion != nil {
			result.Suggestion = *tmsResp.Response.Suggestion
		}
		if tmsResp.Response.Label != nil {
			result.Label = *tmsResp.Response.Label
		}
		if tmsResp.Response.Score != nil {
			result.Score = int32(*tmsResp.Response.Score)
		}
		if tmsResp.Response.RequestId != nil {
			result.RequestId = *tmsResp.Response.RequestId
		}

		// 解析详细结果
		if tmsResp.Response.DetailResults != nil {
			for _, detail := range tmsResp.Response.DetailResults {
				d := &TextModerationDetail{}
				if detail.Label != nil {
					d.Label = *detail.Label
				}
				if detail.Suggestion != nil {
					d.Suggestion = *detail.Suggestion
				}
				if detail.Score != nil {
					d.Score = int32(*detail.Score)
				}
				if detail.Keywords != nil {
					for _, kw := range detail.Keywords {
						if kw != nil {
							d.Keywords = append(d.Keywords, *kw)
						}
					}
				}
				result.DetailResults = append(result.DetailResults, d)
			}
		}
	}

	log.Printf("[ContentSecurity] 文本审核完成，userId: %s, suggestion: %s, label: %s, score: %d, requestId: %s",
		userId, result.Suggestion, result.Label, result.Score, result.RequestId)

	return result, nil
}

// ImageModeration 图片内容安全审核
// 调用腾讯云 IMS（图片内容安全）接口进行图片审核
// fileUrl 和 fileContent 二选一，fileUrl 为图片URL，fileContent 为图片二进制内容
func (s *ContentSecurityService) ImageModeration(ctx context.Context, fileUrl string, fileContent []byte, bizType, userId, dataId string) (*ImageModerationResult, error) {
	log.Printf("[ContentSecurity] 图片审核请求，userId: %s, dataId: %s, hasUrl: %v, hasContent: %v",
		userId, dataId, fileUrl != "", len(fileContent) > 0)

	imsClient := client.GetIMSClient()

	// 构建腾讯云 IMS 请求
	imsReq := ims.NewImageModerationRequest()

	// 设置图片内容（URL 或 base64 二选一）
	if fileUrl != "" {
		imsReq.FileUrl = common.StringPtr(fileUrl)
	} else if len(fileContent) > 0 {
		b64Content := base64.StdEncoding.EncodeToString(fileContent)
		imsReq.FileContent = common.StringPtr(b64Content)
	} else {
		return nil, fmt.Errorf("图片审核请求缺少图片内容（file_url 和 file_content 不能同时为空）")
	}

	// 设置业务类型
	if bizType == "" {
		bizType = "default"
	}
	imsReq.BizType = common.StringPtr(bizType)

	// 设置用户ID和数据ID
	if userId != "" {
		imsReq.User = &ims.User{
			UserId: common.StringPtr(userId),
		}
	}
	if dataId != "" {
		imsReq.DataId = common.StringPtr(dataId)
	}

	// 调用腾讯云 IMS 接口
	imsResp, err := imsClient.ImageModeration(imsReq)
	if err != nil {
		log.Printf("[ContentSecurity] 图片审核调用失败，userId: %s, err: %v", userId, err)
		return nil, fmt.Errorf("图片审核调用失败: %w", err)
	}

	// 构建结果
	result := &ImageModerationResult{}

	if imsResp.Response != nil {
		if imsResp.Response.Suggestion != nil {
			result.Suggestion = *imsResp.Response.Suggestion
		}
		if imsResp.Response.Label != nil {
			result.Label = *imsResp.Response.Label
		}
		if imsResp.Response.Score != nil {
			result.Score = int32(*imsResp.Response.Score)
		}
		if imsResp.Response.RequestId != nil {
			result.RequestId = *imsResp.Response.RequestId
		}
	}

	log.Printf("[ContentSecurity] 图片审核完成，userId: %s, suggestion: %s, label: %s, score: %d, requestId: %s",
		userId, result.Suggestion, result.Label, result.Score, result.RequestId)

	return result, nil
}
