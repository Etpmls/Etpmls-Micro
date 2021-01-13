package em_library

import (
	"context"
	Package_Redis "github.com/go-redis/redis/v8"
	"time"
)

var Instance_Redis *Package_Redis.Client

func Init_Redis(enableCache bool, address string, password string, db int)  {
	//	If the cache is not turned on, skip redis initialization
	//	如果没有开启缓存则跳过redis初始化
	if !enableCache {
		return
	}

	Instance_Redis = Package_Redis.NewClient(&Package_Redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	_, err := Instance_Redis.Ping(context.TODO()).Result()
	if err != nil {
		InitLog.Fatalln("[WARNING]", "Redis initialization failed.", " Error:", err)
	} else {
		InitLog.Println("[Info]", "redis initialized successfully.")
	}
}


type redis struct {}

func NewRedis() *redis {
	return &redis{}
}

// Get Json
// 获取字符串
func (this *redis) GetString (key string) (string, error) {
	return Instance_Redis.Get(context.Background(), key).Result()
}


// Set Json
// 设置字符串
func (this *redis) SetString (key string, value string, time time.Duration) {
	_ = Instance_Redis.Set(context.Background(), key, value, time).Err()
	return
}


// Delete Json
// 删除字符串
func (this *redis) DeleteString (list ...string) {
	_ = Instance_Redis.Del(context.Background(), list...).Err()
	return
}


// Get Hash
// 获取哈希
func (this *redis) GetHash (key string, field string) (string, error) {
	return Instance_Redis.HGet(context.Background(), key, field).Result()
}


// Set Hash
// 设置哈希
func (this *redis) SetHash (key string, value map[string]string) {
	var tmp = make(map[string]interface{})
	for k, v := range value {
		tmp[k] = v
	}
	_ = Instance_Redis.HSet(context.Background(), key, tmp).Err()
	return
}


// Delete Hash
// 删除哈希
func (this *redis) DeleteHash (key string, list ...string) {
	_ = Instance_Redis.HDel(context.Background(), key, list...).Err()
	return
}


// Clear all caches in the current DB
// 清除当前DB内所有缓存
func (this *redis) ClearAllCache() {
	Instance_Redis.FlushDB(context.Background())
	return
}



