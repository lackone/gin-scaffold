package core

// NewInstance 定义了如何创建一个新实例，所有服务容器的创建服务
type NewInstance func(...interface{}) (interface{}, error)

// ServiceProvider 定义一个服务提供者需要实现的接口
type IServiceProvider interface {
	Register(IContainer) NewInstance //服务提供者实例化的方法
	Boot(IContainer) error           //调用实例化时会调用，做一些准备工作
	IsDefer() bool                   //是否在注册时就实例化服务，为true表示延迟实例化
	Params(IContainer) []interface{} //传递给实例化的参数
	Name() string                    //服务提供者的凭证
}
