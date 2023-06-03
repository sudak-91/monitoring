package opcuaservice

import (
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
)

func (opc *OPCUAService) GetNodeDescription(node *opcua.Node) (DataType string, Description string, Err error) {
	attr, err := node.Attributes(ua.AttributeIDDataType, ua.AttributeIDDescription)
	if err != nil {
		Err = err
		return
	}
	if attr[0].Status == ua.StatusOK {
		switch v := attr[0].Value.NodeID().IntID(); v {
		case id.DateTime:
			DataType = "time.Time"
		case id.Boolean:
			DataType = "bool"
		case id.SByte:
			DataType = "int8"
		case id.Int16:
			DataType = "int16"
		case id.Int32:
			DataType = "int32"
		case id.Byte:
			DataType = "byte"
		case id.UInt16:
			DataType = "uint16"
		case id.UInt32:
			DataType = "uint32"
		case id.UtcTime:
			DataType = "time.Time"
		case id.String:
			DataType = "string"
		case id.Float:
			DataType = "float32"
		case id.Double:
			DataType = "float64"
		default:
			DataType = attr[0].Value.NodeID().String()
		}

	} else {

		DataType = "error"
	}
	if attr[1].Status == ua.StatusOK {
		Description = attr[1].Value.String()
	} else {
		Description = "error"
	}
	return
}
