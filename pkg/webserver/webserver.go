package webserver

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Server struct {
	Mutex sync.RWMutex
	users map[uuid.UUID]*WebUser
	ctx   context.Context
}

func NewServer(ctx context.Context) *Server {
	var Server Server
	Server.users = make(map[uuid.UUID]*WebUser)
	Server.ctx = ctx
	return &Server
}

func (s *Server) Start() {
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
		var webUser WebUser
		webUser.Conn = wsConnection
		webUser.ctx, webUser.cancel = context.WithCancel(s.ctx)
		uuid := uuid.New()
		s.users[uuid] = &webUser
		s.Mutex.Unlock()
		go webUser.Run(uuid.String())
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
