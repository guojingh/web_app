package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
	"web_app/settings"
)

var db *sqlx.DB

func Init(cfg *settings.Mysql) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		/*viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),*/

		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	/*	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
		db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))*/

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

func Close() {
	_ = db.Close()
}
