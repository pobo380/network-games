.PHONY: build clean deploy deploy-function

build: bin/on_connect bin/on_disconnect bin/join_room bin/game_action bin/debug

bin/%: handler/cmd/%/main.go
	env GOOS=linux go build -ldflags="-s -w" -o $@ $^

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

deploy-function: clean build
	sls deploy function -f $(FNAME)
