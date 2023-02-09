package core

import (
	"errors"
	"sync"
)

// 服务容器接口
type IContainer interface {
	Bind(provider IServiceProvider) error                          //绑定服务提供者
	IsBind(key string) bool                                        //判断是否绑定
	Make(key string) (interface{}, error)                          //根据关键字获取服务
	MustMake(key string) interface{}                               //根据关键字获取服务，如果未绑定服务提供者，会panic
	MakeNew(key string, params []interface{}) (interface{}, error) //根据关键字和传入的params参数实例化服务
}

// 服务容器具体实现
type Container struct {
	providers map[string]IServiceProvider //存储注册的服务提供者
	instances map[string]interface{}      //存储具体的实例
	lock      sync.RWMutex                //锁住容器变更
}

func NewContainer() *Container {
	return &Container{
		providers: map[string]IServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

func (c *Container) ProviderList() []string {
	ret := make([]string, 0)
	for _, provider := range c.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}

func (c *Container) Bind(provider IServiceProvider) error {
	c.lock.Lock()
	key := provider.Name()
	c.providers[key] = provider
	c.lock.Unlock()

	// 如果服务提供者，不延迟实例化
	if provider.IsDefer() == false {
		if err := provider.Boot(c); err != nil {
			return err
		}
		// 实例化方法
		params := provider.Params(c)
		method := provider.Register(c)
		instance, err := method(params...)
		if err != nil {
			return err
		}
		c.instances[key] = instance
	}
	return nil
}

func (c *Container) IsBind(key string) bool {
	return c.findServiceProvider(key) != nil
}

func (c *Container) findServiceProvider(key string) IServiceProvider {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if sp, ok := c.providers[key]; ok {
		return sp
	}
	return nil
}

func (c *Container) Make(key string) (interface{}, error) {
	return c.make(key, nil, false)
}

func (c *Container) MustMake(key string) interface{} {
	serv, err := c.make(key, nil, false)
	if err != nil {
		panic("container not contain key " + key)
	}
	return serv
}

func (c *Container) MakeNew(key string, params []interface{}) (interface{}, error) {
	return c.make(key, params, true)
}

func (c *Container) newInstance(sp IServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.Boot(c); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(c)
	}
	method := sp.Register(c)
	ins, err := method(params...)
	if err != nil {
		return nil, err
	}
	return ins, err
}

func (c *Container) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	sp := c.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		//强制创建一个新实例
		return c.newInstance(sp, params)
	}

	if ins, ok := c.instances[key]; ok {
		//容器中有，则直接返回
		return ins, nil
	}

	// 容器中还未实例化，则进行一次实例化
	inst, err := c.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	c.instances[key] = inst
	return inst, nil
}
