package client

import (
	"context"

	"github.com/google/uuid"
	update "github.com/sudak-91/monitoring/pkg/message/update"
)

type ChangeUUID struct {
	OldUUID uuid.UUID
	NewUUID uuid.UUID
}

type GetOpcUaNodeTransfer struct {
	Namespace    uint16
	IID          uint32
	SID          string
	ResponseChan chan update.OPCNode
	Cancel       context.CancelFunc
}

type NodeDescriptionTransfer struct {
	DataType    string
	Description string
}

type GetOpcUaNodeDescriptionTransfer struct {
	Namespace    uint16
	SID          string
	ResponseChan chan NodeDescriptionTransfer
	Cancel       context.CancelFunc
}

func NewGetOpcUaNodeTransfer(ns uint16, iid uint32, sid string, respChan chan update.OPCNode, ctx context.Context) (GetOpcUaNodeTransfer, context.Context) {
	var transfer GetOpcUaNodeTransfer
	transfer.Namespace = ns
	transfer.IID = iid
	transfer.SID = sid
	transfer.ResponseChan = respChan
	trnsferCtx, cancel := context.WithCancel(ctx)
	transfer.Cancel = cancel
	return transfer, trnsferCtx
}

func NewGetOpcUaNodeDescriptionTransfer(ns uint16, sid string, respChan chan NodeDescriptionTransfer, ctx context.Context) (GetOpcUaNodeDescriptionTransfer, context.Context) {
	var transfer GetOpcUaNodeDescriptionTransfer
	transfer.Namespace = ns
	transfer.SID = sid
	transfer.ResponseChan = respChan
	transferCtx, cancel := context.WithCancel(ctx)
	transfer.Cancel = cancel
	return transfer, transferCtx
}
