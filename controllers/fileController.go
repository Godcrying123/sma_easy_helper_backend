package controllers

import (
	_ "github.com/astaxie/beego"
)

// FileController is the controller for handling the file requests
type FileController struct {
	BaseController
}

// Get function is for getting the file content from the known host
func (this *FileController) Get() {

}

// Save function is for saving the file content edited by customer in UI
func (this *FileController) Save() {

}

// Compare function is for comparing the origin and new files content
func (this *FileController) Compare() {

}
