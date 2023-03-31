package message

import opcuaservice "github.com/sudak-91/monitoring/pkg/server/opcua_service"

type SendOpcNodes struct {
	Nodes *opcuaservice.OPCUAObjectData
}

func NewSendOpcNodes(nodes *opcuaservice.OPCUAObjectData) *SendOpcNodes {
	var s SendOpcNodes
	s.Nodes = nodes
	return &s
}
