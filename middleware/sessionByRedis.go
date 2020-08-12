package middleware

import (
	"MyServer/app"
	"MyServer/config"
	"fmt"
	"github.com/gin-contrib/sessions"
	session_redis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"time"
)

//  使用 redis 连接池
var RedisPool *redis.Pool

func sessionByRedis() gin.HandlerFunc {
	redisHost := config.GetStringFromConfig("redis.host")
	redisPort := config.GetStringFromConfig("redis.port")
	redisPassword := config.GetStringFromConfig("redis.password")
	reidsProtocol := config.GetStringFromConfig("redis.protocol")
	redisMaxIdle := config.GetIntFromConfig("redis.max_idle")
	redisMaxActive := config.GetIntFromConfig("redis.max_active")
	redisIdleTimeout := time.Duration(config.GetInt64FromConfig("redis.idle_timeout")) * time.Second
	redisTimeout := time.Duration(config.GetInt64FromConfig("redis.timeout")) * time.Second
	redisDB := config.GetIntFromConfig("redis.database")
	redisSecret := config.GetStringFromConfig("redis.secret")
	seesionCookieName := config.GetStringFromConfig("cookie.name")

	//   初始化  redis  连接池
	RedisPool = &redis.Pool{
		MaxIdle:     redisMaxIdle,     // 最大空闲连接数, 设为0表示无限制
		MaxActive:   redisMaxActive,   // 连接池的最大redis连接数。设为0表示无限制。
		IdleTimeout: redisIdleTimeout, // 空闲连接超时时间，超时的空闲连接会被关闭。应该设置一个比redis服务端超时时间更短的时间。设为0表示无限制
		Wait:        true,             // 如果Wait被设置成true，则Get()方法将会阻塞

		// Dial()方法返回一个连接，从在需要创建连接到的时候调用
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				reidsProtocol,                          // redis 使用的协议 tcp/udp
				redisHost+":"+redisPort,                // redis 地址:端口
				redis.DialPassword(redisPassword),      // reids 登录密码
				redis.DialDatabase(redisDB),            // reids 使用 哪个数据库，0-15
				redis.DialConnectTimeout(redisTimeout), // 连接 redis 的超时时间
				redis.DialReadTimeout(redisTimeout),    // 读取 redis 属性的超时时间
				redis.DialWriteTimeout(redisTimeout),   // 写入 redis 属性的超时时间
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

	store, err := session_redis.NewStoreWithPool(RedisPool, []byte(redisSecret))
	if err != nil {
		errMsg := fmt.Sprintf("Initial redis Err: %v", err)
		app.Logger.Error(errMsg)
		panic(err)
	}
	store.Options(sessions.Options{
		HttpOnly: true,
		Secure:   false,
		MaxAge:   73400,
	})

	// 指的是session的名字，也是cookie的名字
	return sessions.Sessions(seesionCookieName, store)
}
