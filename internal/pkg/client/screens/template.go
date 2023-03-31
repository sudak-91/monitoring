package screens

import (
	"context"
	"syscall/js"
)

type Renderer interface {
	Render(ctx context.Context)
}

// Template struct for web page
type CommonScreen struct {
	DOMModel   map[string]js.Value
	renderChan chan<- interface{}
	screenChan <-chan interface{}
}
