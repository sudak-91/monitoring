package render

import (
	"context"
	"log"

	"github.com/sudak-91/monitoring/internal/pkg/screens"
)

type Render struct {
	ScreenList         map[string]screens.Screens
	renderChan         <-chan interface{}
	screenChan         chan<- interface{}
	messageServiceChan chan interface{}
	ctx                context.Context
}

type ActualScreen struct {
	Init           bool
	Screen         screens.Screens
	CancelFunction context.CancelFunc
}

func NewRender(ctx context.Context, renderChan <-chan interface{}, screenChan chan<- interface{}, messageServiceChan chan interface{}) *Render {
	var r Render
	r.ScreenList = make(map[string]screens.Screens)
	r.ctx = ctx
	r.renderChan = renderChan
	r.screenChan = screenChan
	r.messageServiceChan = messageServiceChan
	return &r
}

func (r *Render) AddScreen(key string, screen screens.Screens) {
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

				if render, ok := r.ScreenList[data].(screens.Screens); ok {
					render.Render(ctx)
					k.Screen = r.ScreenList[data]
				}
				continue
			}
			log.Println("Fail")
		case data := <-r.messageServiceChan:
			log.Println("MessageService chan")
			go k.Screen.Update(data)

		}
	}
}
