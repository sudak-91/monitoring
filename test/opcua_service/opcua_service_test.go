package opcuaservice_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	opcuaservice "github.com/sudak-91/monitoring/pkg/opcua_service"
)

type HightLevel struct {
}

func (h *HightLevel) Error() string {
	return fmt.Sprintf("Hight level")
}

var (
	opcuChan   chan interface{}
	updateChan chan interface{}
)

func connectToOpcUa() (*opcuaservice.OPCUAService, error) {

	opcuaService := opcuaservice.NewOpcUaService(context.TODO(), opcuChan, updateChan)
	return opcuaService, nil
}

func TestGetDataValuesNode(t *testing.T) {
	var (
		opcuChan   chan interface{}
		updateChan chan interface{}
	)
	opcuaService := opcuaservice.NewOpcUaService(context.TODO(), opcuChan, updateChan)
	opcuaService.StartOPCUA("opc.tcp://192.168.1.225:4840")
	defer opcuaService.OPCLient.Close()
	RootNodes, err := opcuaService.GetNodes(0, 84, "")
	if err != nil {
		t.Error(err.Error())
	}
	for _, v := range RootNodes.OrganizesNode {

		log.Printf("[Root]ID: %d \t Namespace: %d \t Name: %s \n", v.IID, v.Namespace, v.Name)
		NodeID := ua.NewNumericNodeID(v.Namespace, v.IID)
		Node := opcuaService.OPCLient.Node(NodeID)
		t.Log(Node.ID.String())
		err := subNodes(opcuaService, Node, 0)
		switch {
		case errors.Is(err, &HightLevel{}):
			log.Println(err)
			continue
		case err != nil:
			log.Printf("[Error]| %s\n", err.Error())
			t.Fail()
		}

	}

}

func TestNode(t *testing.T) {
	var (
		opcuChan   chan interface{}
		updateChan chan interface{}
	)
	opcuaService := opcuaservice.NewOpcUaService(context.TODO(), opcuChan, updateChan)
	opcuaService.StartOPCUA("opc.tcp://192.168.1.225:4840")
	defer opcuaService.OPCLient.Close()
	NodeID := ua.NewNumericNodeID(4, 1001)
	Node := opcuaService.OPCLient.Node(NodeID)
	attr, err := opcuaService.GetDataValuesNode(Node)
	if err != nil {
		t.Error(err.Error())
	}
	for _, v := range attr {
		log.Printf("\t Data: %v\n", v.Value)
	}
}

func subNodes(service *opcuaservice.OPCUAService, node *opcua.Node, level int) error {
	if level > 10 {
		return &HightLevel{}
	}
	OrganizesNodesList, err := service.GetOrganizesNodes(node)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	for _, v := range OrganizesNodesList {
		name, err := v.BrowseName()
		if err != nil {
			log.Println(err.Error())
		}
		k := func(level int) string {
			var l string
			for i := 0; i < level; i++ {
				l = l + fmt.Sprint("\t")
			}
			return l
		}(level)
		log.Printf("%s[Organizes] NodeID:%s\t Namespace:%d\t BrowseName:%s\n", k, v.ID.StringID(), v.ID.Namespace(), name.Name)
		attrib, err := service.GetDataValuesNode(v)
		if err == nil {
			if attrib[0].Status == ua.StatusOK {
				log.Printf("[DataType] %d\n", attrib[0].Value.NodeID().IntID())
			}
			if attrib[1].Status == ua.StatusOK {
				log.Printf("[Value] %v\n", attrib[1].Value.Value())
			}
			if attrib[2].Status == ua.StatusOK {
				log.Printf("[NodeClass] %v\n", attrib[2].Value.Value())
			}
		}
		err = subNodes(service, v, level+1)
		if errors.Is(err, &HightLevel{}) {
			log.Println(err.Error())
			break
		}
	}
	ComponentNodesList, err := service.GetHasComponentNodes(node)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	for _, v := range ComponentNodesList {
		name, err := v.BrowseName()
		if err != nil {
			log.Println(err.Error())
		}
		k := func(level int) string {
			var l string
			for i := 0; i < level; i++ {
				l = l + fmt.Sprint("\t")
			}
			return l
		}(level)
		log.Printf("%s[Component] NodeID:%s\t Namespace:%d\t BrowseName:%s\n", k, v.ID.String(), v.ID.Namespace(), name.Name)
		attrib, err := service.GetDataValuesNode(v)
		if err == nil {
			if attrib[0].Status == ua.StatusOK {
				log.Printf("[DataType] %d\n", attrib[0].Value.NodeID().IntID())
			}
			if attrib[1].Status == ua.StatusOK {
				log.Printf("[Value] %v\n", attrib[1].Value.Value())
			}
			if attrib[2].Status == ua.StatusOK {
				log.Printf("[NodeClass] %v\n", attrib[2].Value.Value())
			}
		} else {
			return err
		}
		err = subNodes(service, v, level+1)
		if errors.Is(err, &HightLevel{}) {
			log.Println(err.Error())
			break
		}
	}
	PropertyNodesList, err := service.GetHasPropertyNodes(node)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	for _, v := range PropertyNodesList {
		name, err := v.BrowseName()
		if err != nil {
			log.Println(err.Error())
		}
		k := func(level int) string {
			var l string
			for i := 0; i < level; i++ {
				l = l + fmt.Sprint("\t")
			}
			return l
		}(level)

		log.Printf("%s[Property] NodeID:%s\t Namespace:%d\t BrowseName:%s\n", k, v.ID.String(), v.ID.Namespace(), name.Name)
		err = subNodes(service, v, level+1)
		if errors.Is(err, &HightLevel{}) {
			log.Println(err.Error())
			break
		}
	}

	return nil
}

func TestSubscribue(t *testing.T) {
	opcuaservice, err := connectToOpcUa()
	if err != nil {
		t.Fail()
	}
	err = opcuaservice.StartOPCUA("opc.tcp://192.168.1.225:4840")
	if err != nil {
		t.Fail()
	}
	defer opcuaservice.OPCLient.Close()

	var (
		nodeID    = "|var|CODESYS Control Win V3.Application.PLC_PRG.Status.wStatus"
		namespace = 4
	)
	expendedNode := ua.NewStringExpandedNodeID(uint16(namespace), nodeID)
	log.Println(expendedNode.NodeID, expendedNode.NodeID.IntID())
	Node := opcuaservice.OPCLient.Node(expendedNode.NodeID)
	attrs, err := opcuaservice.GetDataValuesNode(Node)
	if err != nil {
		t.Fail()
	}
	if attrs[0].Status == ua.StatusOK {
		log.Printf("[DataType] %d\n", attrs[0].Value.NodeID().IntID())
	}
	if attrs[1].Status == ua.StatusOK {
		log.Printf("[Value] %v\n", attrs[1].Value.Value())
	}

	notifyCh := make(chan *opcua.PublishNotificationData) //Channel for OPC UA pub message
	sub, err := opcuaservice.OPCLient.Subscribe(&opcua.SubscriptionParameters{Interval: opcua.DefaultSubscriptionInterval}, notifyCh)
	if err != nil {
		t.Fail()
	}
	defer sub.Cancel(context.TODO())
	monitoringItemRequest := opcua.NewMonitoredItemCreateRequestWithDefaults(expendedNode.NodeID, ua.AttributeIDValue, 42)
	_, err = sub.Monitor(ua.TimestampsToReturnBoth, monitoringItemRequest)
	if err != nil {
		t.Fail()
	}
	for {
		select {
		case res := <-notifyCh:
			if res.Error != nil {
				log.Println(res.Error.Error())
				continue
			}
			x, ok := res.Value.(*ua.DataChangeNotification)
			if !ok {
				continue
			}
			for _, v := range x.MonitoredItems {
				data := v.Value.Value.Value()
				log.Println(data)
			}
		default:
			continue
		}
	}

}
