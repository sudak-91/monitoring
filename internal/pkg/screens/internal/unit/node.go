package unit

import (
	"syscall/js"

	update "github.com/sudak-91/monitoring/pkg/message/update"
	"github.com/sudak-91/wasmhtml"
	"github.com/sudak-91/wasmhtml/element"
)

type NodeUnit struct {
	unit      *element.Div
	actiondiv *element.Div
	titile    *element.Div
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
