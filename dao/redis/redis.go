package redis

import (
	"fmt"
	"web_app/settings"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

func Init() (err error) {
	client = redis.NewClient(&redis.Options{
		/*		Addr: fmt.Sprintf("%s:%d",
				viper.GetString("redis.host"),
				viper.GetInt("redis.port"),
			),*/
		/*		Password: viper.GetString("redis.password"), // 密码
				DB:       viper.GetInt("redis.db"),          // 数据库
				PoolSize: viper.GetInt("redis.pool_size"), */ // 连接池大小
		Addr: fmt.Sprintf("%s:%d",
			settings.Conf.Redis.Host,
			settings.Conf.Redis.Port,
		),
		Password: settings.Conf.Redis.Password,
		DB:       settings.Conf.Redis.DB,
		PoolSize: settings.Conf.Redis.PoolSize,
	})

	_, err = client.Ping().Result()
	return
}

func Close() {
	_ = client.Close()
}
