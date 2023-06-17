package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/sudak-91/monitoring/pkg/clientservice"
	opcuaservice "github.com/sudak-91/monitoring/pkg/opcua_service"
	"github.com/sudak-91/monitoring/pkg/webserver"
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
	opcuaChan := make(chan any, 5)
	commandToOpcUaTransferChan := make(chan any, 5)
	opcuaservice := opcuaservice.NewOpcUaService(MainCTX, opcuaChan, commandToOpcUaTransferChan)
	opcuaservice.StartOPCUA(os.Getenv("OPCUA_Server"))
	log.Println("Create Web Service")
	clientService := clientservice.NewClientService(MainCTX, commandToOpcUaTransferChan)
	webService := webserver.NewWebService(MainCTX, clientService)
	log.Println("Create Update Service")
	go webService.Run()
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
