package log

import (
	mylog "gitee.com/lwj8507/nggs/log"

	"gitee.com/lwj8507/light-protoactor-go/actor"
	"gitee.com/lwj8507/light-protoactor-go/mailbox"
	"gitee.com/lwj8507/light-protoactor-go/remote"
)

var L mylog.ILogger = &mylog.ConsoleLogger{}

func Set(logger mylog.ILogger) {
	L = logger
	actor.SetLogger(logger)
	mailbox.SetLogger(logger)
	remote.SetLogger(logger)
}

func Close() {
	L.Close()
}
