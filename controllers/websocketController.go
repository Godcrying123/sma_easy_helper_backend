package controllers

import (
	"net/http"
	"webconsole_sma/utils"

	"sma_easy_helper/models"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// WebSocketController is the controller for handing remote SSH connection
type SSHWebSocketController struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Get function is for connecting the remote SSH machine
func (this *SSHWebSocketController) Get() {
	sshHost := models.Machine{
		MachineID: 1,
		Label:     "Master",
		HostName:  "test-machine-01",
		HostIP:    "16.186.79.151",
		UserName:  "root",
		AuthType:  "password",
		PassWord:  "iso*help",
		AuthKey:   "/test.crt",
	}
	sshClient, err := models.NewSSHClient(sshHost)
	if err != nil {
		beego.Error(err)
	}
	defer sshClient.Close()
	// startTime := time.Now()
	sshConn, err := utils.NewSshConn(120, 32, sshClient)
	if err != nil {
		beego.Error(err)
	}
	defer sshConn.Close()
	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if err != nil {
		beego.Error(err)
	}
	defer ws.Close()
	utils.SSHClients[ws] = true
	quitChan := make(chan bool, 3)
	go sshConn.ReceiveWsMsg(ws, quitChan)
	go sshConn.SendComboOutput(ws, quitChan)
	go sshConn.SessionWait(quitChan)
	<-quitChan
}
