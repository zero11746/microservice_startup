# microservice_startup
初学者，写的go微服务启动代码，这是一个 Go 语言基于 gRPC 构建的微服务基础设施启动程序，核心聚焦于微服务运行所需的各类基础设施初始化与管理，包含 MySQL、Redis、MongoDB 多数据源配置，集成 Jaeger 分布式链路追踪，实现结构化日志输出并同步至 Elasticsearch，同时完成各类中间件、跨服务 gRPC 调用客户端的初始化，为微服务业务层提供、底层支撑。
