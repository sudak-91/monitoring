package main

import (
	"context"
	"log"

	"github.com/sudak-91/monitoring/internal/pkg/client"
	"github.com/sudak-91/monitoring/internal/pkg/client/messageservice"
	"github.com/sudak-91/monitoring/internal/pkg/client/render"
	"github.com/sudak-91/monitoring/internal/pkg/client/screens"
	"github.com/sudak-91/monitoring/pkg/message"
	command "github.com/sudak-91/monitoring/pkg/message/command"

	"github.com/sudak-91/wasmhtml/cookie"
	"github.com/sudak-91/wasmhtml/element"
)

func main() {
	ctx := context.Background()
	cookie := cookie.NewCookie()
	messageServceChan := make(chan interface{})
	MessageService := messageservice.NewMessageService(ctx, messageServceChan)
	client := client.NewClient(ctx, cookie, MessageService)
	done := make(chan bool)
	go client.Run(done)
	<-done
	uuid, err := cookie.GetValue("UUID")
	if err != nil {
		command := command.GetUUIDCommand()
		data, err := message.EncodeData(command)
		if err != nil {
			log.Println(err.Error())
		}
		err = client.Request(data)
		if err != nil {
			panic(err.Error())
		}
	} else {
		command := command.SetUUIDCommand(uuid)
		data, err := message.EncodeData(command)
		if err != nil {
			log.Println(err.Error())
		}
		err = client.Request(data)
		if err != nil {
			panic(err.Error())
		}
	}
	screenChan := make(chan interface{})
	renderChan := make(chan interface{})
	render := render.NewRender(ctx, renderChan, screenChan, messageServceChan)
	MainPage := screens.NewMainScreen(renderChan, screenChan, element.GetBody(), client)
	render.AddScreen("main", MainPage)
	go render.Run()
	renderChan <- "main"
	//body := element.GetBody() //Получение <BODY>
	//body.Generate()           //Генерация старницы
	<-ctx.Done()

}
