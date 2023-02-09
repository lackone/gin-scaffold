package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"strings"
	"sync"
	"time"
)

type MemoryData struct {
	val        interface{}
	createTime time.Time
	ttl        time.Duration
}

type MemoryService struct {
	container core.IContainer
	data      map[string]*MemoryData
	lock      sync.RWMutex
}

func NewMemoryService(params ...interface{}) (interface{}, error) {
	container := params[0].(core.IContainer)

	return &MemoryService{
		container: container,
		data:      map[string]*MemoryData{},
		lock:      sync.RWMutex{},
	}, nil
}

func (m *MemoryService) Get(ctx context.Context, key string) (string, error) {
	var val string
	if err := m.GetObj(ctx, key, &val); err != nil {
		return "", err
	}
	return val, nil
}

func (m *MemoryService) GetObj(ctx context.Context, key string, obj interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if md, ok := m.data[key]; ok {
		if md.ttl != NoneDuration {
			if time.Now().Sub(md.createTime) > md.ttl {
				delete(m.data, key)
				return ErrKeyNotFound
			}
		}

		bt, _ := json.Marshal(md.val)
		err := json.Unmarshal(bt, obj)
		if err != nil {
			return err
		}
		return nil
	}

	return ErrKeyNotFound
}

func (m *MemoryService) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	errs := make([]string, 0, len(keys))
	rets := make(map[string]string)
	for _, key := range keys {
		val, err := m.Get(ctx, key)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		rets[key] = val
	}
	if len(errs) == 0 {
		return rets, nil
	}
	return rets, errors.New(strings.Join(errs, "||"))
}

func (m *MemoryService) Set(ctx context.Context, key string, val string, timeout time.Duration) error {
	return m.SetObj(ctx, key, val, timeout)
}

func (m *MemoryService) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	md := &MemoryData{
		val:        val,
		createTime: time.Now(),
		ttl:        timeout,
	}
	m.data[key] = md
	return nil
}

func (m *MemoryService) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	errs := []string{}
	for k, v := range data {
		err := m.Set(ctx, k, v, timeout)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "||"))
	}
	return nil
}

func (m *MemoryService) SetForever(ctx context.Context, key string, val string) error {
	return m.Set(ctx, key, val, NoneDuration)
}

func (m *MemoryService) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return m.SetObj(ctx, key, val, NoneDuration)
}

func (m *MemoryService) Remember(ctx context.Context, key string, timeout time.Duration, rememberFunc contract.RememberFunc, obj interface{}) error {
	err := m.GetObj(ctx, key, obj)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return err
	}

	objNew, err := rememberFunc(ctx, m.container)
	if err != nil {
		return err
	}

	if err := m.SetObj(ctx, key, objNew, timeout); err != nil {
		return err
	}

	if err := m.GetObj(ctx, key, obj); err != nil {
		return err
	}
	return nil
}

func (m *MemoryService) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if md, ok := m.data[key]; ok {
		md.ttl = timeout
		return nil
	}
	return ErrKeyNotFound
}

func (m *MemoryService) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if md, ok := m.data[key]; ok {
		return md.ttl, nil
	}
	return NoneDuration, ErrKeyNotFound
}

func (m *MemoryService) Calc(ctx context.Context, key string, step int64) (int64, error) {
	var val int64
	err := m.GetObj(ctx, key, &val)
	val = val + step
	if err == nil {
		m.data[key].val = val
		return val, nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return 0, err
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	m.data[key] = &MemoryData{
		val:        val,
		createTime: time.Now(),
		ttl:        NoneDuration,
	}

	return val, nil
}

func (m *MemoryService) Increment(ctx context.Context, key string) (int64, error) {
	return m.Calc(ctx, key, 1)
}

func (m *MemoryService) Decrement(ctx context.Context, key string) (int64, error) {
	return m.Calc(ctx, key, -1)
}

func (m *MemoryService) Del(ctx context.Context, key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.data, key)
	return nil
}

func (m *MemoryService) DelMany(ctx context.Context, keys []string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}
