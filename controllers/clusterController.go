package controllers

import (
	"sma_easy_helper/models"

	_ "github.com/astaxie/beego"
)

// ClusterController is the controller for handling the cluster requests & response
type ClusterController struct {
	BaseController
}

// @Title GetAll
// @Description get all Cluster
// @Success 200 {object} models.Cluster
// @router / [get]
func (c *ClusterController) GetAll() {
	clusters := models.GetAllClusters()
	c.Data["json"] = clusters
	c.ServeJSON()
}

// @Title Get
// @Description get cluster by name
// @Param	cname		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 404 :cname is empty
// @router /:cname [get]
func (c *ClusterController) Get() {
	cname := c.GetString(":cname")
	if cname != "" {
		cluster, err := models.GetCluster(cname)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = cluster
		}
	}
	c.ServeJSON()
}
