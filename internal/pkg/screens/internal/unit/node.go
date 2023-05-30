package unit

import (
	"syscall/js"

	"github.com/sudak-91/wasmhtml"
	"github.com/sudak-91/wasmhtml/element"
)

type NodeUnit struct {
	unit      *element.Div
	actiondiv *element.Div
	titile    *element.Div
}

func NewNodeUnit(parent *element.Div) *NodeUnit {
	var u NodeUnit
	u.unit = parent.AddDiv()
	u.unit.AddClass("node")
	u.unit.SetStyle(`
	white-space:nowrap
	`)
	u.actiondiv = u.unit.AddDiv()
	u.actiondiv.AddClass("action")
	u.actiondiv.SetInnerHtml("+")
	u.actiondiv.SetStyle(`
	display:inline-block
	`)
	u.titile = u.unit.AddDiv()
	u.titile.AddClass("title")
	u.titile.SetStyle(`
	display:inline-block
	`)
	return &u
}

func (n *NodeUnit) AddTitle(title string) {
	n.titile.SetInnerHtml(title)
}

func (n *NodeUnit) SetAttributes(id string, opcns uint16, opcid uint32, opcsid string) {
	n.unit.SetID(id)
	wasmhtml.SetAttribute(n.actiondiv.Object, "opcns", opcns)
	wasmhtml.SetAttribute(n.actiondiv.Object, "opcid", opcid)
	wasmhtml.SetAttribute(n.actiondiv.Object, "opcsid", opcsid)

}
func (n *NodeUnit) AddEventListener(fun js.Func) {
	wasmhtml.AddClickEventListenr(n.actiondiv.Object, fun)
}

func (n *NodeUnit) GetParentDiv() *element.Div {
	return n.unit
}

func (n *NodeUnit) AddClass(class string) {
	n.unit.AddClass(class)
}
