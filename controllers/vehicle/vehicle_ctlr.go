package vehicle_ctlr

import (
	"../../filters/access_control"
	"github.com/astaxie/beego"
	"log"
	"time"
)

type YearController struct {
	beego.Controller
}

func (this *YearController) Get() {
	access_control.Tokenize(this.Ctx.ResponseWriter, this.Ctx.Request)

	log.Println("hit years")
	time.Sleep(100 * time.Millisecond)
	log.Println("finished years")

	this.Ctx.WriteString("get years")
}
