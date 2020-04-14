package controllers

import (
	"runtime"

	"github.com/astaxie/beego"
)

// BaseController for the cookies check
type BaseController struct {
	beego.Controller
}

// CookiesCheck Check the Cookies is still valid
func (this *BaseController) CookiesCheck() {
	cookiesVal := this.GetSession("username")
	if cookiesVal == nil {
		beego.Info("the saved cookie is expired")
		this.Redirect("/login", 304)
	}
}

func GCTest() {
	runtime.GC()
}
