package Log

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func init() {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func NewLogger(filename string, MaxBackups, MaxAge, MaxSize int) *zerolog.Logger {
	//TODO 后续使用cli将项目目录传递进去
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		panic(errors.New("项目根目录未找到"))
	}
	l := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", projectRoot, filename),
		MaxSize:    MaxSize,
		MaxAge:     MaxAge,
		MaxBackups: MaxBackups,
	}
	//输出进日志和前台
	writer := zerolog.MultiLevelWriter(l, os.Stdout)
	logger := zerolog.New(writer).With().Timestamp().Logger()
	return &logger
}
