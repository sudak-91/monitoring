package webserver

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/sudak-91/monitoring/pkg/clientservice"
	"nhooyr.io/websocket"
)

type WebService struct {
	Mutex         sync.RWMutex
	ctx           context.Context
	clientService *clientservice.ClientService
}

func NewWebService(ctx context.Context, clientService *clientservice.ClientService) *WebService {
	var service WebService
	service.ctx = ctx
	service.clientService = clientService
	return &service
}

func (service *WebService) Run() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./template")))
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		service.Mutex.Lock()
		wsConnection, err := websocket.Accept(w, r, nil)
		if err != nil {
			service.Mutex.Unlock()
			log.Fatal(err)
		}
		service.clientService.NewClient(wsConnection)
		log.Println("Connect Done")
	})
	log.Println("HTTP Server start")
	if err := http.ListenAndServe("localhost:8000", mux); err != nil {
		panic(err)
	}
}
