package render

import "context"

type Screen interface {
	Render()
}

type Render struct {
	ScreenList map[string]Screen
	ScreenName chan string
	ctx        context.Context
}

func NewRender(ctx context.Context, screenNameChan chan string) *Render {
	var r Render
	r.ScreenList = make(map[string]Screen)
	r.ctx = ctx
	r.ScreenName = screenNameChan
	return &r
}

func (r *Render) AddScreen(key string, screen Screen) {
	r.ScreenList[key] = screen
}

func (r *Render) Run() {
	for {
		select {
		case <-r.ctx.Done():
		case screen := <-r.ScreenName:
			r.ScreenList[screen].Render()
		}
	}
}
