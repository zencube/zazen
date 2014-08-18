APP_NAME=zazen
VPATH=.:src:$(GOPATH)/src:

update-deps:
	@echo Updating dependencies
	go get github.com/fzzy/radix/redis

build: deps
	go build -o $(APP_NAME) src/main.go

run: deps
	go run src/main.go

clean:
	rm $(APP_NAME)

deps: github.com/fzzy/radix/redis

github.com/fzzy/radix/redis:
	go get github.com/fzzy/radix/redis
