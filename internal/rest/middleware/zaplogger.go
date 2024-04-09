package middleware

import (
	"dogker/lintang/monitor-service/config"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var lg *zap.Logger

// InitLogger 初始化Logger
func InitLogger(cfg config.LogConfig) (err error) {
	writeSyncerStdout, writeSyncerFile := getLogWriter(cfg.MaxBackups, cfg.MaxAge)
	encoder, fileEncoder := getEncoder()

	logLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	logLevelProd := zap.NewAtomicLevelAt(zap.InfoLevel)
	coreConsole := zapcore.NewCore(encoder, writeSyncerStdout, logLevel)
	coreFile := zapcore.NewCore(fileEncoder, writeSyncerFile, logLevelProd)
	core := zapcore.NewTee(
		coreConsole,
		coreFile,
	)
	lg = zap.New(core)
	zap.ReplaceGlobals(lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可

	return
}

func getEncoder() (encoder, fileEncoder zapcore.Encoder) {
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	productionCfg.EncodeDuration = zapcore.SecondsDurationEncoder
	productionCfg.EncodeCaller = zapcore.ShortCallerEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder = zapcore.NewJSONEncoder(productionCfg)
	return consoleEncoder, fileEncoder
}

func getLogWriter(maxBackup, maxAge int) (writeSyncerStdout zapcore.WriteSyncer, writeSyncerFile zapcore.WriteSyncer) {
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename: "./logs/app.log",

		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	})
	stdout := zapcore.AddSync(os.Stdout)

	return stdout, file
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		lg.Info(path,
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

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
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
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) //nolint: errcheck // asds
					c.Abort()
					return
				}

				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
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
