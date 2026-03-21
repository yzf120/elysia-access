package service_impl

import (
	"context"
	"fmt"

	pb "github.com/yzf120/elysia-access/proto/content_security"
	"github.com/yzf120/elysia-access/service"
)

// ContentSecurityServiceImpl 内容安全审核服务实现（出入参处理层）
type ContentSecurityServiceImpl struct {
	contentSecurityService *service.ContentSecurityService
}

// NewContentSecurityServiceImpl 创建内容安全审核服务实例
func NewContentSecurityServiceImpl() *ContentSecurityServiceImpl {
	return &ContentSecurityServiceImpl{
		contentSecurityService: service.NewContentSecurityService(),
	}
}

// TextModeration 文本内容安全审核
func (s *ContentSecurityServiceImpl) TextModeration(ctx context.Context, req *pb.TextModerationRequest) (*pb.TextModerationResponse, error) {
	// 参数校验
	if req.Content == "" {
		return &pb.TextModerationResponse{
			Code:    -1,
			Message: "文本内容不能为空",
		}, nil
	}

	// 调用 service 层
	result, err := s.contentSecurityService.TextModeration(ctx, req.Content, req.BizType, req.UserId, req.DataId)
	if err != nil {
		return &pb.TextModerationResponse{
			Code:    -1,
			Message: fmt.Sprintf("文本审核调用失败: %v", err),
		}, nil
	}

	// 构建响应（service 结果 → proto 响应）
	resp := &pb.TextModerationResponse{
		Code:       0,
		Message:    "success",
		Suggestion: result.Suggestion,
		Label:      result.Label,
		Score:      result.Score,
		RequestId:  result.RequestId,
	}

	// 转换详细结果
	for _, detail := range result.DetailResults {
		resp.DetailResults = append(resp.DetailResults, &pb.TextModerationDetail{
			Label:      detail.Label,
			Suggestion: detail.Suggestion,
			Keywords:   detail.Keywords,
			Score:      detail.Score,
		})
	}

	return resp, nil
}

// ImageModeration 图片内容安全审核
func (s *ContentSecurityServiceImpl) ImageModeration(ctx context.Context, req *pb.ImageModerationRequest) (*pb.ImageModerationResponse, error) {
	// 参数校验
	if req.FileUrl == "" && len(req.FileContent) == 0 {
		return &pb.ImageModerationResponse{
			Code:    -1,
			Message: "图片审核请求缺少图片内容（file_url 和 file_content 不能同时为空）",
		}, nil
	}

	// 调用 service 层
	result, err := s.contentSecurityService.ImageModeration(ctx, req.FileUrl, req.FileContent, req.BizType, req.UserId, req.DataId)
	if err != nil {
		return &pb.ImageModerationResponse{
			Code:    -1,
			Message: fmt.Sprintf("图片审核调用失败: %v", err),
		}, nil
	}

	// 构建响应（service 结果 → proto 响应）
	resp := &pb.ImageModerationResponse{
		Code:       0,
		Message:    "success",
		Suggestion: result.Suggestion,
		Label:      result.Label,
		Score:      result.Score,
		RequestId:  result.RequestId,
	}

	return resp, nil
}
