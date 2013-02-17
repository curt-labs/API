package home

import (
	"github.com/astaxie/beego"
)

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Get() {
	this.Redirect("http://labs.curtmfg.com", 302)
}
