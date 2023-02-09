package redis

import (
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/redis/go-redis/v9"
	"sync"
)

type RedisService struct {
	container core.IContainer          // 服务容器
	clients   map[string]*redis.Client // key为uniqKey, value为redis.Client (连接池）
	lock      *sync.RWMutex            // 读写锁
}

func NewRedisService(params ...interface{}) (interface{}, error) {
	return &RedisService{
		container: params[0].(core.IContainer),
		clients:   make(map[string]*redis.Client),
		lock:      &sync.RWMutex{},
	}, nil
}

func (r *RedisService) GetClient(option ...contract.RedisOption) (*redis.Client, error) {
	// 读取默认配置
	config := GetDefaultConfig(r.container)

	// option对opt进行修改
	for _, opt := range option {
		if err := opt(r.container, config); err != nil {
			return nil, err
		}
	}

	// 如果最终的config没有设置dsn,就生成dsn
	key := config.UniqKey()

	// 判断是否已经实例化了redis.Client
	r.lock.RLock()
	if db, ok := r.clients[key]; ok {
		r.lock.RUnlock()
		return db, nil
	}
	r.lock.RUnlock()

	// 没有实例化gorm.DB，那么就要进行实例化操作
	r.lock.Lock()
	defer r.lock.Unlock()

	// 实例化gorm.DB
	client := redis.NewClient(config.Options)

	// 挂载到map中，结束配置
	r.clients[key] = client

	return client, nil
}
