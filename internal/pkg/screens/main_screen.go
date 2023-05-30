package screens

import (
	"context"
	"log"
	"strconv"
	"syscall/js"

	"github.com/sudak-91/monitoring/internal/pkg/screens/internal/unit"
	wsusecase "github.com/sudak-91/monitoring/internal/pkg/ws/use_case"
	"github.com/sudak-91/monitoring/pkg/message"
	command "github.com/sudak-91/monitoring/pkg/message/command"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/wasmhtml"
	"github.com/sudak-91/wasmhtml/element"
)

type MainScreen struct {
	CommonScreen
	ws          wsusecase.SocketUseCase
	body        *element.Body
	MainDiv     *element.Div
	Navigate    *element.Div
	Browse      *element.Div
	NodeTree    *element.Div
	NodeInfo    *element.Div
	StatusDiv   *element.Div
	OPCUAStatus *element.Div
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

	m.NodeTree = m.Browse.AddDiv()
	m.NodeTree.AddClass("nodeview")

	m.NodeInfo = m.Browse.AddDiv()
	m.NodeInfo.AddClass("nodeinfo")
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
	m.DOMModel["body"] = m.body
	m.DOMModel[m.MainDiv.GetID()] = m.MainDiv
	m.DOMModel[m.StatusDiv.GetID()] = m.StatusDiv
	m.DOMModel[m.OPCUAStatus.GetID()] = m.OPCUAStatus

	return &m
}

func (m *MainScreen) Render(ctx context.Context) {

	m.body.Generate()
	go m.update(ctx)
}

func (m *MainScreen) update(ctx context.Context) {
	command := command.GetOpcUaNodeCommand()
	data, err := message.EncodeData(command)
	if err != nil {
		log.Println(err.Error())
	}
	err = m.ws.Request(context.TODO(), data)
	if err != nil {
		log.Println(err.Error())
	}

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

func (m *MainScreen) Update(data any) {
	switch upd := data.(type) {
	case *update.OPCNodes:
		m.MainDiv.Child = nil
		var (
			organizeNodeDiv  *element.Div
			componentNodeDiv *element.Div
			propertyNodeDiv  *element.Div
		)
		log.Println("[MainScreen]Update|OPCNodes")
		if len(upd.Nodes.OrganizesNode) != 0 {
			organizeNodeDiv = m.NodeTree.AddDiv()
			organizeNodeDiv.SetID("organizeNode")
			m.DOMModel[organizeNodeDiv.GetID()] = organizeNodeDiv
		}

		if len(upd.Nodes.ComponentNode) != 0 {
			componentNodeDiv = m.NodeTree.AddDiv()
			componentNodeDiv.SetID("componentNode")
			m.DOMModel[componentNodeDiv.GetID()] = componentNodeDiv
		}
		if len(upd.Nodes.PropertyNode) != 0 {
			propertyNodeDiv = m.NodeTree.AddDiv()
			propertyNodeDiv.SetID("propertyNode")
			m.DOMModel[propertyNodeDiv.GetID()] = propertyNodeDiv
		}
		for _, v := range upd.Nodes.OrganizesNode {

			elem := unit.NewNodeUnit(organizeNodeDiv)
			switch v.NodeType {
			case 0:
				elem.AddClass("foldernode")
			default:
				elem.AddClass("datanode")
			}
			elem.AddTitle(v.Name)
			elem.SetAttributes(v.Name, v.Namespace, v.IID, "")
			elem.AddEventListener(js.FuncOf(m.GetValue))
			m.DOMModel[v.Name] = elem.GetParentDiv()
		}
		for _, v := range upd.Nodes.ComponentNode {
			l := componentNodeDiv.AddLi()
			l.SetID(v.Name)
			l.SetInnerHtml(v.Name)
			wasmhtml.SetAttribute(l.Object, "opcid", v.IID)
			wasmhtml.SetAttribute(l.Object, "opcns", v.Namespace)
			wasmhtml.AddClickEventListenr(l.Object, js.FuncOf(m.GetValue))
			m.DOMModel[v.Name] = l
		}
		for _, v := range upd.Nodes.PropertyNode {
			l := propertyNodeDiv.AddLi()
			l.SetID(v.Name)
			l.SetInnerHtml(v.Name)
			wasmhtml.SetAttribute(l.Object, "opcid", v.IID)
			wasmhtml.SetAttribute(l.Object, "opcns", v.Namespace)
			wasmhtml.AddClickEventListenr(l.Object, js.FuncOf(m.GetValue))
			m.DOMModel[v.Name] = l
		}
		m.MainDiv.Generate()
		return
	case *update.SubNodes:
		log.Println("[MainScreen]|Update|SubNodes")
		log.Printf("[MainScreen] Parent is: %s", upd.Parent)
		elem := m.DOMModel[upd.Parent]

		parent, ok := elem.(*element.Div)
		//parent.RemoveAllChild()
		if !ok {
			log.Println("parent element is not div")
			return
		}
		for _, v := range upd.Nodes.OrganizesNode {

			node := unit.NewNodeUnit(parent)
			log.Println(v.NodeType)
			switch v.NodeType {
			case 0:
				node.AddClass("foldernode")
			default:
				node.AddClass("datanode")
			}
			node.AddTitle(v.Name)
			node.SetAttributes(v.Name, v.Namespace, v.IID, v.SID)
			node.AddEventListener(js.FuncOf(m.GetValue))
			m.DOMModel[v.Name] = node.GetParentDiv()
		}
		for _, v := range upd.Nodes.ComponentNode {
			node := unit.NewNodeUnit(parent)
			log.Println(v.NodeType)
			switch v.NodeType {
			case 0:
				node.AddClass("foldernode")
			default:
				node.AddClass("datanode")
			}
			node.AddTitle(v.Name)
			node.SetAttributes(v.Name, v.Namespace, v.IID, v.SID)
			node.AddEventListener(js.FuncOf(m.GetValue))
			m.DOMModel[v.Name] = node.GetParentDiv()
		}
		for _, v := range upd.Nodes.PropertyNode {
			node := unit.NewNodeUnit(parent)
			log.Println(v.NodeType)
			switch v.NodeType {
			case 0:
				node.AddClass("foldernode")
			default:
				node.AddClass("datanode")
			}
			node.AddTitle(v.Name)
			node.SetAttributes(v.Name, v.Namespace, v.IID, v.SID)
			node.AddEventListener(js.FuncOf(m.GetValue))
			m.DOMModel[v.Name] = node.GetParentDiv()
		}
		parent.Generate()
		return

	}
}

func (m *MainScreen) GetValue(this js.Value, args []js.Value) any {
	idRaw := this.Call("getAttribute", "opcid")
	nodeID, err := strconv.ParseUint(idRaw.String(), 10, 32)
	if err != nil {
		return err
	}
	namespcaeRaw := this.Call("getAttribute", "opcns").String()
	namespace, err := strconv.ParseUint(namespcaeRaw, 10, 16)
	if err != nil {
		return err
	}
	sidRaw := this.Call("getAttribute", "opcsid").String()
	parent := this.Get("parentElement")
	log.Println(parent)
	parentID := parent.Get("id")
	log.Printf("[Function]|Parent is: %s", parentID.String())
	cmd := command.GetSubNodeCommande(parentID.String(), uint32(nodeID), uint16(namespace), sidRaw)
	data, err := message.EncodeData(cmd)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	err = m.ws.Request(context.TODO(), data)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	log.Printf("namespace = %d, id = %d", cmd.GetSubNode.Namespace, cmd.GetSubNode.IID)
	return nil
}
