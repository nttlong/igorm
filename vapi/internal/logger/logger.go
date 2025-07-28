package logger

import (
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	FileName   string
	MaxSize    int // MB
	MaxAge     int // ngày
	MaxBackups int
	Compress   bool
	AlsoStdout bool
}
type LoggerService struct {
	Output io.Writer
	Config LoggerConfig
}

func NewLoggerService(config *LoggerConfig) *LoggerService {
	if config == nil {
		config = &LoggerConfig{
			FileName:   "app.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
			AlsoStdout: true,
		}
	}
	logWriter := &lumberjack.Logger{
		Filename:   config.FileName,
		MaxSize:    config.MaxSize,    // MB
		MaxBackups: config.MaxBackups, // giữ lại 3 file log cũ
		MaxAge:     config.MaxAge,     // ngày
		Compress:   config.Compress,   // gzip
	}

	var output io.Writer = logWriter
	if config.AlsoStdout {
		output = io.MultiWriter(os.Stdout, logWriter)
	}

	return &LoggerService{Output: output}
}

func (l *LoggerService) Apply(e *echo.Echo) {
	e.Logger.SetOutput(l.Output)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: l.Output,
		Format: `${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human}` + "\n",
	}))
	e.Use(middleware.Recover())
}

// Apply gắn logger cho Echo
