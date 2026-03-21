package router

import (
	"encoding/json"
	"net/http"

	pb "github.com/yzf120/elysia-access/proto/content_security"
	"github.com/yzf120/elysia-access/service_impl"
)

var (
	contentSecurityServiceImpl *service_impl.ContentSecurityServiceImpl
)

// InitServices 初始化路由层所需的服务实例
func InitServices() {
	contentSecurityServiceImpl = service_impl.NewContentSecurityServiceImpl()
}

// RegisterRoutes 注册所有 HTTP 路由
func RegisterRoutes(mux *http.ServeMux) {
	// 内容安全审核接口
	mux.HandleFunc("/api/content-security/text", handleTextModeration)
	mux.HandleFunc("/api/content-security/image", handleImageModeration)

	// 健康检查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}

// handleTextModeration 处理文本审核 HTTP 请求
func handleTextModeration(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req pb.TextModerationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := contentSecurityServiceImpl.TextModeration(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    500,
			"message": "服务内部错误: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// handleImageModeration 处理图片审核 HTTP 请求
func handleImageModeration(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req pb.ImageModerationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resp, err := contentSecurityServiceImpl.ImageModeration(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    500,
			"message": "服务内部错误: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(resp)
}
