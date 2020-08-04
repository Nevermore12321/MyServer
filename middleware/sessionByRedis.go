package middleware

import (
	"MyServer/config"
	"github.com/gin-contrib/sessions"
	session_redis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func sessionByRedis() gin.HandlerFunc {
	reidsHost := config.GetStringFromConfig("redis.host")
	reidsPort := config.GetStringFromConfig("redis.port")
	redisPassword := config.GetStringFromConfig("redis.password")
	reidsProtocol := config.GetStringFromConfig("redis.protocol")

	//  使用 redis 连接池
	var redisPool *redis.Pool

}
