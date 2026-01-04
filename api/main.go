package main

import (
	"context"
	"log"

	"api/config"
	"api/web"
	srv "common"
	"common/env"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in  header
// @name  Authorization
// @description  请输入 JWT Token，格式为：Bearer {你的Token值}
func main() {
	err := env.InitEnvConfig()
	if err != nil {
		return
	}
	err = config.InitConfig("config.toml")
	if err != nil {
		return
	}
	server, _ := web.Start()

	//开启pprof 默认的访问路径是/debug/pprof
	//pprof.Register(r)
	////测试代码
	//r.GET("/mem", func(c *gin.Context) {
	//	// 业务代码运行
	//	outCh := make(chan int)
	//	// 每秒起10个goroutine，goroutine会阻塞，不释放内存
	//	tick := time.Tick(time.Second / 10)
	//	i := 0
	//	for range tick {
	//		i++
	//		fmt.Println(i)
	//		alloc1(outCh) // 不停的有goruntine因为outCh堵塞，无法释放
	//	}
	//})
	srv.Run(func() {
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Gin server graceful shutdown failed: %v", err)
		}
	})
}
