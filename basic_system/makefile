APP_NAME=rtc

all: clean build docker

build:
	go build -o $(APP_NAME)

run: build
	./$(APP_NAME)

docker:
	docker build -t $(APP_NAME) .

docker-run:
	docker run -d -p 8080:8080 --name $(APP_NAME) $(APP_NAME)

clean:
	docker system prune
ifneq ("$(shell docker images | grep $(APP_NAME) | awk '{ print $$3; }')", "") 
	docker images | grep $(APP_NAME) | awk '{ print $$3; }' | xargs -I {} docker rmi {}
endif
	rm -f $(APP_NAME)
