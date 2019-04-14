package mailbox

import (
	mylog "gitee.com/lwj8507/nggs/log"
)

var logger mylog.ILogger = &mylog.ConsoleLogger{}

func SetLogger(l mylog.ILogger) {
	logger = l
}
