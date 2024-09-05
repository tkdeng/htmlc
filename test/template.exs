defmodule App do
  @map_layout %{
		"layout" => :_layout_a9uflWsPjFdA3F2p,
  } #_MAP_LAYOUT

  @map_widget %{
		"app" => :_app_a9uflFuS0NiFj3b5,
  } #_MAP_WIDGET

  @map_page %{
		"404" => :_404_a9ufl4Kpx1qdSFxr,
		"error" => :_error_a9ufln4jqU07HHPO,
		"index" => :_index_a9uflfMZOMlnGSSp,
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
          error: 'Page Not Found!'
        })])
      @map_page["#{name}/error"] ->
        apply(PAGE, @map_page["#{name}/error"], [layout, Map.merge(args, %{
          status: 404,
          error: 'Page Not Found!'
        })])
      @map_page[String.replace(name, ~r/\/[^\/]+$/, "/404")] ->
        apply(PAGE, @map_page[String.replace(name, ~r/\/[^\/]+$/, "/404")], [layout, Map.merge(args, %{
          status: 404,
          error: 'Page Not Found!'
        })])
      @map_page[String.replace(name, ~r/\/[^\/]+$/, "/error")] ->
        apply(PAGE, @map_page[String.replace(name, ~r/\/[^\/]+$/, "/error")], [layout, Map.merge(args, %{
          status: 404,
          error: 'Page Not Found!'
        })])
      @map_page["404"] ->
        apply(PAGE, @map_page["404"], [layout, Map.merge(args, %{
          status: 404,
          error: 'Page Not Found!'
        })])
      @map_page["error"] ->
        apply(PAGE, @map_page["error"], [layout, Map.merge(args, %{
          status: 404,
          error: 'Page Not Found!'
        })])
      true ->
        '<h1>Error 404</h1>\n<h2>Page Not Found!</h2>'
    end
  end

  def widget(widget, args) do
    if @map_widget[widget] do
      apply(WIDGET, @map_widget[widget], [args])
    else
      '{Error 500: Widget Not Found!}'
    end
  end

  def layout(layout, args, cont) do
    if @map_layout[layout] do
      apply(LAYOUT, @map_layout[layout], [args, cont])
    else
      Enum.reduce(cont, '', fn {_, val}, a ->
        '#{a} #{val}'
      end)
    end
  end

  def escapeHTML(arg) do
    #todo: add regex to escape html
    '#{arg}'
  end

  def escapeArg(arg) do
    #todo: add regex to escape html arg in string
    '#{arg}'
  end
end

defmodule LAYOUT do
	def _layout_a9uflWsPjFdA3F2p(args, cont) do
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
</html>
'
	end
end #_LAYOUT

defmodule WIDGET do
	def _app_a9uflFuS0NiFj3b5(args) do
'#{App.escapeHTML args[:n]} * 2 = <%
  args.n * 2
%>
'
	end
end #_WIDGET

defmodule PAGE do
	def _404_a9ufl4Kpx1qdSFxr(layout, args) do
		App.layout layout, args, %{
			body:
'
  <h1>Error 404</h1>
  <h2>Page Not Found!</h2>
',
		}
	end
	def _error_a9ufln4jqU07HHPO(layout, args) do
		App.layout layout, args, %{
			body:
'
  <h1>Error #{App.escapeHTML args[:status]}</h1>
  <h2>#{App.escapeHTML args[:error]}</h2>
',
		}
	end
	def _index_a9uflfMZOMlnGSSp(layout, args) do
		App.layout layout, args, %{
			head:
'
  <link rel="stylesheet" href="/style.css">
',
			body:
'
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

  <ul each="menu">
    <li><a href="#{this[:url]}">#{this[:name]}</a></li>
  </ul>
',
		}
	end
end #_PAGE
