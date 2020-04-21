package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego"
	"github.com/pkg/sftp"
	"sma_easy_helper/models"
)

var sftpConn  *sftp.Client

// FileController is the controller for handling the file requests
type FileController struct {
	BaseController
}

// @Title ListAll
// @Description get all Files in this Directory Path
// @Success 200 {object} models.DirectoryList
// @router / [get]
// List function is for getting the file content from the known host
func (c *FileController) List() {
	dirPath := c.Input().Get("dir")
	if dirPath != "" {
		dirList, err := models.SFTPFileDirList(dirPath, sftpConn)
		if err != nil {
			beego.Error(err)
			c.Data["json"] = err
			c.ServeJSON()
		} else {
			c.Data["json"] = dirList
			c.ServeJSON()
		}
	}
}

// @Title Read
// @Description get the File all attribute for
// @Success 200 {object} models.File
// @router /read/ [get]
//Read function is for handling file content read
func (c *FileController) Read() {
	filePath := c.Input().Get("path")
	var readFile models.File
	readFile.FileName = filePath
	readFile.FilePath = filePath
	readFile, err := models.SFTPFileRead(readFile, sftpConn)
	if err != nil {
		beego.Error(err)
		c.Data["json"] = err
		c.ServeJSON()
	} else {
		c.Data["json"] = readFile
		c.ServeJSON()
	}
}

// Save function is for saving the file content edited by customer in UI
func (c *FileController) Save() {

}

// Compare function is for comparing the origin and new files content
func (c *FileController) Compare() {

}
