package controllers

import "github.com/astaxie/beego"

type SetupController struct {
	beego.Controller
}

func (c *SetupController) Get(){
	c.TplName = "index.html"
}