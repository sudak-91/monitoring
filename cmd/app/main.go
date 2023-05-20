package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/sudak-91/monitoring/pkg/server"
	"github.com/sudak-91/monitoring/pkg/server/clients"
	opcuaservice "github.com/sudak-91/monitoring/pkg/server/opcua_service"
)

func main() {
	log.Println("Star Monitoring Server")
	viper.SetConfigName("service")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Println("Config read")
	debugMode := viper.GetBool("Debug")
	if debugMode {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}
	MainCTX := context.Background()
	ClientList := clients.NewClientList()
	updateToClientChan := make(chan any, 5)
	updateToOpcUaChan := make(chan any, 6)
	opcuaChan := make(chan any, 5)
	clientChan := make(chan any, 5)
	log.Println("Create Web Service")
	webService := server.NewWebService(MainCTX, updateToClientChan, clientChan, ClientList)
	log.Println("Create Update Service")
	//updateService := updateservice.NewUpdateService(MainCTX, clientChan, opcuaChan, updateToClientChan, updateToOpcUaChan, ClientList)
	//go updateService.Update()
	opcuaservice := opcuaservice.NewOpcUaService(MainCTX, opcuaChan, updateToOpcUaChan)
	opcuaservice.StartOPCUA(os.Getenv("OPCUA_Server"))
	go webService.Run(opcuaservice)

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
