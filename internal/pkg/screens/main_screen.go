package screens

import (
	"context"
	"log"
	"strconv"
	"syscall/js"

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
	StatusDiv   *element.Div
	OPCUAStatus *element.Div
}

func NewMainScreen(renderChan chan<- interface{}, screenChan <-chan interface{}, body *element.Body, ws wsusecase.SocketUseCase) Screens {
	var m MainScreen
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
			organizeNodeDiv = m.MainDiv.AddDiv()
			organizeNodeDiv.SetID("organizeNode")
			m.DOMModel[organizeNodeDiv.GetID()] = organizeNodeDiv
		}
		if len(upd.Nodes.ComponentNode) != 0 {
			componentNodeDiv = m.MainDiv.AddDiv()
			componentNodeDiv.SetID("componentNode")
			m.DOMModel[componentNodeDiv.GetID()] = componentNodeDiv
		}
		if len(upd.Nodes.PropertyNode) != 0 {
			propertyNodeDiv = m.MainDiv.AddDiv()
			propertyNodeDiv.SetID("propertyNode")
			m.DOMModel[propertyNodeDiv.GetID()] = propertyNodeDiv
		}
		for _, v := range upd.Nodes.OrganizesNode {

			elem := organizeNodeDiv.AddLi()
			elem.SetID(v.Name)
			elem.SetInnerHtml(v.Name)
			wasmhtml.SetAttribute(elem.Object, "opcid", v.IID)
			wasmhtml.SetAttribute(elem.Object, "opcns", v.Namespace)
			wasmhtml.AddClickEventListenr(elem.Object, js.FuncOf(m.GetValue))
			m.DOMModel[v.Name] = elem
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

		parent, ok := elem.(*element.Li)
		parent.RemoveAllChild()
		list := parent.AddUl()
		if !ok {
			log.Println("parent element is not div")
			return
		}
		for _, v := range upd.Nodes.OrganizesNode {

			node := list.AddLi()
			node.SetID(v.Name)
			node.SetInnerHtml(v.Name)
			wasmhtml.SetAttribute(node.Object, "opcid", v.IID)
			wasmhtml.SetAttribute(node.Object, "opcsid", v.SID)
			wasmhtml.SetAttribute(node.Object, "opcns", v.Namespace)
			wasmhtml.AddClickEventListenr(node.Object, js.FuncOf(m.GetValue))
			m.DOMModel[v.Name] = node
		}
		for _, v := range upd.Nodes.ComponentNode {
			node := list.AddLi()
			node.SetID(v.Name)
			node.SetInnerHtml(v.Name)
			wasmhtml.SetAttribute(node.Object, "opcid", v.IID)
			wasmhtml.SetAttribute(node.Object, "opcsid", v.SID)
			wasmhtml.SetAttribute(node.Object, "opcns", v.Namespace)
			wasmhtml.AddClickEventListenr(node.Object, js.FuncOf(m.GetValue))
			m.DOMModel[v.Name] = node
		}
		for _, v := range upd.Nodes.PropertyNode {
			node := list.AddLi()
			node.SetID(v.Name)
			node.SetInnerHtml(v.Name)
			wasmhtml.SetAttribute(node.Object, "opcid", v.IID)
			wasmhtml.SetAttribute(node.Object, "opcsid", v.SID)
			wasmhtml.SetAttribute(node.Object, "opcns", v.Namespace)
			wasmhtml.AddClickEventListenr(node.Object, js.FuncOf(m.GetValue))
			m.DOMModel[v.Name] = node
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
	cmd := command.GetSubNodeCommande(this.Get("id").String(), uint32(nodeID), uint16(namespace), sidRaw)
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
