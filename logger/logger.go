package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
	"web_app/settings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(cfg *settings.Log, mode string) (err error) {

	//配置日志写入规则
	writeSyncer := getLogWriter(
		/*		viper.GetString("log.filename"),
				viper.GetInt("log.max_size"),
				viper.GetInt("log.max_backups"),
				viper.GetInt("log.max_age"),*/

		/*		settings.Conf.Log.FileName,
				settings.Conf.Log.MaxSize,
				settings.Conf.Log.MaxBackups,
				settings.Conf.Log.MaxAge,*/
		cfg.FileName,
		cfg.MaxAge,
		cfg.MaxSize,
		cfg.MaxBackups,
	)
	//配置日志输出格式
	encoder := getEncoder()
	var l = new(zapcore.Level)
	//err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	err = l.UnmarshalText([]byte(settings.Conf.Log.Level))

	if err != nil {
		return
	}

	var core zapcore.Core
	if mode == "dev" {
		//开发模式日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	//logger = zap.New(core, zap.AddCaller())
	lg := zap.New(core, zap.AddCaller())

	//替换全局的 logger
	zap.ReplaceGlobals(lg)
	return
}

// 配置日志输出格式
func getEncoder() zapcore.Encoder {
	//json 格式
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//控制台格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
	//return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter(fileName string, maxSize int, maxBackups int, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,    // M
		MaxBackups: maxBackups, //最大备份数量
		MaxAge:     maxAge,     //最大备份天数
		Compress:   false,      // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
