package contract

const AppKey = "app:app"

type App interface {
	AppID() string //表示当前app的唯一ID

	Version() string //当前版本

	BaseFolder() string //项目基础目录

	ConfigFolder() string //配置文件目录

	LogFolder() string //日志文件目录

	ProviderFolder() string //服务提供者目录

	MiddlewareFolder() string //中间件目录

	CommandFolder() string //命令目录

	RuntimeFolder() string //运行时信息

	TestFolder() string //测试信息

	StorageFolder() string //存储目录

	DeployFolder() string //存放部署的时候创建的文件夹

	AppFolder() string //AppFolder定义业务代码所在的目录，用于监控文件变更使用

	DocsFolder() string //文档目录

	RoutesFolder() string //路由目录

	LoadAppConfig(configs map[string]string) //加载新的APP配置
}
