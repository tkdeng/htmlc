build all:
	make install
	go mod tidy
	go build -C exec -o ../htmlc

install:
ifeq (,$(wildcard $(/usr/bin/dnf)))
	sudo dnf install pcre-devel
	sudo dnf install elixir erlang
else ifeq (,$(wildcard $(/usr/bin/apt)))
	sudo apt install libpcre3-dev
	sudo add-apt-repository ppa:rabbitmq/rabbitmq-erlang
	sudo apt update
	sudo apt install elixir erlang-dev erlang-xmerl
else ifeq (,$(wildcard $(/usr/bin/apt)))
	sudo yum install pcre-dev
	pacman -S elixir
endif

clean:
	rm ./htmlc
