package discovery

import (
	"encoding/json"
	"log"
	"sync/atomic"
)

// Service to register
type Service struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`

	lastIndex uint64
	quit      chan struct{}
}

// NewService returns a new Service
func NewService(name, host string, port int) *Service {
	return &Service{
		Name: name,
		Host: host,
		Port: port,
		quit: make(chan struct{}, 1),
	}
}

// Stop stops the heartbeat
func (srv *Service) Stop() {
	var s struct{}
	srv.quit <- s
}

//SetIndex sets lastIndex
func (srv *Service) SetIndex(index uint64) {
	atomic.StoreUint64(&srv.lastIndex, index)
}

// GetIndex get lastIndex
func (srv *Service) GetIndex() uint64 {
	return atomic.LoadUint64(&srv.lastIndex)
}

func (srv *Service) String() string {
	b, err := json.Marshal(srv)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(b)
}
