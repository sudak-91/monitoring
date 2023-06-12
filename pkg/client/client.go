package client

import (
	"context"
	"log"

	"github.com/google/uuid"
	message "github.com/sudak-91/monitoring/pkg/message"
	command "github.com/sudak-91/monitoring/pkg/message/command"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	opcuaservice "github.com/sudak-91/monitoring/pkg/opcua_service"
	"nhooyr.io/websocket"
)

type Client struct {
	Conn                *websocket.Conn
	opcuaService        *opcuaservice.OPCUAService
	clientToServiceChan chan<- any
	UUID                uuid.UUID
	IsUUIDTemp          bool
}

func NewClient(conn *websocket.Conn, clientToServiceChan chan<- any, opcuaService *opcuaservice.OPCUAService) *Client {
	var c Client
	c.UUID = uuid.New()
	c.IsUUIDTemp = true
	c.Conn = conn
	c.opcuaService = opcuaService
	c.clientToServiceChan = clientToServiceChan
	return &c
}

func (c *Client) Run() {
	log.Println("Start read message")
	for {
		MessageType, data, err := c.Conn.Read(context.TODO())
		log.Printf("%s get message %v\n", c.UUID, MessageType)
		log.Printf("%v", data)
		command, err := message.Decode[command.Command](data)
		if err != nil {
			log.Println(err)
			continue
		}
		go c.messageRouter(command)
	}
	/*
	   mailoop:

	   	for {
	   		select {
	   		case data := <-cs.updateChan:
	   			log.Println("Get data from updateChan")
	   			go c.updateRouter(data)
	   		case <-ctx.Done():
	   			log.Println("Connection is odne")
	   			break mailoop
	   		default:

	   		}
	   	}
	   	log.Println("End")
	*/
}

func (c *Client) messageRouter(data command.Command) {
	log.Println("StartMessageRouting")
	switch {
	case data.GetUUID != nil:
		log.Println("GetUUID")
		data, err := c.GetUUIDHandle(c.UUID)
		if err != nil {
			log.Println(err)
		}
		c.Conn.Write(context.TODO(), websocket.MessageBinary, data)
	case data.SetUUID != nil:
		log.Println("setUUID")
		log.Println(c.UUID)
		err := c.SetUUIDHandle(data.SetUUID.UUID)
		if err != nil {
			log.Println(err)
		}
		log.Println(c.UUID)
		log.Println("finish")
	case data.GetOpcUaNode != nil:
		log.Println("getOPCNode")
		log.Println(c.UUID)
		data, err := c.getOpcUaNodeHandle()
		if err != nil {
			log.Println(err)
		}
		c.Conn.Write(context.TODO(), websocket.MessageBinary, data)
		log.Println("Command send")
	case data.GetSubNode != nil:
		log.Println("GetSubNode Commnd")
		log.Println(c.UUID)
		data, err := c.GetSubNodeHandle(data.GetSubNode.DOMParrentID, data.GetSubNode.IID, data.GetSubNode.SID, data.GetSubNode.Namespace)
		if err != nil {
			log.Println(err.Error())
			return
		}
		c.Conn.Write(context.TODO(), websocket.MessageBinary, data)
	case data.GetNodeDescription != nil:
		log.Println("[Command]|Get Node Description")
		log.Printf("[Client UUID]|ClientUUID: %s", c.UUID)
		data, err := c.GetNodeDescriptionHandler(data.GetNodeDescription.NS, data.GetNodeDescription.SID)
		if err != nil {
			log.Println(err.Error())
			return
		}
		c.Conn.Write(context.TODO(), websocket.MessageBinary, data)

	}
}

/*func (c *Client) updateRouter(data any) {
	switch v := data.(type) {
	case update.OPCNodes:
		var upd update.Update
		upd.OpcNodes = &v
		data, err := message.EncodeData(upd)
		if err != nil {
			log.Println(err)
		}
		err = c.Conn.Write(context.TODO(), websocket.MessageBinary, data)
		if err != nil {
			log.Println(err)
		}
		log.Println("sendOPC complete")
	}
}*/

func (c *Client) GetNodeDescriptionHandler(ns uint16, sid string) ([]byte, error) {
	node := c.opcuaService.GetNodeBySID(ns, sid)
	DataType, Description, err := c.opcuaService.GetNodeDescription(node)
	if err != nil {
		log.Printf("[GetNodeDescription]|%s", err.Error())
		return nil, err
	}
	update := update.NewNodeDescriptionUpdate(DataType, Description)
	return message.EncodeData(update)

}

func (c *Client) getOpcUaNodeHandle() ([]byte, error) {
	var upd update.Update
	data, err := c.opcuaService.GetNodes(0, 84, "")
	if err != nil {
		return nil, err
	}
	upd.OpcNodes = update.NewSendOpcNodes(&data)
	return message.EncodeData(upd)

}

func (c *Client) GetSubNodeHandle(parrentId string, id uint32, sid string, nodeNamespace uint16) ([]byte, error) {
	subNodes := update.NewOPCSubNodeUpdate(parrentId)
	OPCNodes, err := c.opcuaService.GetNodes(nodeNamespace, id, sid)
	if err != nil {
		return nil, err
	}
	subNodes.Nodes = OPCNodes
	update := subNodes.GetUpdate()
	return message.EncodeData(update)

}

func (c *Client) SetUUIDHandle(UUID string) error {
	newUUID, err := uuid.Parse(UUID)
	if err != nil {
		return err
	}
	var command ChangeUUID
	oldUUID := c.UUID
	c.UUID = newUUID
	command.NewUUID = newUUID
	command.OldUUID = oldUUID
	c.clientToServiceChan <- command
	return nil
}

func (c *Client) GetUUIDHandle(uuid uuid.UUID) ([]byte, error) {
	var upd update.Update
	upd.SendUUID = update.NewSendUUID(uuid)
	return message.EncodeData(upd)

}
