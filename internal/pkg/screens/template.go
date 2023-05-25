package screens

import (
	"context"
)

type Screens interface {
	Render(ctx context.Context)
	Update(data any)
}

// Template struct for web page
type CommonScreen struct {
	DOMModel   map[string]any
	renderChan chan<- interface{}
	screenChan <-chan interface{}
}
