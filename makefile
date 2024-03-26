IMAGE_NAME=w-api-image
CONTAINER_NAME=w-api-container

api:
	@go run main.go

api-hot:
	@$$(go env GOPATH)/bin/air

docker-build:
	@docker build --platform linux/amd64 -t ${IMAGE_NAME} .
	@echo "docker image ${IMAGE_NAME} build successful..."

docker-run: docker-build
	@if docker ps -a --format '{{.Names}}' | grep -q ${CONTAINER_NAME}; then \
        echo "Container exists.."; \
		docker rm -f ${CONTAINER_NAME}; \
    else \
        echo "Container does not exist"; \
    fi
	@docker run -d -p 9000:8000 --name ${CONTAINER_NAME} ${IMAGE_NAME}
	@echo "docker container ${CONTAINER_NAME} running successfully..."

docker-start:
	@docker start ${CONTAINER_NAME}
	@echo "docker container ${CONTAINER_NAME} running..."

docker-stop:
	@docker stop ${CONTAINER_NAME}
	@echo "docker container ${CONTAINER_NAME} stopped..."