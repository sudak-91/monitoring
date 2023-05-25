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
	err := loadConfig()
	if err != nil {
		panic(err)
	}
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
	opcuaservice := opcuaservice.NewOpcUaService(MainCTX, opcuaChan, updateToOpcUaChan)
	opcuaservice.StartOPCUA(os.Getenv("OPCUA_Server"))
	go webService.Run(opcuaservice)
	l := make(chan bool)
	<-l
}

func loadConfig() error {
	viper.SetConfigName("service")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	log.Println("Config read")
	return nil
}
