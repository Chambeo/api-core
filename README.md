# Chambeo-api-core

Este es el repo del API Core de Chambeo.

El mismo corresponde a un proyecto realizado en Go y es un API monolitica.

El stack tecnologico, esta conformado por lo siguiente:

- Go 1.21
- Gin Gonic
- Postgresql 12
- AWS SQS/SNS
- Localstack (AWS localmente)
- Docker

En el siguiente diagrama puede observarse, la interaccion de estos servicios.

### Pasos para levantar el proyecto

- Levantar el docker-compose ejecutando el comando ```docker-compose up```
- Situarse en el directorio ```cmd/api``` y levantar la aplicacion ejecutando el comando ```go run main.go```
- Hacer un request http GET a ```localhost:8080/ping```, deberia retornar 

 ```
{
    "message": "pong"
} 
```
