# HTMLC (HTML Compiler)

[<img src="./assets/icon.png" alt="icon" height="100"/>](./assets/icon.png)

Compiles HTML to Elixir.
Who says HTML is not a programming language?

This module is just another templating engine.
With elixir, we can leverage its ability to call individual functions in `iex` mode.
Elixir can quickly build templates on the fly.

> Notice: This Project Is Still In Beta.

## Installation

```shell
# install the go module
go get github.com/tkdeng/htmlc

# or install the binary
git clone https://github.com/tkdeng/htmlc.git &&\
cd htmlc &&\
make &&\
cd ../ && rm -r htmlc

# install into /usr/bin (default)
make install

# install locally
make local

# build without dependency installation
make build

# install dependencies
make deps

# uninstall htmlc
make clean
```

## Golang Usage

```go
import (
  "github.com/tkdeng/htmlc"
)

func main(){
  htmlc.Compile("./src", "./output.exs")
}
```

## Binary Usage

You can opptionally just use the binary instead of importing the module.

```shell
htmlc --src="./src" --out="./output.exs"
```

You can also specify a port number, to automatically start a static-like http server.

```shell
htmlc --port="3000"
```

Note: by default, "--src" is set to the current working directory,
and "--out" is set to the same directory, with the file name set to the base folder name.

You can also call this method without the "--src" or "--port"

```shell
htmlc --src="/var/www/html"
# is equivalent to
htmlc /var/www/html

htmlc --port="3000"
# is equivalent to
htmlc 3000

# so you can use the method like this
htmlc /var/www/html 3000

# and the order doesnt matter, as long as the port number is a valid uint16
htmlc 3000 /var/www/html

# note: the output file must still be specified with "--out"
htmlc --out="html.exs" /var/www/html 3000
```

## Elixir Template Usage

```shell
# compile
./htmlc

# start template engine
elixir "./html.exs"

# render page
> render, mypage/home, mylayout/layout, eyJqc29uIjogImFyZ3MifQ== # base64({"json": "args"})

# render widget (optional)
> widget, mywidget/app, eyJqc29uIjogImFyZ3MifQ== # base64({"json": "args"})

# render layout (optional)
> layout, mylayout/layout, eyJqc29uIjogImFyZ3MifQ==, eyJqc29uIjogImh0bWwgY29udGVudCJ9
#                          base64({"json": "args"}), base64({"json": "html content"})

# stop template engine
stop
```

### IEX Template Usage

```elixir
# compile
./htmlc --iex

# start template engine
iex "./html.exs"

# render page
iex> App.render "mypage/home", "mylayout/layout", %{args: "myarg"}

# render widget (optional)
iex> App.widget "mywidget/app", %{args: "myarg"}

# render layout (optional)
iex> App.layout "mylayout/layout", %{args: "myarg"}, %{body: "page embed"}

# stop template engine
iex> System.halt
```

## HTML Templates

### Layout

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, height=device-height, initial-scale=1.0, minimum-scale=1.0"/>
  <meta name="description" content="{desc}"/>
  <title>{title}<!-- {variable} --></title>
  {@head} <!-- embed page head -->
</head>
<body>
  {@body} <!-- embed page body -->
</body>
</html>
```

### Page

```html
<_@head> <!-- embed into layout {@head} reference -->
  <link rel="stylesheet" href="/style.css">
</_@head>

<_@body> <!-- embed into {@body} -->
  <h1>Hello, World</h1>

  <!-- use `<_#name>` to embed widgets -->
  <_#app n="2">
    widget body
  </_#app>

  <main>
    {&main} <!-- {&variable} use `&` to allow raw HTML -->
  </main>
</_@body>
```

### Widget

```html
<div class="widget">
  <!-- use <% scripts %> to run elixir (feature not yet implemented) -->
  {n} * 2 = <%
    args.n * 2
  %>
</div>

<!-- markdown not yet implemented -->
<md>
  Markdown
</md>

<markdown>
  Markdown
</markdown>
```
