package render

import (
	"context"
	"log"

	"github.com/sudak-91/monitoring/internal/pkg/client/screens"
)

type Render struct {
	ScreenList map[string]screens.Renderer
	renderChan <-chan interface{}
	screenChan chan<- interface{}
	ctx        context.Context
}

type ActualScreen struct {
	Init           bool
	CancelFunction context.CancelFunc
}

func NewRender(ctx context.Context, renderChan <-chan interface{}, screenChan chan<- interface{}) *Render {
	var r Render
	r.ScreenList = make(map[string]screens.Renderer)
	r.ctx = ctx
	r.renderChan = renderChan
	r.screenChan = screenChan
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
				r.ScreenList[data].Render(ctx)
				continue
			}
			log.Println("Fail")

		}
	}
}
