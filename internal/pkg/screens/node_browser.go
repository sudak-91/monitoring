package screens

import (
	"context"
	"log"
	"strconv"
	"syscall/js"

	"github.com/sudak-91/monitoring/internal/pkg/screens/internal/unit"
	wsusecase "github.com/sudak-91/monitoring/internal/pkg/ws/use_case"
	message "github.com/sudak-91/monitoring/pkg/message"
	command "github.com/sudak-91/monitoring/pkg/message/command"
	update "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/wasmhtml"
	"github.com/sudak-91/wasmhtml/element"
)

type NodeBrowser struct {
	CommonScreen
	ws       wsusecase.SocketUseCase
	Parent   *element.Div
	NodeTree *element.Div
	NodeInfo *element.Div
}

func NewNodeBrowser(renderChan chan<- any, screenChan <-chan any, parent *element.Div, ws wsusecase.SocketUseCase) Screens {
	var screen NodeBrowser
	screen.Parent = parent
	screen.ws = ws
	screen.DOMModel = make(map[string]any)
	screen.screenChan = screenChan
	screen.renderChan = renderChan
	screen.NodeTree = parent.AddDiv()
	screen.NodeTree.AddClass("nodeview")
	screen.NodeInfo = parent.AddDiv()
	screen.NodeInfo.AddClass("nodeinfo")
	return &screen
}

func (n *NodeBrowser) Render(ctx context.Context) {
	n.Parent.Generate()
	go n.update(ctx)
}

func (n *NodeBrowser) update(ctx context.Context) {
	command := command.GetOpcUaNodeCommand()
	data, err := message.EncodeData(command)
	if err != nil {
		log.Println(err.Error())
	}
	err = n.ws.Request(context.TODO(), data)
	if err != nil {
		log.Println(err.Error())
	}

mainloop:
	for {
		select {
		case data := <-n.screenChan:
			log.Println(data)
		case <-ctx.Done():
			log.Println("Update complete")
			break mainloop
		}
	}
}

func (n *NodeBrowser) Update(data any) {
	switch upd := data.(type) {
	case *update.OPCNodes:
		var (
			organizeNodeDiv  *element.Div
			componentNodeDiv *element.Div
			propertyNodeDiv  *element.Div
		)
		log.Println("[MainScreen]Update|OPCNodes")
		if len(upd.Nodes.OrganizesNode) != 0 {
			organizeNodeDiv = n.NodeTree.AddDiv()
			organizeNodeDiv.SetID("organizeNode")
			n.DOMModel[organizeNodeDiv.GetID()] = organizeNodeDiv
		}

		if len(upd.Nodes.ComponentNode) != 0 {
			componentNodeDiv = n.NodeTree.AddDiv()
			componentNodeDiv.SetID("componentNode")
			n.DOMModel[componentNodeDiv.GetID()] = componentNodeDiv
		}
		if len(upd.Nodes.PropertyNode) != 0 {
			propertyNodeDiv = n.NodeTree.AddDiv()
			propertyNodeDiv.SetID("propertyNode")
			n.DOMModel[propertyNodeDiv.GetID()] = propertyNodeDiv
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
			elem.AddEventListener(js.FuncOf(n.GetValue))
			n.DOMModel[v.Name] = elem.GetParentDiv()
		}
		for _, v := range upd.Nodes.ComponentNode {
			l := componentNodeDiv.AddLi()
			l.SetID(v.Name)
			l.SetInnerHtml(v.Name)
			wasmhtml.SetAttribute(l.Object, "opcid", v.IID)
			wasmhtml.SetAttribute(l.Object, "opcns", v.Namespace)
			wasmhtml.AddClickEventListenr(l.Object, js.FuncOf(n.GetValue))
			n.DOMModel[v.Name] = l
		}
		for _, v := range upd.Nodes.PropertyNode {
			l := propertyNodeDiv.AddLi()
			l.SetID(v.Name)
			l.SetInnerHtml(v.Name)
			wasmhtml.SetAttribute(l.Object, "opcid", v.IID)
			wasmhtml.SetAttribute(l.Object, "opcns", v.Namespace)
			wasmhtml.AddClickEventListenr(l.Object, js.FuncOf(n.GetValue))
			n.DOMModel[v.Name] = l
		}
		n.Parent.Generate()
		return
	case *update.SubNodes:
		log.Println("[MainScreen]|Update|SubNodes")
		log.Printf("[MainScreen] Parent is: %s", upd.Parent)
		elem := n.DOMModel[upd.Parent]

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
			node.AddEventListener(js.FuncOf(n.GetValue))
			n.DOMModel[v.Name] = node.GetParentDiv()
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
			node.AddEventListener(js.FuncOf(n.GetValue))
			n.DOMModel[v.Name] = node.GetParentDiv()
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
			node.AddEventListener(js.FuncOf(n.GetValue))
			n.DOMModel[v.Name] = node.GetParentDiv()
		}
		parent.Generate()
		return

	}
}

func (n *NodeBrowser) GetValue(this js.Value, args []js.Value) any {
	parent := this.Get("parentElement")
	idRaw := parent.Call("getAttribute", "opcid")
	nodeID, err := strconv.ParseUint(idRaw.String(), 10, 32)
	if err != nil {
		return err
	}
	namespcaeRaw := parent.Call("getAttribute", "opcns").String()
	namespace, err := strconv.ParseUint(namespcaeRaw, 10, 16)
	if err != nil {
		return err
	}
	sidRaw := parent.Call("getAttribute", "opcsid").String()

	parentID := parent.Get("id")
	log.Printf("[Function]|Parent is: %s", parentID.String())
	cmd := command.GetSubNodeCommande(parentID.String(), uint32(nodeID), uint16(namespace), sidRaw)
	data, err := message.EncodeData(cmd)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	err = n.ws.Request(context.TODO(), data)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	log.Printf("namespace = %d, id = %d", cmd.GetSubNode.Namespace, cmd.GetSubNode.IID)
	return nil
}
