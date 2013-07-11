package githubservicehook

import (
	"container/list"
	"io/ioutil"
	"net/http"
	"sync"
)

type PayloadProcessor func(Payload)

type hookProcess struct {
	list         *list.List
	listMutex    sync.RWMutex
	processMutex sync.Mutex
	processor    PayloadProcessor
	addr         string

	server *http.Server
}

func (this *hookProcess) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	payload, err := parsePayload(body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	this.listMutex.Lock()
	this.list.PushBack(payload)
	this.listMutex.Unlock()

	go this.processNextPayload()
}

func (this *hookProcess) processNextPayload() {
	this.processMutex.Lock()
	defer this.processMutex.Unlock()

	// grab the front element
	this.listMutex.RLock()
	payload := this.list.Front()
	this.listMutex.RUnlock()

	if payload == nil {
		return
	}

	// do the thing
	this.processor(payload.Value.(Payload))

	// remove the payload
	this.listMutex.Lock()
	this.list.Remove(payload)
	this.listMutex.Unlock()

}

// this will block
func (this *hookProcess) Start() error {
	this.server = &http.Server{
		Addr:    this.addr,
		Handler: this,
	}

	return this.server.ListenAndServe()
}

func New(addr string, f PayloadProcessor) *hookProcess {
	return &hookProcess{
		list:         list.New(),
		listMutex:    sync.RWMutex{},
		processMutex: sync.Mutex{},
		processor:    f,
		addr:         addr,
	}
}
