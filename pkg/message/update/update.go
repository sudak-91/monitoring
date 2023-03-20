package message

type Updates interface {
	SendUUID()
}

type Update struct {
	SendUUID *SendUUID
}
