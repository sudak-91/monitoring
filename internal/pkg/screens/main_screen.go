package screens

import (
	"context"

	wsusecase "github.com/sudak-91/monitoring/internal/pkg/ws/use_case"
	"github.com/sudak-91/wasmhtml"
	"github.com/sudak-91/wasmhtml/element"
)

type MainScreen struct {
	CommonScreen
	ws            wsusecase.SocketUseCase
	ActualBrowser Screens
	body          *element.Body
	MainDiv       *element.Div
	Navigate      *element.Div
	Browse        *element.Div
}

func NewMainScreen(renderChan chan<- interface{}, screenChan <-chan interface{}, body *element.Body, ws wsusecase.SocketUseCase) Screens {
	var m MainScreen
	head := wasmhtml.Document.Get("head")
	style := wasmhtml.Document.Call("createElement", "style")
	head.Call("appendChild", style)
	style.Set("innerHTML", `
	.node .node{
		margin-left: 3em;
	  }

	  .main {
		display: flex;
		flex-direction: row;
		align-items: flex-start;
		padding: 0px 10px;
		gap: 10px;
		position: relative;
		width: 1920px;
		height: 1080px;
	  }

.navigate{
display: flex;
flex-direction: row;
align-items: flex-start;
padding: 10px;
gap: 10px;

width:300px;
height: 1080px;
background-color: green;
}

.browse{
	display: flex;
	flex-direction: row;
	align-items: flex-start;
	padding: 10px;
	gap: 10px;
	
	width: minmax(600,auto);
	height: 1080px;
	
	
	}

	.nodeview{
		display: flex;
		flex-direction: row;
		align-items: flex-start;
		padding: 10px;
		gap: 10px;
		
		width: 1260px;
		height: 1060px;
		
		
		}

		.nodeinfo{
			display: flex;
			flex-direction: row;
			align-items: flex-start;
			padding: 10px;
			gap: 10px;
			
			width: 300px;
			height: 1060px;
			  background-color: gray;
			}
	.foldernode{
		background-color: blue;
	}
	.datanode{
		background-color: yellow;
	}
	`)
	m.DOMModel = make(map[string]any)
	m.ws = ws
	m.screenChan = screenChan
	m.renderChan = renderChan
	m.body = body
	m.body.SetStyle(`position: relative;
	width: 1920px;
	height: 1080px;
	background: #D2D2D2;`)
	m.MainDiv = m.body.AddDiv()
	m.MainDiv.SetID("maindiv")
	m.MainDiv.AddClass("main")

	m.Navigate = m.MainDiv.AddDiv()
	m.Navigate.AddClass("navigate")

	m.Browse = m.MainDiv.AddDiv()
	m.Browse.AddClass("browse")

	m.DOMModel["body"] = m.body
	m.DOMModel[m.MainDiv.GetID()] = m.MainDiv

	return &m
}

func (m *MainScreen) Render(ctx context.Context) {
	m.body.Generate()
	nodeBrowse := NewNodeBrowser(m.renderChan, m.screenChan, m.Browse, m.ws)
	m.ActualBrowser = nodeBrowse
	nodeBrowse.Render(ctx)
}

func (m *MainScreen) Update(data any) {
	m.ActualBrowser.Update(data)
}
