package msusecase

import (
	update "github.com/sudak-91/monitoring/pkg/message/update"
)

func (u *MSUseCase) OpcNodesUpdate(update *update.OPCNodes) {
	u.messageChannel <- update
}

func (u *MSUseCase) OpcSubNodesupdate(update *update.SubNodes) {
	u.messageChannel <- update
}
