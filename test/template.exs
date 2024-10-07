defmodule App do
  @map_layout %{
		"layout" => :_layout_aafwyYT8dMvBxtMd,
  } #_MAP_LAYOUT

  @map_widget %{
		"app" => :_app_aafwy4x92s5eCxl1,
		"md:text" => :md_text_aafwya1IXOyk5DsM,
  } #_MAP_WIDGET

  @map_page %{
		"404" => :_404_aafwyGj8qmIxdzTH,
		"error" => :_error_aafwydKqfLsk5KuA,
		"index" => :_index_aafwylZlsFjTvJDA,
  } #_MAP_PAGE

  def render(name, layout, args) do
    cond do
      @map_page[name] ->
        apply(PAGE, @map_page[name], [layout, args])
      @map_page["#{name}/index"] ->
        apply(PAGE, @map_page["#{name}/index"], [layout, args])
      @map_page["#{name}/404"] ->
        apply(PAGE, @map_page["#{name}/404"], [layout, Map.merge(args, %{
          status: 404,
          error: "Page Not Found!"
        })])
      @map_page["#{name}/error"] ->
        apply(PAGE, @map_page["#{name}/error"], [layout, Map.merge(args, %{
          status: 404,
          error: "Page Not Found!"
        })])
      @map_page[String.replace(name, ~r/\/[^\/]+$/, "/404")] ->
        apply(PAGE, @map_page[String.replace(name, ~r/\/[^\/]+$/, "/404")], [layout, Map.merge(args, %{
          status: 404,
          error: "Page Not Found!"
        })])
      @map_page[String.replace(name, ~r/\/[^\/]+$/, "/error")] ->
        apply(PAGE, @map_page[String.replace(name, ~r/\/[^\/]+$/, "/error")], [layout, Map.merge(args, %{
          status: 404,
          error: "Page Not Found!"
        })])
      @map_page["404"] ->
        apply(PAGE, @map_page["404"], [layout, Map.merge(args, %{
          status: 404,
          error: "Page Not Found!"
        })])
      @map_page["error"] ->
        apply(PAGE, @map_page["error"], [layout, Map.merge(args, %{
          status: 404,
          error: "Page Not Found!"
        })])
      true ->
        "<h1>Error 404</h1>\n<h2>Page Not Found!</h2>"
    end
  end

  def widget(widget, args) do
    if @map_widget[widget] do
      apply(WIDGET, @map_widget[widget], [args])
    else
      "{Error 500: Widget Not Found!}"
    end
  end

  def layout(layout, args, cont) do
    if @map_layout[layout] do
      apply(LAYOUT, @map_layout[layout], [args, cont])
    else
      Enum.reduce(cont, "", fn {_, val}, a ->
        "#{a} #{val}"
      end)
    end
  end

  def escapeHTML(arg) do
    if is_bitstring(arg) do
      String.replace(arg, ~r/[<>&]/, fn (c) ->
        cond do
          c == "<" ->
            "&lt;"
          c == ">" ->
            "&gt;"
          c == "&" ->
            "&amp;"
          true ->
            ""
        end
      end) |> String.replace(~r/&amp;(amp;)*/, "&amp;")
    else
      arg
    end
  end

  def escapeArg(arg) do
    if is_bitstring(arg) do
      String.replace(arg, ~r/([\\]*)([\\"'])/, fn (c) ->
        if rem(String.length(c), 2) == 0 do
          "#{c}"
        else
          "\\#{c}"
        end
      end)
    else
      arg
    end
  end
end

defmodule LAYOUT do
	def _layout_aafwyYT8dMvBxtMd(args, cont) do
'<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, height=device-height, initial-scale=1.0, minimum-scale=1.0"/>
  <meta name="description" content="#{App.escapeArg args[:desc]}"/>
  <title>#{App.escapeHTML args[:title]}</title>
  #{cont[:head]}
</head>
<body>
  #{cont[:body]}
</body>
</html>'
	end
end #_LAYOUT

defmodule WIDGET do
	def _app_aafwy4x92s5eCxl1(args) do
'<div class="widget">
  #{App.escapeHTML args[:n]} * 2 = <%
    args.n * 2
  %>
</div>'
	end
	def md_text_aafwya1IXOyk5DsM(args) do
'\# Hello, Markdown'
	end
end #_WIDGET

defmodule PAGE do
	def _404_aafwyGj8qmIxdzTH(layout, args) do
		App.layout layout, args, %{
			body: '
  <h1>Error 404</h1>
  <h2>Page Not Found!</h2>',
		}
	end
	def _error_aafwydKqfLsk5KuA(layout, args) do
		App.layout layout, args, %{
			body: '
  <h1>Error #{App.escapeHTML args[:status]}</h1>
  <h2>#{App.escapeHTML args[:error]}</h2>',
		}
	end
	def _index_aafwylZlsFjTvJDA(layout, args) do
		App.layout layout, args, %{
			head: '
  <link rel="stylesheet" href="/style.css">',
			body: '
  <h1>Hello, World</h1>
  #{App.widget "app", Map.merge(args, %{
	n: 2,
	body: "
    widget body
  ",
})}

  <main if="main">
    #{args[:main]}
  </main>

  #{App.widget "md:text", args}

  <md>
    [md link](/)
  </md>

  <script>
    console.log("args.msg", args.msg)
  </script>

  <style>
    body {
      color: var(--text, black);
    }
  </style>

  <!-- <ul each="menu">
    <li><a href="{\#url}">{\#name}</a></li>
  </ul> -->',
		}
	end
end #_PAGE
