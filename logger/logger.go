package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Init 初始化全局Logger
func Init() (err error) {
	// 1. 获取日志写入器 (整合 lumberjack 日志切割)
	writeSyncer := getLogWriter(
		viper.GetString("log.filename"),
		viper.GetInt("log.max_size"),
		viper.GetInt("log.max_backups"),
		viper.GetInt("log.max_age"),
	)

	// 2. 获取日志编码器 (定义日志的输出格式)
	encoder := getEncoder()

	// 3. 动态解析日志级别
	level := new(zapcore.Level)
	err = level.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil {
		return err
	}

	// 4. 创建核心 Core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 5. 生成 Logger，zap.AddCaller() 让输出附带调用者的文件名和行号
	log := zap.New(core, zap.AddCaller())

	// 6. 替换 zap 库中全局的 logger，以后可以在任意包直接使用 zap.L().Info(...)
	zap.ReplaceGlobals(log)

	return nil // 补充缺失的 return
}

// getEncoder 配置日志的编码格式
func getEncoder() zapcore.Encoder {
	// 获取 Zap 内置的生产环境默认配置 (自带一些字段键名的标准化)
	encoderConfig := zap.NewProductionEncoderConfig()

	// ----------------- 常见自定义调整 -----------------

	// 1. 时间格式：默认是不可读的 Unix 时间戳，改为 ISO8601 (人类可读，例如: 2023-10-25T12:00:00.000Z)
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 2. 日志级别格式：大写输出 (例如 INFO, ERROR 代替 info, error)，视觉上更醒目
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 3. 函数调用信息：短路径格式 (例如 main.go:25 代替一长串的绝对路径)
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 返回 JSON 格式的编码器。
	// 生产环境极力推荐 JSON，因为极其方便 Logstash、Fluentd 等日志采集工具进行解析。
	// 若是纯开发环境想在终端看彩色日志，可将其改为：zapcore.NewConsoleEncoder(encoderConfig)
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 配置日志文件的写入逻辑与自动切割
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	// 引入 lumberjack 库来实现基于文件大小和时间的日志切割
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 日志文件的全路径位置 (如 "./logs/app.log")
		MaxSize:    maxSize,   // 文件触发切割的最大大小（单位：MB）
		MaxBackups: maxBackup, // 系统保留的旧日志文件最大个数
		MaxAge:     maxAge,    // 旧日志文件保留的最大天数
		Compress:   false,     // 是否使用 gzip 压缩旧文件 (如果服务器 CPU 资源紧张可设为 false，推荐通过定时任务转移日志)
	}

	// 为了确保同时将日志写入文件和终端（方便本地排查），企业里有时会做多路输出：
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout))

	// 这里按原意，仅返回文件的写入同步器
	//return zapcore.AddSync(lumberJackLogger)
}

//	两个中间件，替换Gin默认的Recovery和Logger
//
// GinLogger 接收 Gin 框架默认的日志，替换原生的 gin.Logger()
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 执行下一个中间件及路由处理函数
		c.Next()

		cost := time.Since(start)

		// 记录请求相关信息
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),                                 // HTTP 状态码
			zap.String("method", c.Request.Method),                               // 请求方法
			zap.String("path", path),                                             // 请求路径
			zap.String("query", query),                                           // URL 参数
			zap.String("ip", c.ClientIP()),                                       // 客户端 IP
			zap.String("user-agent", c.Request.UserAgent()),                      // UserAgent
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // 内部错误信息
			zap.Duration("cost", cost),                                           // 耗时
		)
	}
}

// GinRecovery 捕获项目中可能出现的 panic，替换原生的 gin.Recovery()
// stack 参数决定是否在日志中记录完整的堆栈追踪信息
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取产生 panic 的请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 检查是否为底层的断开连接错误（Broken pipe）
				// 如果客户端断开了连接，这个时候我们不应该继续向其返回状态码
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 由于连接已断开，直接 abort，不写状态码
					c.Error(err.(error))
					c.Abort()
					return
				}

				// 根据是否需要记录堆栈信息，决定日志打印的详细程度
				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())), // 记录堆栈信息
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				// 向客户端返回 500 内部服务器错误
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		// 执行后续业务逻辑
		c.Next()
	}
}
