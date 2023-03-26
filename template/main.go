package main

import (
	"context"
	"log"

	"github.com/sudak-91/monitoring/internal/pkg/client"
	"github.com/sudak-91/monitoring/internal/pkg/client/command"
	"github.com/sudak-91/monitoring/internal/pkg/client/render"
	screen "github.com/sudak-91/monitoring/internal/pkg/client/render/template"
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
	screenChan := make(chan string)
	render := render.NewRender(ctx, screenChan)
	MainPage := screen.NewMianPage()
	render.AddScreen("main", MainPage)
	go render.Run()
	screenChan <- "main"
	body := element.GetBody() //Получение <BODY>
	body.Generate()           //Генерация старницы
	<-ctx.Done()

}
