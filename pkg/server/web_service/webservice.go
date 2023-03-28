package webservice

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/sudak-91/monitoring/pkg/server"
	clientservice "github.com/sudak-91/monitoring/pkg/server/client_service"
	"nhooyr.io/websocket"
)

// Service for web socket connection
type WebService struct {
	Mutex      sync.RWMutex
	ctx        context.Context
	server     *server.Server
	updateChan chan interface{}
	clientChan chan any
}

func NewWebService(ctx context.Context, updateChan chan any, server *server.Server, clientChan chan any) *WebService {
	var service WebService
	service.ctx = ctx
	service.server = server
	service.updateChan = updateChan
	service.clientChan = clientChan
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
		log.Println("Connect Done")
		var client server.Client
		uuid := uuid.New()
		client.UUID = uuid
		client.IsUUIDTemp = true
		client.Conn = wsConnection
		cs := clientservice.NewClientService(service.ctx, &client, service.clientChan, service.updateChan)
		service.server.Users[uuid] = &client
		service.Mutex.Unlock()
		go cs.Run()
	})
	log.Println("HTTP Server start")
	if err := http.ListenAndServe("localhost:8000", mux); err != nil {
		panic(err)
	}
}
