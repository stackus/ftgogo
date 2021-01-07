# ftgogo - event-driven architecture demonstration application 

## Prerequisites

[Docker](https://www.docker.com/)

Optionally [Postman](https://www.postman.com/) can be used to load the api collection for an easy demo.

## Execution

Open a command prompt and then execute the following docker command

### Mac and Linux Users

    COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose up

### Windows Users

    set "COMPOSE_DOCKER_CLI_BUILD=1" & set "DOCKER_BUILDKIT=1" & docker-compose up

Use `Ctrl-C` to stop all services.

## API

The api for each service has been defined in an `openapi.yaml` file. Look for them under each bounded context folder in `cmd/service`.

> A swagger UI may be added sometime later.

## Quick Demo

With the application running using one of the commands above you can make the following calls to see the processes, events, and entities involved.

- Consumer Service: Register Consumer
- Restaurant Service: Create Restaurant
- Order: Create Order

Loading `FTGOGO.postman_collection.json` into Postman will provide pre-built calls with semi-random data that can be used to test the above.

## Origins

This application is a rewrite in Golang of the [FTGO](https://github.com/microservices-patterns/ftgo-application) application. I tried to duplicate the application as closely as possible but there will be some differences due to misinterpretation, misunderstanding, and/or bugs.