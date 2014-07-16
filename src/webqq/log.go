package webqq

import (
	"github.com/astaxie/beego/logs"
)

var lg *logs.BeeLogger

func init() {
	lg = logs.NewLogger(1000)
	lg.SetLogger("console", `{}`)
	lg.SetLevel(logs.LevelTrace)
}
