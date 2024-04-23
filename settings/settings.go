package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Conf = new(ConfApp)

type ConfApp struct {
	App   *App   `mapstructure:"App"`
	Log   *Log   `mapstructure:"Log"`
	Mysql *Mysql `mapstructure:"Mysql"`
	Redis *Redis `mapstructure:"Redis"`
}

type App struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Host      string `mapstructure:"host"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`
	Page      int64  `mapstructure:"page"`
	Size      int64  `mapstructure:"size"`
}

type Log struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type Mysql struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(fileName string) (err error) {
	//viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	//viper.SetConfigName("config")        // 配置文件名称(无扩展名)
	//viper.SetConfigType("yaml")          // 如果配置文件的名称中没有扩展名，则需要配置此项 基本上配合远程配置中心使用

	viper.SetConfigFile(fileName)
	//viper.AddConfigPath("./")  // 查找配置文件所在的路径
	//viper.AddConfigPath(".")   // 还可以在工作目录中查找配置
	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误
		//读取配置信息失败
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}

	if err = viper.Unmarshal(Conf); err != nil {
		zap.L().Error("unmarshal conf failed", zap.Error(err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")
		if err = viper.Unmarshal(Conf); err != nil {
			zap.L().Error("unmarshal conf failed", zap.Error(err))
		}
	})
	return
}
