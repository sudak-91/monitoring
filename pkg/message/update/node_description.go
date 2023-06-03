package message

type NodeDescription struct {
	DataType    string
	Description string
}

func NewNodeDescriptionUpdate(dataType string, description string) Update {
	var u Update
	var n NodeDescription
	n.DataType = dataType
	n.Description = description
	u.NodeDescription = &n
	return u
}
