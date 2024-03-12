package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"web_app/settings"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
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

	_, err = rdb.Ping().Result()
	return
}

func Close() {
	_ = rdb.Close()
}
