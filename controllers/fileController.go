package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego"
	"github.com/pkg/sftp"
	"sma_easy_helper/models"
)

var (
	sftpConn *sftp.Client
)

// FileController is the controller for handling the file requests
type FileController struct {
	BaseController
}

func SFTPConnGet(c *FileController) (sftpConn *sftp.Client, err error) {

	sshConn, err := models.NewSshClient(SSHHost)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	sftpConn, err = sftp.NewClient(sshConn)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return sftpConn, nil
}

// @Title ListAll
// @Description get all Files in this Directory Path
// @Success 200 {object} models.DirectoryList
// @router /list [get]
// List function is for getting the file content from the known host
func (c *FileController) List() {

	sftpConn, _ = SFTPConnGet(c)
	dirPath := c.Input().Get("dir")
	if dirPath != "" {
		dirList, err := models.SFTPFileDirList(dirPath, sftpConn)
		if err != nil {
			beego.Error(err)
			c.Data["json"] = err.Error()
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
// @router /read [get]
//Read function is for handling file content read
func (c *FileController) Read() {

	sftpConn, _ = SFTPConnGet(c)
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

// @Title save
// @Description write the file with assigned content
// @Success 200 {object} models.File
// @router /write [post]
// Save function is for saving the file content edited by customer in UI
func (c *FileController) Save() {

	sftpConn, _ = SFTPConnGet(c)
	filePath := c.Input().Get("path")
	var writeFile models.File
	writeFile.FileName = filePath
	writeFile.FilePath = filePath
	err := models.SFTPFileWrite(writeFile, sftpConn)
	if err != nil {
		beego.Error(err)
		c.Data["json"] = err.Error()
		c.ServeJSON()
	}
}

// Compare function is for comparing the origin and new files content
func (c *FileController) Compare() {

}
