package webserver

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	clientservice "github.com/sudak-91/monitoring/internal/pkg/client_service"
	"nhooyr.io/websocket"
)

type Server struct {
	Mutex      sync.RWMutex
	users      map[uuid.UUID]*clientservice.ClientService
	ctx        context.Context
	updateData chan interface{}
}

func NewServer(ctx context.Context) *Server {
	var Server Server
	Server.users = make(map[uuid.UUID]*clientservice.ClientService)
	Server.ctx = ctx
	Server.updateData = make(chan interface{}, 5)
	return &Server
}

func (s *Server) Start() {
	go s.Update()
	log.Println("Start HttpServer")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./template")))
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.Mutex.Lock()
		wsConnection, err := websocket.Accept(w, r, nil)
		if err != nil {
			s.Mutex.Unlock()
			log.Fatal(err)
		}
		log.Println("Connect Done")
		uuid := uuid.New()
		client := clientservice.NewClientService(s.ctx, uuid, wsConnection, s.updateData)
		s.users[uuid] = client
		s.Mutex.Unlock()
		go client.Run()
		/*defer wsConnection.Close(websocket.StatusInternalError, "falling")
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()
		var v any
		err = wsjson.Read(ctx, wsConnection, &v)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("%v", v)
		wsConnection.Close(websocket.StatusNormalClosure, "")*/
	})
	if err := http.ListenAndServe("localhost:8000", mux); err != nil {
		log.Fatal(err)
	}

}

func (s *Server) changeUUID(oldUUID uuid.UUID, newUUID uuid.UUID) error {
	s.users[newUUID] = s.users[oldUUID]
	delete(s.users, oldUUID)
	return nil
}
