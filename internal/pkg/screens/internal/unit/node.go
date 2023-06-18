package unit

import (
	"log"
	"strings"
	"syscall/js"

	update "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/wasmhtml"
	"github.com/sudak-91/wasmhtml/element"
)

type NodeUnit struct {
	unit      *element.Div
	actiondiv *element.Div
	titile    *element.Div
	checkbox  *element.Input
}

func NewNodeUnit(parent *element.Div, nodeDef update.NodeDef, folderFunc js.Func, dataFunc js.Func) *NodeUnit {
	var u NodeUnit
	u.unit = parent.AddDiv()
	u.unit.AddClass("node")
	u.unit.SetStyle(`
	white-space:nowrap
	`)
	u.actiondiv = u.unit.AddDiv()
	u.actiondiv.AddClass("action")
	switch nodeDef.NodeType {
	case 0:
		u.actiondiv.SetInnerHtml("+")
		u.unit.AddClass("foldernode")
		u.addEventListener(u.actiondiv.Object, folderFunc)

	default:
		u.actiondiv.SetInnerHtml(">")
		u.unit.AddClass("datanode")
		u.addEventListener(u.actiondiv.Object, dataFunc)
	}

	u.actiondiv.SetStyle(`
	display:inline-block
	`)
	u.titile = u.unit.AddDiv()
	u.titile.AddClass("title")
	u.titile.SetStyle(`
	display:inline-block
	`)
	u.addTitle(nodeDef.Name)
	u.setAttributes(nodeDef.Name, nodeDef.Namespace, nodeDef.IID, nodeDef.SID)
	u.checkbox = u.unit.AddInput()
	u.checkbox.SetType("checkbox")
	wasmhtml.AddClass(u.checkbox.Object, "nodeSelector")
	wasmhtml.AddClickEventListenr(u.checkbox.Object, js.FuncOf(EnableCheckBox))
	return &u
}

func (n *NodeUnit) addTitle(title string) {
	n.titile.SetInnerHtml(title)
}

func (n *NodeUnit) setAttributes(id string, opcns uint16, opcid uint32, opcsid string) {
	n.unit.SetID(id)
	wasmhtml.SetAttribute(n.unit.Object, "opcns", opcns)
	wasmhtml.SetAttribute(n.unit.Object, "opcid", opcid)
	wasmhtml.SetAttribute(n.unit.Object, "opcsid", opcsid)

}
func (n *NodeUnit) addEventListener(object js.Value, fun js.Func) {
	wasmhtml.AddClickEventListenr(object, fun)
}

func (n *NodeUnit) GetParentDiv() *element.Div {
	return n.unit
}

func (n *NodeUnit) AddClass(class string) {
	n.unit.AddClass(class)
}

func selectChildren(object js.Value) (js.Value, error) {
	childrenArr, err := wasmhtml.GetChildren(object)
	if err != nil {
		return js.Value{}, err
	}
	return childrenArr, nil
}

func enableCheckBox(object js.Value) any {
	children, err := selectChildren(object)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(children.Length())
	for i := 0; i < children.Length(); i++ {
		obj := children.Index(i)
		className := strings.Split(obj.Get("className").String(), " ")
		switch className[0] {
		case "nodeSelector":
			log.Println("Node Selector")
			wasmhtml.Set(obj, "checked", "true")
			wasmhtml.SetAttribute(obj, "checked", "true")
		case "node":
			enableCheckBox(obj)
		}
	}
	return nil
}

func EnableCheckBox(this js.Value, args []js.Value) any {
	parent, _ := wasmhtml.GetParent(this)
	enableCheckBox(parent)
	return nil
}
