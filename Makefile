.PHONY: app-start app-stop

app-start:
	docker-compose up --build --remove-orphans --force-recreate
app-stop:
	docker-compose down

lint:
	golangci-lint run ./...