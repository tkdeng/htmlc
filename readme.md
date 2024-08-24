# HTMLC (HTML Conpiler)

Compiles HTML to Elixir.
Who said HTML is not a programming language?

This module is actually just another templating engine (and a really fast one).

## Installation

```shell
go get github.com/tkdeng/htmlc
```

## Usage

```go
import (
  "github.com/tkdeng/htmlc"
)

func main(){
  htmlc.Compile("./src", "./html.exs")
}
```

## Using The Binary

You can opptionally just use the binary instead of importing the module.

```shell
./htmlc --src="./src" --dist="./output.exs"
```

You can also specify a root path, and use the defaults.

```shell
./htmlc --root="/var/www/html"

# or
./htmlc /var/www/html

# note: this will assume --src="/var/www/html/src" --dist="/var/www/html/html.exs"
```

You can also specify a port number, to automatically start a static-like http server.

```shell
./htmlc 3000

# or with a root
./htmlc /var/www/html 3000

# and the order doesnt matter (as long as the port number is a valid uint16)
./htmlc 3000 /var/www/html
```
