APP_NAME=object-storage-service
PORT=8080

.PHONY: run docker-run test clean

run:
	@echo "Running server locally on port $(PORT)..."
	PORT=$(PORT) go run main.go

docker-run:
	@echo "Building and running Docker container..."
	docker build -t $(APP_NAME) .
	docker run -p $(PORT):$(PORT) -e PORT=$(PORT) $(APP_NAME)

test:
	@echo "Running tests..."
	go test ./... -v

clean:
	@echo "Cleaning up..."
	-rm -f $(APP_NAME)
	-docker rmi $(APP_NAME)