package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
	"github.com/mitchellh/mapstructure"
	"runtime"
	"sma_easy_helper/models"
)

// BaseController for the basic operations
type BaseController struct {
	beego.Controller
}

var MachineMap map[string]models.Machine

// CookiesCheck Check the Cookies is still valid
func (this *BaseController) CookiesCheck() {
	cookiesVal := this.GetSession("username")
	if cookiesVal == nil {
		beego.Info("the saved cookie is expired")
		this.Redirect("/login", 304)
	}
}

//init function is for initalizing all data
func init() {
	MachineMap = make(map[string]models.Machine)
}

func GCTest() {
	runtime.GC()
}

// @Title Get
// @Description get init data
// @Success 200 {object} models.Cluster
// @Failure 404 :error
// @router /cluster/read [get]
//InfoRead function is for reading the init data from file for cluster
func (c *BaseController) InfoClusterRead() {

	clusterString, err := models.FileRead("./saved_infos/init/init_cluster.json")
	if err != nil {
		MachineMap = nil
		beego.Error(err)
	} else {
		jsons, _ := simplejson.NewJson([]byte(clusterString))
		for _, jsonEntity := range jsons.MustArray() {
			cluster := models.Cluster{}
			err = mapstructure.WeakDecode(jsonEntity.(map[string]interface{}), &cluster)
			if err != nil {
				MachineMap = nil
				beego.Error(err)
			} else {
				for _, machineEntity := range cluster.Machines {
					MachineMap[machineEntity.Label] = machineEntity
				}
				c.Data["json"] = cluster
				c.ServeJSON()
			}
		}
	}
	c.Data["json"] = err
	c.ServeJSON()
}

// @Title Get
// @Description get init data
// @Success 200 {object} models.Operation
// @Failure 404 :error
// @router /operation/read [get]
//InfoOperationRead function is for handling the init data from file for operation
func (c *BaseController) InfoOperationRead(){
	operationString, err := models.FileRead("./saved_infos/init/init_operation.json")
	if err != nil {
		beego.Error(err)
	} else {
		jsons, _ := simplejson.NewJson([]byte(operationString))
		for _, jsonEntity := range jsons.MustArray() {
			operation := models.Operation{}
			err = mapstructure.WeakDecode(jsonEntity.(map[string]interface{}), &operation)
			if err != nil {
				beego.Error(err)
			} else {
				c.Data["json"] = operation
				c.ServeJSON()
			}
		}
	}
	c.Data["json"] = err
	c.ServeJSON()
}

// @Title InfoOperationWrite
// @Description write init operation data
// @Success 200 {object} models.Cluster
// @Failure 404 :error
// @router /cluster/write [post]
//InfoClusterWrite function is for writing the init data to file for cluster
func (c *BaseController) InfoClusterWrite(){
	var clusterEntity models.Cluster
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &clusterEntity)
	if err != nil {
		beego.Error(err)
	}
	dataString := "[" + string(data) + "]"
	beego.Info(dataString)
	err = models.FileWrite("./saved_infos/init/init_cluster.json", dataString)
	if err != nil {
		beego.Error(err)
	}
	c.Data["json"] = clusterEntity
	c.ServeJSON()
}

// @Title InfoOperationWrite
// @Description write init operation data
// @Success 200 {object} models.Operation
// @Failure 404 :error
// @router /operation/write [post]
//InfoOperationWrite function is for writing the init data to file for operation
func (c *BaseController) InfoOperationWrite() {
	var operationEntity models.Operation
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &operationEntity)
	if err != nil {
		beego.Error(err)
	}
	dataString := "[" + string(data) + "]"
	beego.Info(dataString)
	err = models.FileWrite("./saved_infos/init/init_operation.json", dataString)
	if err != nil {
		beego.Error(err)
	}
	c.Data["json"] = operationEntity
	c.ServeJSON()
}