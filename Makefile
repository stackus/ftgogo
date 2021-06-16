.DEFAULT_GOAL := run
.PHONY: tidy shiny build run build-monolith run-monolith

tidy:
	@echo Tidying all mod files
	@cd accounting & go mod tidy & cd .. & echo accounting is tidy
	@cd consumer & go mod tidy & cd .. & echo consumer is tidy
	@cd delivery & go mod tidy & cd ..  & echo delivery is tidy
	@cd kitchen & go mod tidy & cd .. & echo kitchen is tidy
	@cd order & go mod tidy & cd .. & echo order is tidy
	@cd order-history & go mod tidy & cd .. & echo order-history is tidy
	@cd restaurant & go mod tidy & cd .. & echo restaurant is tidy
	@cd serviceapis & go mod tidy & cd .. & echo serviceapis is tidy
	@cd shared-go & go mod tidy & cd .. & echo shared-go is tidy
	@cd customer-web & go mod tidy & cd .. & echo customer-web is tidy
	@cd monolith & go mod tidy & cd .. & echo monolith is tidy

shiny:
	@echo Updating all dependencies
	@cd accounting & go get -u ./... & cd .. & echo accounting has shiny new packages
	@cd consumer & go get -u ./... & cd .. & echo consumer has shiny new packages
	@cd delivery & go get -u ./... & cd ..  & echo delivery has shiny new packages
	@cd kitchen & go get -u ./... & cd .. & echo kitchen has shiny new packages
	@cd order & go get -u ./... & cd .. & echo order has shiny new packages
	@cd order-history & go get -u ./... & cd .. & echo order-history has shiny new packages
	@cd restaurant & go get -u ./... & cd .. & echo restaurant has shiny new packages
	@cd serviceapis & go get -u ./... & cd .. & echo serviceapis has shiny new packages
	@cd shared-go & go get -u ./... & cd .. & echo shared-go has shiny new packages
	@cd customer-web & go get -u ./... & cd .. & echo customer-web has shiny new packages
	@cd monolith & go get -u ./... & cd .. & echo monolith has shiny new packages

build:
	@docker-compose build

run:
	@docker-compose up

build-monolith:
	@docker-compose -f docker-compose-monolith.yml build

run-monolith:
	@docker-compose -p monolith -f docker-compose-monolith.yml up
