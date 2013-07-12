package githubservicehook

import (
	"net/http"
	"sync"
)

type payloadProcessor func(Payload)

type hookProcess struct {
	processMutex sync.Mutex
	processor    payloadProcessor
	server       *http.Server
}

func (this *hookProcess) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := r.FormValue("payload")
	payload, err := parsePayload(body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go this.processNextPayload(payload)
}

func (this *hookProcess) processNextPayload(payload Payload) {
	this.processMutex.Lock()
	defer this.processMutex.Unlock()

	// do the thing
	this.processor(payload)
}

// this will block
func (this *hookProcess) Start(addr string) error {
	this.server = &http.Server{
		Addr:    addr,
		Handler: this,
	}

	return this.server.ListenAndServe()
}

func New(f payloadProcessor) *hookProcess {
	return &hookProcess{
		processMutex: sync.Mutex{},
		processor:    f,
	}
}
