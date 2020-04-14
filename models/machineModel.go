package models

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// Machine model is structing the machine model with below attribute
type Machine struct {
	MachineID int
	Label     string
	HostName  string
	HostIP    string
	UserName  string
	AuthType  string
	PassWord  string
	AuthKey   string
}

// wsBufferWriter is with attribute of byte buffer and mutext locker
type wsBufferWriter struct {
	buffer bytes.Buffer
	mux    sync.Mutex
}

type SSHConn struct {
	StdinPipe io.WriteCloser
	// Write() be called to receive data from ssh server
	ComboOutput *wsBufferWriter
	Session     *ssh.Session
}

type WsMsg struct {
	Type string `json:"type"`
	Cmd  string `json:"cmd"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

var (
	SSHMachines map[string]Machine
	MachineID   *int
)

const (
	wsMsgCmd    = "cmd"
	wsMsgResize = "resize"
)

func init() {

}

func NewSSHClient(machine Machine) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            machine.UserName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if machine.AuthType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(machine.PassWord)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(machine.PassWord)}
	}
	address := fmt.Sprintf("%s:%s", machine.HostIP, "22")
	conn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

func hostKeycallBackFunc(host string) ssh.HostKeyCallback {
	hostPath, err := homedir.Expand("~/.ssh/known_hosts")
	if err != nil {
		logrus.WithError(err).Error("find known_hosts's home dir failed")
	}
	file, err := os.Open(hostPath)
	if err != nil {
		logrus.WithError(err).Error("can't find known_host file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}
	if hostKey == nil {
		log.Fatalf("no hostkey for %s,%v", host, err)
	}
	return ssh.FixedHostKey(hostKey)
}

// implement Write interface to write bytes from ssh server into bytes.Buffer.
func (w *wsBufferWriter) Write(p []byte) (int, error) {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.buffer.Write(p)
}

func runCommand(client *ssh.Client, command string) (stdout string, err error) {
	session, err := client.NewSession()
	if err != nil {
		//log.Print(err)
		return
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		//log.Print(err)
		return
	}
	stdout = string(buf.Bytes())

	return
}

//flushComboOutput flush ssh.session combine output into websocket response
func flushComboOutput(w *wsBufferWriter, wsConn *websocket.Conn) error {
	if w.buffer.Len() != 0 {
		err := wsConn.WriteMessage(websocket.TextMessage, w.buffer.Bytes())
		if err != nil {
			return err
		}
		w.buffer.Reset()
	}
	return nil
}

// setup ssh shell session
// set Session and StdinPipe here,
// and the Session.Stdout and Session.Sdterr are also set.
// NewSSHConn function is for building the SSH Connection based on SSH Client
func NewSSHConn(cols, rows int, sshClient *ssh.Client) (*SSHConn, error) {
	sshSession, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	StdinPipe, err := sshSession.StdinPipe()
	if err != nil {
		return nil, err
	}

	comboWriter := new(wsBufferWriter)
	//ssh.stdout and stderr will write output into comboWriter
	sshSession.Stdout = comboWriter
	sshSession.Stderr = comboWriter

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     //disable echo
		ssh.TTY_OP_ISPEED: 14400, //input speed = 14.4Kbaud
		ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4Kbaud
	}
	// Request pseudo terminal
	if err := sshSession.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil, err
	}
	// Start remote Shell
	if err := sshSession.Shell(); err != nil {
		return nil, err
	}
	return &SSHConn{StdinPipe: StdinPipe, ComboOutput: comboWriter, Session: sshSession}, nil
}

//ReceiveWsMsg  receive websocket msg do some handling then write into ssh.session.stdin
func (ssConn *SSHConn) ReceiveWsMsg(wsConn *websocket.Conn, exitCh chan bool) {
	//tells other go routine quit
	defer setQuit(exitCh)
	for {
		select {
		case <-exitCh:
			return
		default:
			//read websocket msg
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				logrus.WithError(err).Error("reading webSocket message failed")
				return
			}
			//unmashal bytes into struct
			msgObj := WsMsg{}
			if err := json.Unmarshal(wsData, &msgObj); err != nil {
				beego.Error(err)
				logrus.WithError(err).WithField("wsData", string(wsData)).Error("unmarshal websocket message failed")
			}
			switch msgObj.Type {
			case wsMsgResize:
				//handle xterm.js size change
				if msgObj.Cols > 0 && msgObj.Rows > 0 {
					if err := ssConn.Session.WindowChange(msgObj.Rows, msgObj.Cols); err != nil {
						logrus.WithError(err).Error("ssh pty change windows size failed")
					}
				}
			case wsMsgCmd:
				//handle xterm.js stdin
				decodeBytes, err := base64.StdEncoding.DecodeString(msgObj.Cmd)
				if err != nil {
					logrus.WithError(err).Error("websock cmd string base64 decoding failed")
				}
				if _, err := ssConn.StdinPipe.Write(decodeBytes); err != nil {
					logrus.WithError(err).Error("ws cmd bytes write to ssh.stdin pipe failed")
				}
				//write input cmd to log buffer
				// if _, err := logBuff.Write(decodeBytes); err != nil {
				// 	logrus.WithError(err).Error("write received cmd into log buffer failed")
				// }
			}
		}
	}
}

func (ssConn *SSHConn) SendComboOutput(wsConn *websocket.Conn, exitCh chan bool) {
	//tells other go routine quit
	defer setQuit(exitCh)

	//every 120ms write combine output bytes into websocket response
	tick := time.NewTicker(time.Millisecond * time.Duration(120))
	//for range time.Tick(120 * time.Millisecond){}
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			//write combine output bytes into websocket response
			if err := flushComboOutput(ssConn.ComboOutput, wsConn); err != nil {
				logrus.WithError(err).Error("ssh sending combo output to webSocket failed")
				return
			}
		case <-exitCh:
			return
		}
	}
}

func (sshConn *SSHConn) SessionWait(quitChan chan bool) {
	if err := sshConn.Session.Wait(); err != nil {
		logrus.WithError(err).Error("ssh session wait failed")
		setQuit(quitChan)
	}
}

func setQuit(ch chan bool) {
	ch <- true
}

func (s *SSHConn) close() {
	if s.Session != nil {
		s.Session.Close()
	}
}
