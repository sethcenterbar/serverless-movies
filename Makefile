.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/get get/main.go
	# env GOOS=linux go build -ldflags="-s -w" -o bin/put put/main.go
	# env GOOS=linux go build -ldflags="-s -w" -o bin/post post/main.go
	# env GOOS=linux go build -ldflags="-s -w" -o bin/delete delete/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
