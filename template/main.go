package main

import (
	"context"
	"log"

	"github.com/sudak-91/monitoring/internal/pkg/client"
	"github.com/sudak-91/monitoring/internal/pkg/client/command"
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
	body := element.GetBody()     //Получение <BODY>
	leftdiv := body.AddDiv()      //Создание дочернего <DIV>
	leftdiv.Id = "lefyDiv"        //Добавление ID
	leftdiv.AddClass("container") //Добавление клааса
	rightdiv := body.AddDiv()
	rightdiv.Id = "rightDiv"
	rightdiv.InnerHtml = "Click me"
	rightdiv.OnClick = "alert(\"tada\")"
	rightdiv.AddClass("container")
	body.Generate() //Генерация старницы
	<-ctx.Done()

}
