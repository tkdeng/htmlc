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
make install &&\
cd ../ && rm -r htmlc

# install into /usr/bin
make install

# install locally (with dependencies)
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

```shell
htmlc ./src --out="./output.exs"
```

You can also specify a port number, to automatically start a static-like http server.

```shell
htmlc ./src --port="3000"
```

Note: by default, "--out" is set to the same directory, with the file name set to the base folder name.

## Elixir Template Usage

```shell
# compile
htmlc ./src

# start template engine
elixir ./html.exs

# render page
> render, mypage/home, mylayout/layout, eyJqc29uIjogImFyZ3MifQ==
#                                       base64({"json": "args"})

# render widget (optional)
> widget, mywidget/app, eyJqc29uIjogImFyZ3MifQ==
#                       base64({"json": "args"})

# render layout (optional)
> layout, mylayout/layout, eyJqc29uIjogImFyZ3MifQ==, eyJqc29uIjogImh0bWwgY29udGVudCJ9
#                          base64({"json": "args"}), base64({"json": "html content"})

# stop template engine
stop
```

### IEX Template Usage

```elixir
# compile
htmlc --iex ./src

# start template engine
iex ./html.exs

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
