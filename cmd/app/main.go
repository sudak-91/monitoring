package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
	"github.com/pkg/errors"
	"github.com/sudak-91/monitoring/pkg/webserver"
)

type NodeDef struct {
	NodeID      *ua.NodeID
	NodeClass   ua.NodeClass
	BrowseName  string
	Description string
	AccessLevel ua.AccessLevelType
	Path        string
	DataType    string
	Writable    bool
	Unit        string
	Scale       string
	Min         string
	Max         string
}
type OPCUAObjectData struct {
	Name []string
}

func main() {
	ctx1 := context.Background()
	Server := webserver.NewServer(ctx1)
	go Server.Start()
	var (
		endpoint = "opc.tcp://192.168.1.225:4840"
		ctx      = context.Background()
		c        = opcua.NewClient(endpoint)

		Names []string
	)

	if err := c.Connect(ctx); err != nil {
		panic(err)
	}
	defer c.CloseWithContext(ctx)

	uid, err := ua.ParseNodeID("ns=0;i=84")
	if err != nil {
		panic(err)
	}

	node := c.Node(uid)
	nodesList, err := node.ReferencedNodesWithContext(context.Background(), id.Organizes, ua.BrowseDirectionForward, ua.NodeClassAll, true)
	if err != nil {
		panic(err)
	}
	for _, v := range nodesList {
		name, err := v.BrowseName()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		Names = append(Names, name.Name)
		fmt.Println(name.Name)
		fmt.Printf("Roots node list is %v\n", v.ID)
		subnodeLists, err := v.ReferencedNodesWithContext(context.Background(), id.Organizes, ua.BrowseDirectionForward, ua.NodeClassAll, true)
		if err != nil {
			continue
		}
	subnodeloop:
		for _, k := range subnodeLists {
			name, err := k.BrowseName()
			if err != nil {
				continue subnodeloop
			}
			fmt.Println(name.Name)
		}
	}

	/*mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./template/index.html")
		if err != nil {
			log.Println(err.Error())
		}
		OPCUAObjects := OPCUAObjectData{
			Name: Names,
		}
		tmpl.Execute(w, OPCUAObjects)
	})

	if err = http.ListenAndServe("localhost:8000", mux); err != nil {
		log.Fatal(err.Error())
	}
	Endpoints, err := opcua.GetEndpoints(context.Background(), endpoint)
	if err != nil {
		panic(err)
	}
	for _, ep := range Endpoints {
		fmt.Println(ep.EndpointURL)

	}
	endpesp, err := c.GetEndpointsWithContext(ctx)
	if err != nil {
		panic(err)
	}
	for _, points := range endpesp.Endpoints {
		fmt.Printf("%+v for endpoints browser\n", points.EndpointURL)
	}
	arr := c.Namespaces()
	for _, ar := range arr {
		fmt.Println(ar)
	}

	DataValue, err := node.AttributesWithContext(ctx, ua.AttributeIDBrowseName, ua.AttributeIDNodeClass)
	for _, v := range DataValue {
		fmt.Printf("Attribute is %v\n", v)
	}

	refNodes, err := node.ReferencesWithContext(ctx, id.BaseObjectType, ua.BrowseDirectionForward, ua.NodeClassObjectType, true)

	if err != nil {
		panic(err)
	}
	for _, v := range refNodes {
		fmt.Printf("Reference is %v\n", v)
	}
	childNode, err := node.ChildrenWithContext(ctx, id.BaseObjectType, ua.NodeClassObjectType)
	if err != nil {
		panic(err)
	}
	for _, v := range childNode {
		fmt.Printf("Child node is %v", v)
	}

	/*nd, err := browse(ctx, c.Node(uid), "", 0)
	if err != nil {
		panic(err)
	}
	for _, n := range nd {
		fmt.Println(n.BrowseName)
	}
	*/
	l := make(chan bool)
	<-l
}

func (n NodeDef) Records() []string {
	return []string{n.BrowseName, n.DataType, n.NodeID.String(), n.Unit, n.Scale, n.Min, n.Max, strconv.FormatBool(n.Writable), n.Description}
}

func join(a, b string) string {
	if a == "" {
		return b
	}
	return a + "." + b
}

func browse(ctx context.Context, n *opcua.Node, path string, level int) ([]NodeDef, error) {
	// fmt.Printf("node:%s path:%q level:%d\n", n, path, level)
	if level > 10 {
		return nil, nil
	}

	attrs, err := n.AttributesWithContext(ctx, ua.AttributeIDNodeClass, ua.AttributeIDBrowseName, ua.AttributeIDDescription, ua.AttributeIDAccessLevel, ua.AttributeIDDataType)
	if err != nil {
		return nil, err
	}

	var def = NodeDef{
		NodeID: n.ID,
	}

	switch err := attrs[0].Status; err {
	case ua.StatusOK:
		def.NodeClass = ua.NodeClass(attrs[0].Value.Int())
	default:
		return nil, err
	}

	switch err := attrs[1].Status; err {
	case ua.StatusOK:
		def.BrowseName = attrs[1].Value.String()
	default:
		return nil, err
	}

	switch err := attrs[2].Status; err {
	case ua.StatusOK:
		def.Description = attrs[2].Value.String()
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return nil, err
	}

	switch err := attrs[3].Status; err {
	case ua.StatusOK:
		def.AccessLevel = ua.AccessLevelType(attrs[3].Value.Int())
		def.Writable = def.AccessLevel&ua.AccessLevelTypeCurrentWrite == ua.AccessLevelTypeCurrentWrite
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return nil, err
	}

	switch err := attrs[4].Status; err {
	case ua.StatusOK:
		switch v := attrs[4].Value.NodeID().IntID(); v {
		case id.DateTime:
			def.DataType = "time.Time"
		case id.Boolean:
			def.DataType = "bool"
		case id.SByte:
			def.DataType = "int8"
		case id.Int16:
			def.DataType = "int16"
		case id.Int32:
			def.DataType = "int32"
		case id.Byte:
			def.DataType = "byte"
		case id.UInt16:
			def.DataType = "uint16"
		case id.UInt32:
			def.DataType = "uint32"
		case id.UtcTime:
			def.DataType = "time.Time"
		case id.String:
			def.DataType = "string"
		case id.Float:
			def.DataType = "float32"
		case id.Double:
			def.DataType = "float64"
		default:
			def.DataType = attrs[4].Value.NodeID().String()
		}
	case ua.StatusBadAttributeIDInvalid:
		// ignore
	default:
		return nil, err
	}

	def.Path = join(path, def.BrowseName)
	// fmt.Printf("%d: def.Path:%s def.NodeClass:%s\n", level, def.Path, def.NodeClass)

	var nodes []NodeDef
	if def.NodeClass == ua.NodeClassVariable {
		nodes = append(nodes, def)
	}

	browseChildren := func(refType uint32) error {
		refs, err := n.ReferencedNodesWithContext(ctx, refType, ua.BrowseDirectionForward, ua.NodeClassAll, true)
		if err != nil {
			return errors.Errorf("References: %d: %s", refType, err)
		}
		// fmt.Printf("found %d child refs\n", len(refs))
		for _, rn := range refs {
			children, err := browse(ctx, rn, def.Path, level+1)
			if err != nil {
				return errors.Errorf("browse children: %s", err)
			}
			nodes = append(nodes, children...)
		}
		return nil
	}

	if err := browseChildren(id.HasComponent); err != nil {
		return nil, err
	}
	if err := browseChildren(id.Organizes); err != nil {
		return nil, err
	}
	if err := browseChildren(id.HasProperty); err != nil {
		return nil, err
	}
	return nodes, nil
}
