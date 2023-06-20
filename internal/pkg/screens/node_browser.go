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
			organizeNodeDiv *element.Div
			// componentNodeDiv *element.Div
			// propertyNodeDiv  *element.Div
		)
		log.Println("[MainScreen]Update|OPCNodes")
		if len(upd.Nodes.OrganizesNode) != 0 {
			organizeNodeDiv = n.createMainNodeDiv("organizeNode")
		}

		// if len(upd.Nodes.ComponentNode) != 0 {
		// 	componentNodeDiv = n.createMainNodeDiv("componentNode")
		// }
		// if len(upd.Nodes.PropertyNode) != 0 {
		// 	propertyNodeDiv = n.createMainNodeDiv("propertyNode")
		// }

		n.createNodeDiv(organizeNodeDiv, upd.Nodes)

		// for _, v := range upd.Nodes.ComponentNode {
		// 	n.createNodeUnit(componentNodeDiv, v)
		// }
		// for _, v := range upd.Nodes.PropertyNode {
		// 	n.createNodeUnit(propertyNodeDiv, v)
		// }

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

		parent.Generate()
		return
	case *update.NodeDescription:
		log.Printf("[NodeDescription]:%s,%s", upd.DataType, upd.Description)

	}
}

func (n *NodeBrowser) createNodeDiv(parent *element.Div, nodes *update.OPCNode) {
	for _, v := range nodes.OrganizesNode {
		n.createNodeUnit(parent, v)
		parDiv, _ := n.DOMModel[v.Name].(*element.Div)
		n.createNodeDiv(parDiv, &v.ChildNode)
	}
	for _, v := range nodes.ComponentNode {
		n.createNodeUnit(parent, v)
		parDiv, _ := n.DOMModel[v.Name].(*element.Div)
		n.createNodeDiv(parDiv, &v.ChildNode)
	}
	for _, v := range nodes.PropertyNode {
		n.createNodeUnit(parent, v)
		parDiv, _ := n.DOMModel[v.Name].(*element.Div)
		n.createNodeDiv(parDiv, &v.ChildNode)
	}
	parent.Generate()
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

func (n *NodeBrowser) GetNodeDescription(this js.Value, args []js.Value) any {
	log.Println("[GetNodeDescription]")
	parent := this.Get("parentElement")
	namespaceRaw := parent.Call("getAttribute", "opcns").String()
	namespace, err := strconv.ParseUint(namespaceRaw, 10, 16)
	if err != nil {
		return err
	}
	sidRaw := parent.Call("getAttribute", "opcsid").String()

	cmd := command.GetNodeDescriptionCommand(uint16(namespace), sidRaw)
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
	return nil
}

func (n *NodeBrowser) createMainNodeDiv(id string) *element.Div {
	mainNode := n.NodeTree.AddDiv()
	mainNode.SetID(id)
	n.DOMModel[mainNode.GetID()] = mainNode
	return mainNode
}

func (n *NodeBrowser) createNodeUnit(parentNode *element.Div, v update.NodeDef) {
	folderFunc := js.FuncOf(n.GetValue)
	nodeFunc := js.FuncOf(n.GetNodeDescription)
	elem := unit.NewNodeUnit(parentNode, v, folderFunc, nodeFunc)

	n.DOMModel[v.Name] = elem.GetParentDiv()

}
