package updateservice

/*type ChangeUUID struct {
	OldID uuid.UUID
	NewID uuid.UUID
}

type GetOpcUaNode struct {
	Info string
}

type UpdateService struct {
	ctx               context.Context
	clientChan        <-chan any
	opcuaChan         <-chan any
	updateToClienChan chan<- any
	updateToOpcUaChan chan<- any
	server            *server.ClientList
	OPCUAClient       *opcuaservice.OPCUAService
	Client            *clientservice.ClientService
}

func NewUpdateService(ctx context.Context, clientChan <-chan any, opcuaChan <-chan any, updateToClientChan chan<- any, updateToOpcUaChan chan<- any, server *server.ClientList) *UpdateService {
	var u UpdateService
	u.ctx = ctx
	u.clientChan = clientChan
	u.opcuaChan = opcuaChan
	u.updateToClienChan = updateToClientChan
	u.updateToOpcUaChan = updateToOpcUaChan
	u.server = server
	return &u
}

func (s *UpdateService) Update() {
	log.Println("Update service start")
	for {
		select {
		case data := <-s.clientChan:
			go s.clientRouter(data)
		case data := <-s.opcuaChan:
			log.Println("Get data from opcUA Chan")
			go s.opcuaRouter(data)
		case <-s.ctx.Done():
			log.Println("Connection done")
			return

		}
	}
}

func (s *UpdateService) clientRouter(data any) {
	switch v := data.(type) {
	case ChangeUUID:

	case GetOpcUaNode:
		//s.updateToOpcUaChan <- v
		/*Nodes, err := s.OPCUAClient.GetNodes()
		if err!=nil{
			log.Println(err.Error())
			return
		}
		data, err:=message.EncodeData(Nodes)
		if err!=nil{
			log.Println(err.Error())
			return
		}

	}
}

func (s *UpdateService) opcuaRouter(data any) {
	switch v := data.(type) {
	case update.SendOpcNodes:
		log.Println("v")
		s.updateToClienChan <- v
		log.Println("sendOPCUa Node")
	}
}*/
