package message

import (
	"bytes"
	"encoding/gob"

	commandMessage "github.com/sudak-91/monitoring/pkg/message/command"
	updateMessage "github.com/sudak-91/monitoring/pkg/message/update"
)

type DecodeEncode interface {
	commandMessage.Command | updateMessage.Update
}

func EncodeData(data any) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
func Decode[V DecodeEncode](decodeData []byte) (V, error) {
	var (
		decode V
	)
	reader := bytes.NewReader(decodeData)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&decode)
	if err != nil {
		return decode, err
	}
	return decode, nil
}
