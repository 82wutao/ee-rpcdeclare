package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr     string
	Port     int16
	Password string
	DB       int
}

// func getRedisConfig() *RedisConfig {
// 	cfg := RedisConfig{
// 		Addr:     "localhost",
// 		Port:     6379,
// 		Password: "",
// 		DB:       0,
// 	}
// 	return &cfg
// }

var cli *redis.Client

func Connect2redisServer(cfg *RedisConfig, ctx context.Context) {
	opt := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	for running := true; running; running = false {
		if cli == nil {
			break
		}
		if err := cli.Close(); err != nil {
			log.Printf("close redis cli error %+v\n", err)
		}
	}
	cli = redis.NewClient(&opt)
	cli.Conn(ctx)
	log.Printf("debug : redis cli %+v\n", cli)
}
func CloseRedisConnection() {
	if cli == nil {
		return
	}
	cli.Close()
	cli = nil
}

func GetRedisCli() interface{} {
	return cli
}
func RedisTryAcquire(ctx context.Context, k string, v interface{}, expire time.Duration) (bool, error) {
	boolCmd := cli.Conn(context.Background()).SetNX(ctx, k, v, expire)
	return boolCmd.Val(), boolCmd.Err()
}

func RedisSet(ctx context.Context, k RedisKey, v interface{}) error {
	statusCmd := cli.Set(ctx, k.Key, v, k.Expire)
	return statusCmd.Err()
}
func RedisIncreaseBy(ctx context.Context, k string, v int64) (int64, error) {
	intCmd := cli.IncrBy(ctx, k, v)
	return intCmd.Val(), intCmd.Err()
}
func RedisDecreaseBy(ctx context.Context, k string, v int64) (int64, error) {
	intCmd := cli.DecrBy(ctx, k, v)
	return intCmd.Val(), intCmd.Err()
}

func RedisDel(ctx context.Context, k string) (bool, error) {
	intCmd := cli.Del(ctx, k)
	return intCmd.Val() == 1, intCmd.Err()
}
func RedisExpire(ctx context.Context, k string, e time.Duration) (bool, error) {
	boolCmd := cli.Expire(ctx, k, e)
	return boolCmd.Val(), boolCmd.Err()
}
func RedisList(ctx context.Context, k RedisKey, v []interface{}) (int64, error) {
	intCmd := cli.LPush(ctx, k.Key, v...)
	if intCmd.Err() != nil {
		return 0, intCmd.Err()
	}
	_, err := RedisExpire(ctx, k.Key, k.Expire)
	if err != nil {
		return intCmd.Val(), err
	}

	return intCmd.Val(), nil
}

// list 列表（实现队列,元素不唯一，先入先出原则）
// set 集合（各不相同的元素）
// hash hash散列值（hash的key必须是唯一的）
// sort set 有序集合

// 3.list类型支持的常用命令：
// lpush:从左边推入
// lpop:从右边弹出
// rpush：从右变推入
// rpop:从右边弹出
// llen：查看某个list数据类型的长度

// 4.set类型支持的常用命令：
// sadd:添加数据
// scard:查看set数据中存在的元素个数
// sismember:判断set数据中是否存在某个元素
// srem:删除某个set数据中的元素

// 5.hash数据类型支持的常用命令:
// hset:添加hash数据
// hget:获取hash数据
// hmget:获取多个hash数据

// 6.sort set和hash很相似,也是映射形式的存储:
// zadd:添加
// zcard:查询
// zrange:数据排序
