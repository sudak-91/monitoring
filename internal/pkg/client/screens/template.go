package screens

import "context"

type Renderer interface {
	Render(ctx context.Context)
}

// Template struct for web page
type CommonScreen struct {
	renderChan chan<- interface{}
	screenChan <-chan interface{}
}
