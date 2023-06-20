package opcuaservice

import (
	"context"
	"log"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/spf13/viper"
	"github.com/sudak-91/monitoring/pkg/client"
)

type OPCUAService struct {
	opcuaChan              chan<- interface{}
	fromCommandToOpcuaChan <-chan interface{} //
	ctx                    context.Context
	OPCLient               *opcua.Client
}

func NewOpcUaService(ctx context.Context, opcuaChan chan<- interface{}, fromCommandToOpcuaChan <-chan interface{}) *OPCUAService {
	var opc OPCUAService
	opc.ctx = ctx
	opc.fromCommandToOpcuaChan = fromCommandToOpcuaChan
	opc.opcuaChan = opcuaChan
	return &opc
}

func (opc *OPCUAService) StartOPCUA(endpoint string) error {
	opts := []opcua.Option{}
	endpoints, err := opcua.GetEndpoints(context.TODO(), endpoint)
	secPolicy := ua.SecurityPolicyURINone
	secMode := ua.MessageSecurityModeNone

	authMode := ua.UserTokenTypeAnonymous
	authOption := opcua.AuthUsername("Administrator", "5dae40eb*")

	opts = append(opts, authOption)
	if err != nil {
		return err
	}
	var finallyEndpoit *ua.EndpointDescription
	for _, k := range endpoints {
		if k.SecurityPolicyURI == secPolicy || k.SecurityMode == secMode {
			finallyEndpoit = k
			log.Println(k.EndpointURL, k.SecurityMode, k.SecurityPolicyURI)
		}
	}
	secPolicy = finallyEndpoit.SecurityPolicyURI
	secMode = finallyEndpoit.SecurityMode
	secM := opcua.SecurityMode(secMode)
	opts = append(opts, secM)
	opts = append(opts, opcua.SecurityFromEndpoint(finallyEndpoit, authMode))

	opc.OPCLient = opcua.NewClient(endpoint)
	if err := opc.OPCLient.Connect(opc.ctx); err != nil {
		return err
	}
	log.Println("OPC UA Server start")
	go opc.CommandController()
	return nil

}

func (opc *OPCUAService) CommandController() {
	if viper.GetBool("Debug") {
		log.Println("OPCUA Command Controller is Run")
	}
	for {
		select {
		case data := <-opc.fromCommandToOpcuaChan:
			go func() {
				switch v := data.(type) {
				case client.GetOpcUaNodeTransfer:
					data, err := opc.GetNodes(v.Namespace, v.IID, v.SID)
					if err != nil {
						log.Printf("OPCUA Command controller has errod: %s", err.Error())
						v.Cancel()
						return
					}
					v.ResponseChan <- data
				case client.GetOpcUaNodeDescriptionTransfer:
					node := opc.GetNodeBySID(v.Namespace, v.SID)
					dataType, description, err := opc.GetNodeDescription(node)
					if err != nil {
						v.Cancel()
					}
					var transfer client.NodeDescriptionTransfer
					transfer.DataType = dataType
					transfer.Description = description
					v.ResponseChan <- transfer

				}
			}()
		}
	}
}
