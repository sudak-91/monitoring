package main

import (
	"context"
	"log"

	clientcontroller "github.com/sudak-91/monitoring/internal/pkg/client/controller"
	cliententity "github.com/sudak-91/monitoring/internal/pkg/client/entity"
	clientusecase "github.com/sudak-91/monitoring/internal/pkg/client/use_case"
	msusecase "github.com/sudak-91/monitoring/internal/pkg/messageservice/use_case"
	"github.com/sudak-91/monitoring/internal/pkg/render"
	"github.com/sudak-91/monitoring/internal/pkg/screens"
	wscontroller "github.com/sudak-91/monitoring/internal/pkg/ws/controller"
	wsservice "github.com/sudak-91/monitoring/internal/pkg/ws/service"
	wsusecase "github.com/sudak-91/monitoring/internal/pkg/ws/use_case"
	"github.com/sudak-91/wasmhtml/element"
)

func main() {
	ctx := context.Background()
	messageServceChan := make(chan interface{})
	client := cliententity.NewClient()
	ws, err := wsservice.NewWSService(ctx)
	if err != nil {
		panic(err)
	}
	messageUseCase := msusecase.NewMSUseCase(ctx, client, messageServceChan)
	clientUseCase := clientusecase.NewUserUseCase(client)
	clientController := clientcontroller.NewClientController(clientUseCase, messageUseCase, ws)
	wsuc := wsusecase.NewWSUseCase(ctx, ws)

	wsController := wscontroller.NewWSController(wsuc, messageUseCase)
	go wsController.Run()

	err = clientController.InitUUID()
	if err != nil {
		err := clientController.GetUUIDFromServer(ctx)
		if err != nil {
			panic(err)
		}
	} else {
		err := clientController.SetUUID(ctx)
		if err != nil {
			log.Panicf("[Main]Has error: %s", err.Error())
			panic(err)
		}
	}

	screenChan := make(chan interface{})
	renderChan := make(chan interface{})
	render := render.NewRender(ctx, renderChan, screenChan, messageServceChan)
	MainPage := screens.NewMainScreen(renderChan, screenChan, element.GetBody(), ws)
	render.AddScreen("main", MainPage)
	go render.Run()
	renderChan <- "main"
	//body := element.GetBody() //Получение <BODY>
	//body.Generate()           //Генерация старницы
	<-ctx.Done()

}
