# Chambeo-api-core

Este es el repo del API Core de Chambeo.

El mismo corresponde a un proyecto realizado en Go y es un API monolitica.

El stack tecnologico, esta conformado por lo siguiente:

- Go 1.21
- AWS DynamoDB
- Gin Gonic
- AWS SQS/SNS
- Localstack (para emular AWS localmente)
- Docker

En el siguiente diagrama puede observarse, la interaccion de estos servicios.

### Pasos para levantar el proyecto

- Levantar el docker-compose
- Levantar la aplicacion
- Hacer un request http GET a localhost:8080/ping, deberia retornar pong y 200