package message

type Updates interface {
	SendUUID()
	SendOpcNodes()
}

type Update struct {
	SendUUID     *SendUUID
	SendOpcNodes *SendOpcNodes
}
