package webqq

import (
	"github.com/astaxie/beego/logs"
)

var lg *logs.BeeLogger

func init() {
	lg = logs.NewLogger(1000)
	lg.SetLogger("console", `{"level":2}`)
	lg.SetLogger("file", `{"filename":"log.log"}`)
	lg.SetLevel(logs.LevelTrace)
}
