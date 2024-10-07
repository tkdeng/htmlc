build all:
	go mod tidy
	go build -C ./exec -o ../htmlc

install:
	make dependencies
	make build
	sudo cp ./htmlc /usr/bin/htmlc

local:
	make dependencies
	make build

dev:
	make build
	sudo cp ./htmlc /usr/bin/htmlc

deps dependencies:
ifeq (,$(wildcard $(/usr/bin/dnf)))
	sudo dnf install pcre-devel
	sudo dnf install elixir erlang
	sudo dnf install go
else ifeq (,$(wildcard $(/usr/bin/apt)))
	sudo apt install libpcre3-dev
	sudo add-apt-repository ppa:rabbitmq/rabbitmq-erlang
	sudo apt update
	sudo apt install elixir erlang-dev erlang-xmerl
else ifeq (,$(wildcard $(/usr/bin/yum)))
	sudo yum install pcre-dev
	pacman -S elixir
endif

clean:
	rm ./htmlc
	sudo rm /usr/bin/htmlc

test-fiver:
	cd ./gofiber && go test
