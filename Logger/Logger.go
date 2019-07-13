package Logger

import (
	"fmt"
	"github.com/op/go-logging"
	"os"
)

// Password只是实现编校器接口的示例类型。任何
// 记录此日志时，将调用 Redacted()函数。
type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

func GetLogger(name string) (*logging.Logger) {
	var logger = logging.MustGetLogger(name)
	//示例格式字符串。除了消息之外，所有内容都有自定义颜色
	//取决于日志级别。许多字段都有自定义输出
	//格式化也一样。时间返回到毫秒。
	var format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	// 为os.Stderr创建两个后端.
	errBackend := logging.NewLogBackend(os.Stderr, "", 0)
	conBackend := logging.NewLogBackend(os.Stderr, "", 0)
	_, err := os.Stat(fmt.Sprintf("%s.log", name))
	var logFile  *os.File
	if err != nil{
		logFile, err = os.Create(fmt.Sprintf("%s.log", name))
		if err != nil{
			fmt.Printf("日志文件创建异常! Error=%v\n", err)
		}
	}else {
		logFile, err = os.OpenFile(fmt.Sprintf("%s.log", name), os.O_WRONLY, 0644)
		if err != nil{
			fmt.Printf("日志文件打开异常! Error=%v\n", err)
		}
	}
	fileBackend := logging.NewLogBackend(logFile,"", 0)

	//写入backend2的消息，添加一些额外的内容 包括使用的日志级别和名称
	conBackendFormatter := logging.NewBackendFormatter(conBackend, format)
	fileBackendFormatter := logging.NewBackendFormatter(fileBackend, format)

	// 错误和更严重的消息才发送到backend1
	errBackendLeveled := logging.AddModuleLevel(errBackend)
	errBackendLeveled.SetLevel(logging.ERROR, "")

	// 设置要使用的后端
	logging.SetBackend(errBackendLeveled, conBackendFormatter, fileBackendFormatter)
	logger.Info("Name=%s logger初始化完成!", name)
	return logger
}