package contract

import "time"

const ConfigKey = "app:config"

type Config interface {
	IsExist(key string) bool // IsExist检查一个属性是否存在

	Get(key string) interface{} // Get获取一个属性值

	GetBool(key string) bool // GetBool获取一个bool属性

	GetInt(key string) int // GetInt获取一个int属性

	GetFloat64(key string) float64 // GetFloat64获取一个float64属性

	GetTime(key string) time.Time // GetTime获取一个time属性

	GetString(key string) string // GetString获取一个string属性

	GetIntSlice(key string) []int // GetIntSlice获取一个int数组属性

	GetStringSlice(key string) []string // GetStringSlice获取一个string数组

	GetStringMap(key string) map[string]interface{} // GetStringMap获取一个string为key，interface为val的map

	GetStringMapString(key string) map[string]string // GetStringMapString获取一个string为key，string为val的map

	GetStringMapStringSlice(key string) map[string][]string // GetStringMapStringSlice获取一个string为key，数组string为val的map

	Load(key string, val interface{}) error // Load加载配置到某个对象
}
