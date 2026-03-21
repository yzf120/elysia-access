package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/yzf120/elysia-access/config"
	pb "github.com/yzf120/elysia-access/proto/content_security"
	"github.com/yzf120/elysia-access/router"
	"github.com/yzf120/elysia-access/service_impl"
	"trpc.group/trpc-go/trpc-go"
	thttp "trpc.group/trpc-go/trpc-go/http"
)

func main() {
	log.Println("Elysia Access Service Starting...")

	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，使用系统环境变量")
	}

	// 初始化配置
	config.InitConfig()

	// 初始化服务实例（router 层使用）
	router.InitServices()

	// 创建 trpc 服务器
	s := trpc.NewServer()

	// 注册 HTTP 路由
	mux := http.NewServeMux()
	router.RegisterRoutes(mux)
	thttp.RegisterNoProtocolServiceMux(s.Service("trpc.elysia.access.http"), mux)

	// 注册 trpc RPC 服务（供 elysia-backend 等服务通过 trpc 协议调用）
	pb.RegisterContentSecurityServiceService(
		s.Service("trpc.elysia.access.content_security.ContentSecurityService"),
		service_impl.NewContentSecurityServiceImpl(),
	)

	log.Println("Elysia Access Service Started")
	log.Println("支持的接口:")
	log.Println("  - POST /api/content-security/text  (文本内容安全审核)")
	log.Println("  - POST /api/content-security/image  (图片内容安全审核)")

	if err := s.Serve(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}

	log.Println("Elysia Access Service Stopped")
}