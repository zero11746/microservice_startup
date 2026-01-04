package web

import (
	"errors"
	"fmt"
	"net/http"

	"api/config"
	"api/router"

	"github.com/gin-gonic/gin"
)

func Start() (*http.Server, error) {
	engine := gin.New()
	gin.SetMode(gin.DebugMode)

	// 创建jaeger trace
	/*tp, tpErr := tracer.JaegerTraceProvider(
		config.GetConfig().Jaeger.Endpoints,
		config.GetConfig().Server.Name,
		config.GetConfig().Server.Env,
		config.GetConfig().Jaeger.IsOpenOnlySamplerError,
	)
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))*/

	// 设置中间件
	engine.Use(
		gin.Logger(),
		gin.Recovery(),
		//otelgin.Middleware(config.GetConfig().Server.Name),
	)
	// 加载路由
	router.InitRouter(engine)

	// 启动
	addr := fmt.Sprintf("%s:%d", config.GetConfig().Server.Host, config.GetConfig().Server.Port)

	ginServer := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	go func() {
		err := ginServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) { // 排除正常关闭的错误
			panic(fmt.Sprintf("Gin server start failed: %v", err))
		}
	}()

	fmt.Printf("web server is running on %s\n", addr)
	return ginServer, nil
}
