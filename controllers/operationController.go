package controllers

import (
	"sma_easy_helper/models"
	"strconv"

	_ "github.com/astaxie/beego"
)

// OperationController is the controller for handling the operation requests
type OperationController struct {
	BaseController
}

// @Title GetAll
// @Description get all operations
// @Success 200 {object} models.Operation
// @router / [get]
func (o *OperationController) GetAll() {
	operations := models.GetAllOperations()
	o.Data["json"] = operations
	o.ServeJSON()
}

// @Title Get
// @Description get operation by name
// @Param	oname		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Operation
// @Failure 403 :cname is empty
// @router /:oname [get]
func (o *OperationController) Get() {
	oname := o.GetString(":oname")
	if oname != "" {
		intOname, err := strconv.Atoi(oname)
		operation, err := models.GetOperatoin(intOname)
		if err != nil {
			o.Data["json"] = err.Error()
		} else {
			o.Data["json"] = operation
		}
	}
	o.ServeJSON()
}

// Import function is for importing all saved operations informations to files
func (this *OperationController) Import() {

}

// Export function is for exporting all saved operations informations to client side
func (this *OperationController) Export() {

}
