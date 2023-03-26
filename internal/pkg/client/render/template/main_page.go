package template

import (
	"github.com/sudak-91/wasmhtml/element"
)

type MainPage struct {
	MainDiv *element.Div
}

func NewMianPage() *MainPage {
	body := element.GetBody()
	body.SetStyle(`position: relative;
	width: 1920px;
	height: 1080px;
	background: #D2D2D2;`)
	var m MainPage
	m.MainDiv = body.AddDiv()
	m.MainDiv.SetID("maindiv")
	return &m

}

func (m *MainPage) Render() {
	element.GetBody().Generate()
}
