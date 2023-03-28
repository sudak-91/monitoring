package main

import (
	"context"
	"log"

	"github.com/sudak-91/monitoring/internal/pkg/client"
	"github.com/sudak-91/monitoring/internal/pkg/client/command"
	"github.com/sudak-91/monitoring/internal/pkg/client/render"
	"github.com/sudak-91/monitoring/internal/pkg/client/screens"

	"github.com/sudak-91/wasmhtml/cookie"
	"github.com/sudak-91/wasmhtml/element"
)

func main() {
	ctx := context.Background()
	cookie := cookie.NewCookie()
	client := client.NewClient(ctx, cookie)
	done := make(chan bool)
	go client.Run(done)
	<-done
	command := command.NewCommand(client)
	uuid, err := cookie.GetValue("UUID")
	if err != nil {
		err = command.GetUUID()
	} else {
		err = command.SetUUID(uuid)
	}
	if err != nil {
		log.Println(err)
	}
	screenChan := make(chan interface{})
	renderChan := make(chan interface{})
	render := render.NewRender(ctx, renderChan, screenChan)
	MainPage := screens.NewMainScreen(renderChan, screenChan, element.GetBody(), command)
	render.AddScreen("main", MainPage)
	go render.Run()
	renderChan <- "main"
	//body := element.GetBody() //Получение <BODY>
	//body.Generate()           //Генерация старницы
	<-ctx.Done()

}
