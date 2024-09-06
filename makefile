build all:
	go get -u
	go mod tidy
	go build exec -o htmlc
