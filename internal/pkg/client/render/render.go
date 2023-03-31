package render

import (
	"context"
	"log"

	"github.com/sudak-91/monitoring/internal/pkg/client/screens"
	message "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/wasmhtml"
)

type Render struct {
	ScreenList         map[string]interface{}
	renderChan         <-chan interface{}
	screenChan         chan<- interface{}
	messageServiceChan chan interface{}
	ctx                context.Context
}

type ActualScreen struct {
	Init           bool
	Screen         interface{}
	CancelFunction context.CancelFunc
}

func NewRender(ctx context.Context, renderChan <-chan interface{}, screenChan chan<- interface{}, messageServiceChan chan interface{}) *Render {
	var r Render
	r.ScreenList = make(map[string]interface{})
	r.ctx = ctx
	r.renderChan = renderChan
	r.screenChan = screenChan
	r.messageServiceChan = messageServiceChan
	return &r
}

func (r *Render) AddScreen(key string, screen screens.Renderer) {
	r.ScreenList[key] = screen
}

func (r *Render) Run() {
	k := &ActualScreen{}
	for {
		select {
		case <-r.ctx.Done():
		case screen := <-r.renderChan:
			if data, ok := screen.(string); ok {
				if k.Init {
					k.CancelFunction()
				}
				ctx, cancel := context.WithCancel(r.ctx)
				k.CancelFunction = cancel
				k.Init = true

				if render, ok := r.ScreenList[data].(screens.Renderer); ok {
					render.Render(ctx)
					k.Screen = r.ScreenList[data]
				}
				continue
			}
			log.Println("Fail")
		case data := <-r.messageServiceChan:
			log.Println("MessageService chan")
			if NewData, ok := data.(*message.SendOpcNodes); ok {
				log.Println("SendOPC NODE")
				if screen, ok := k.Screen.(*screens.MainScreen); ok {
					for _, v := range NewData.Nodes.Nodes {
						l := screen.MainDiv.AddDiv()
						l.SetID(v.Name)
						l.SetInnerHtml(v.Name)
						wasmhtml.SetAttribute(l.Object, "opcid", string(v.ID))
					}
					screen.MainDiv.Generate()
					return
				}
			}
		}
	}
}
