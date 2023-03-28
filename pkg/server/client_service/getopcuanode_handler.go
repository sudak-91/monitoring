package clientservice

import "github.com/sudak-91/monitoring/pkg/server/updateservice"

func (cs *ClientService) getOpcUaNodeHandle() error {
	var updateData updateservice.GetOpcUaNode
	updateData.Info = "get_opcua"
	cs.clientChan <- updateData
	return nil
}
