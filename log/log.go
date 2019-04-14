package log

import (
	mylog "gitee.com/lwj8507/nggs/log"

	"github.com/hei6775/light-protoactor-go/actor"
	"github.com/hei6775/light-protoactor-go/mailbox"
	"github.com/hei6775/light-protoactor-go/remote"
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
