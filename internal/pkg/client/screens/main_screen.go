package screens

import (
	"context"
	"log"

	"github.com/sudak-91/monitoring/internal/pkg/client/command"
	"github.com/sudak-91/wasmhtml"
	"github.com/sudak-91/wasmhtml/element"
)

type MainScreen struct {
	CommonScreen
	command     *command.Command
	body        *element.Body
	MainDiv     *element.Div
	StatusDiv   *element.Div
	OPCUAStatus *element.Div
}

func NewMainScreen(renderChan chan<- interface{}, screenChan <-chan interface{}, body *element.Body, command *command.Command) Renderer {
	var m MainScreen

	m.screenChan = screenChan
	m.renderChan = renderChan
	m.command = command
	m.body = body
	m.body.SetStyle(`position: relative;
	width: 1920px;
	height: 1080px;
	background: #D2D2D2;`)
	m.MainDiv = m.body.AddDiv()
	m.MainDiv.SetID("maindiv")
	m.MainDiv.SetStyle(`position: absolute;
	width: 1510px;
	height: 1030px;
	left: 0px;
	top: 0px;`)
	m.StatusDiv = m.body.AddDiv()
	m.StatusDiv.SetID("statusdiv")
	m.StatusDiv.SetStyle(`position: absolute;
	width: 1920px;
	height: 50px;
	left: 0px;
	top: 1030px;`)
	m.OPCUAStatus = m.StatusDiv.AddDiv()
	m.OPCUAStatus.SetID("opcuastatus")
	m.OPCUAStatus.SetStyle(`display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
	gap: 10px;
	position: absolute;
	width: 175px;
	height: 50px;	
	background: #8FD189;`)
	m.OPCUAStatus.SetInnerHtml("OPC UA")
	wasmhtml.SetAttribute(m.OPCUAStatus.Object, "onclick", "alert(\"Click\")")
	return &m
}

func (m *MainScreen) Render(ctx context.Context) {

	m.body.Generate()
	go m.update(ctx)
}

func (m *MainScreen) update(ctx context.Context) {
	err := m.command.GetOpcUaNode()
	if err != nil {
		log.Printf("Main screen update has error:%s\n", err.Error())
	}
mainloop:
	for {
		select {
		case data := <-m.screenChan:
			log.Println(data)
		case <-ctx.Done():
			log.Println("Update complete")
			break mainloop
		}
	}
	log.Println("Screen closed")
}
