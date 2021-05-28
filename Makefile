build:
	@echo Building all containers
	@docker build --build-arg service=accounting -t accounting-service:latest .
	@docker build --build-arg service=accounting --build-arg cmd=cdc -t accounting-cdc:latest .

	@docker build --build-arg service=consumer -t consumer-service:latest .
	@docker build --build-arg service=consumer --build-arg cmd=cdc -t consumer-cdc:latest .

	@docker build --build-arg service=delivery -t delivery-service:latest .

	@docker build --build-arg service=kitchen -t kitchen-service:latest .
	@docker build --build-arg service=kitchen --build-arg cmd=cdc -t kitchen-cdc:latest .

	@docker build --build-arg service=order -t order-service:latest .
	@docker build --build-arg service=order --build-arg cmd=cdc -t order-cdc:latest .

	@docker build --build-arg service=order-history -t order-history-service:latest .

	@docker build --build-arg service=restaurant -t restaurant-service:latest .
	@docker build --build-arg service=restaurant --build-arg cmd=cdc -t restaurant-cdc:latest .

	@docker build --build-arg service=web-bff -t web-bff-service:latest .

tidy:
	@echo Tidying all mod files
	@cd accounting & go mod tidy & cd ..
	@cd consumer & go mod tidy & cd ..
	@cd delivery & go mod tidy & cd ..
	@cd kitchen & go mod tidy & cd ..
	@cd order & go mod tidy & cd ..
	@cd order-history & go mod tidy & cd ..
	@cd restaurant & go mod tidy & cd ..
	@cd serviceapis & go mod tidy & cd ..
	@cd shared-go & go mod tidy & cd ..
	@cd web-bff & go mod tidy & cd ..

shiny:
	@echo Updating all dependencies
	@cd accounting & go get -u ./... & cd ..
	@cd consumer & go get -u ./... & cd ..
	@cd delivery & go get -u ./... & cd ..
	@cd kitchen & go get -u ./... & cd ..
	@cd order & go get -u ./... & cd ..
	@cd order-history & go get -u ./... & cd ..
	@cd restaurant & go get -u ./... & cd ..
	@cd serviceapis & go get -u ./... & cd ..
	@cd shared-go & go get -u ./... & cd ..
	@cd web-bff & go get -u ./... & cd ..
