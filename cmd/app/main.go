package main

import (
	"context"
	"log"

	"github.com/sudak-91/monitoring/pkg/server"
	opcuaservice "github.com/sudak-91/monitoring/pkg/server/opcua_service"
	"github.com/sudak-91/monitoring/pkg/server/updateservice"
	webservice "github.com/sudak-91/monitoring/pkg/server/web_service"
)

func main() {
	log.Println("Star Monitoring Server")
	MainCTX := context.Background()
	server := server.NewServer(MainCTX)
	updateToClientChan := make(chan any, 5)
	updateToOpcUaChan := make(chan any, 6)
	opcuaChan := make(chan any, 5)
	clientChan := make(chan any, 5)
	log.Println("Create Web Service")
	webService := webservice.NewWebService(MainCTX, updateToClientChan, server, clientChan)
	go webService.Run()
	log.Println("Create Update Service")
	updateService := updateservice.NewUpdateService(MainCTX, clientChan, opcuaChan, updateToClientChan, updateToOpcUaChan, server)
	go updateService.Update()
	opcuaservice := opcuaservice.NewOpcUaService(MainCTX, opcuaChan, updateToOpcUaChan)
	go opcuaservice.StartOPCUA()

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
