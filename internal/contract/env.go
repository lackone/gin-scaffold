package contract

const (
	EnvProd = "prod" //代表生产环境
	EnvTest = "test" //代表测试环境
	EnvDev  = "dev"  //代表开发环境

	EnvKey = "app:env"
)

type Env interface {
	AppEnv() string //获取当前环境

	IsExist(string) bool // IsExist判断一个环境变量是否有被设置

	Get(string) string // Get获取某个环境变量，如果没有设置，返回""

	All() map[string]string // All获取所有的环境变量，.env和运行环境变量融合后结果
}
