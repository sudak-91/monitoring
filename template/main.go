package main

import (
	"context"

	"github.com/sudak-91/monitoring/internal/pkg/client"
	"github.com/sudak-91/wasmhtml/element"
)

func main() {
	ctx := context.Background()
	client := client.NewClient(ctx)
	go client.Run()
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
