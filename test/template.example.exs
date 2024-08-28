defmodule App do
  @map_layout %{
    "layouts/layout" => :layout,
  } #_map_layout

  @map_widget %{
    "widgets/app" => :app,
  } #_map_widget

  @map_page %{
    "pages/index" => :index,
  } #_map_page

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
      apply(WIDGET, widget, [args])
    else
      '{Error 500: Widget Not Found!}'
    end
  end

  def layout(layout, args, cont) do
    if @map_layout[layout] do
      apply(LAYOUT, layout, [args, cont])
    else
      Enum.reduce(cont, '', fn {_, val}, a ->
        '#{a} #{val}'
      end)
    end
  end
end

defmodule LAYOUT do
  def layout(args, cont) do
'<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, height=device-height, initial-scale=1.0, minimum-scale=1.0"/>
  <meta name="description" content="#{args[:desc]}"/>
  <title>#{args[:title]}</title>
  #{cont[:head]}
</head>
<body>
  #{cont[:body]}
</body>
</html>'
  end
end #_LAYOUT

defmodule WIDGET do
  def app(args) do
'#{args[:n]} * 2 = #{
  args[:n] * 2
}'
  end
end #_WIDGET

defmodule PAGE do
  def index(layout, args) do
    App.layout layout, args, %{
      head:
'<link rel="stylesheet" href="/style.css">',
      body:
'<h1>Hello, World</h1>
#{App.widget :app, Map.merge(args, %{
  n: 2,
})}',
    }
  end
end #_PAGE
