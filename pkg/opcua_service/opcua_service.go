package opcuaservice

import (
	"context"
	"log"
	"strconv"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
	"github.com/pkg/errors"
	"github.com/sudak-91/monitoring/pkg/client"
)

type NodeDef struct {
	NodeID      *ua.NodeID
	NodeClass   ua.NodeClass
	BrowseName  string
	Description string
	AccessLevel ua.AccessLevelType
	Path        string
	DataType    string
	Writable    bool
	Unit        string
	Scale       string
	Min         string
	Max         string
}

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
	//defer opc.c.CloseWithContext(opc.ctx)
	//@TODO: Пока заблокированный цикл
	/*for {
		select {
		case data := <-opc.updateChan:
			log.Println("new update")
			opc.router(data)
		default:
			continue

		}
	}*/
	return nil

}

func (opc *OPCUAService) CommandController() {
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

/*func (opc *OPCUAService) router(data any) {
	switch v := data.(type) {
	case updateservice.GetOpcUaNode:
		log.Println(v.Info)
		go opc.GetNodes()

	}
}*/

func (n NodeDef) Records() []string {
	return []string{n.BrowseName, n.DataType, n.NodeID.String(), n.Unit, n.Scale, n.Min, n.Max, strconv.FormatBool(n.Writable), n.Description}
}

func join(a, b string) string {
	if a == "" {
		return b
	}
	return a + "." + b
}

func browse(ctx context.Context, n *opcua.Node, path string, level int) ([]NodeDef, error) {
	// fmt.Printf("node:%s path:%q level:%d\n", n, path, level)
	if level > 10 {
		return nil, nil
	}

	attrs, err := n.AttributesWithContext(ctx, ua.AttributeIDNodeClass, ua.AttributeIDBrowseName, ua.AttributeIDDescription, ua.AttributeIDAccessLevel, ua.AttributeIDDataType)
	if err != nil {
		return nil, err
	}

	var def = NodeDef{
		NodeID: n.ID,
	}

	switch err := attrs[0].Status; err {
	case ua.StatusOK:
		def.NodeClass = ua.NodeClass(attrs[0].Value.Int())
	default:
		return nil, err
	}

	switch err := attrs[1].Status; err {
	case ua.StatusOK:
		def.BrowseName = attrs[1].Value.String()
	default:
		return nil, err
	}

	switch err := attrs[2].Status; err {
	case ua.StatusOK:
		def.Description = attrs[2].Value.String()
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return nil, err
	}

	switch err := attrs[3].Status; err {
	case ua.StatusOK:
		def.AccessLevel = ua.AccessLevelType(attrs[3].Value.Int())
		def.Writable = def.AccessLevel&ua.AccessLevelTypeCurrentWrite == ua.AccessLevelTypeCurrentWrite
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return nil, err
	}

	switch err := attrs[4].Status; err {
	case ua.StatusOK:
		switch v := attrs[4].Value.NodeID().IntID(); v {
		case id.DateTime:
			def.DataType = "time.Time"
		case id.Boolean:
			def.DataType = "bool"
		case id.SByte:
			def.DataType = "int8"
		case id.Int16:
			def.DataType = "int16"
		case id.Int32:
			def.DataType = "int32"
		case id.Byte:
			def.DataType = "byte"
		case id.UInt16:
			def.DataType = "uint16"
		case id.UInt32:
			def.DataType = "uint32"
		case id.UtcTime:
			def.DataType = "time.Time"
		case id.String:
			def.DataType = "string"
		case id.Float:
			def.DataType = "float32"
		case id.Double:
			def.DataType = "float64"
		default:
			def.DataType = attrs[4].Value.NodeID().String()
		}
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return nil, err
	}

	def.Path = join(path, def.BrowseName)
	// fmt.Printf("%d: def.Path:%s def.NodeClass:%s\n", level, def.Path, def.NodeClass)

	var nodes []NodeDef
	if def.NodeClass == ua.NodeClassVariable {
		nodes = append(nodes, def)
	}

	browseChildren := func(refType uint32) error {
		refs, err := n.ReferencedNodesWithContext(ctx, refType, ua.BrowseDirectionForward, ua.NodeClassAll, true)
		if err != nil {
			return errors.Errorf("References: %d: %s", refType, err)
		}
		// fmt.Printf("found %d child refs\n", len(refs))
		for _, rn := range refs {
			children, err := browse(ctx, rn, def.Path, level+1)
			if err != nil {
				return errors.Errorf("browse children: %s", err)
			}
			nodes = append(nodes, children...)
		}
		return nil
	}

	if err := browseChildren(id.HasComponent); err != nil {
		return nil, err
	}
	if err := browseChildren(id.Organizes); err != nil {
		return nil, err
	}
	if err := browseChildren(id.HasProperty); err != nil {
		return nil, err
	}
	return nodes, nil
}
