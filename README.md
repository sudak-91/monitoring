# monitoring


# Обмен на стороне сервера
При запуске сервера стартуют следующие сервисы:
* WebServer
* OpcUaServer
* UpdateServer


При подключении web server записывает данные о клиенте, создает объект клиента и запускает client service который обрабатывает все данные на прием передачу.\
WebServer -> Client ->ClientService

Взаимодействие между сервисами происходит по каналам.\
ClientService -- *Clilent Chan*-->          UpdateService ---- *Update To OpcUA Chan*-----> OpcUa Service\
ClientService <- *Update To Client Chan*--  UpdateService       <--- *Opc Ua Chan*-----     OpcUa Service \

## UpdateService
Сервис отвечает за обработку всех событий и комманд, а также изменяет серверную информацию. Является связующим звеном.\