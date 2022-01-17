package redis

import (
	"time"
)

type RedisKey struct {
	Key    string
	Expire time.Duration
}

const RedisKey_UserOrders_Str string = "RK_USER_ORDERS"
const RedisKey_UserOrders_Expire time.Duration = time.Hour * 7

var RedisKey_UserOrders_Key RedisKey = RedisKey{Key: RedisKey_UserOrders_Str, Expire: RedisKey_UserOrders_Expire}
