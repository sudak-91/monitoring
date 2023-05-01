package message

func NewOPCSubNodeUpdate(parent string) SubNodes {

	return SubNodes{Parent: parent}
}

func (s *SubNodes) GetUpdate() Update {
	var Update Update
	Update.OPCSubNode = s
	return Update
}
