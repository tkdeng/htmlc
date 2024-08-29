defmodule App do
  @map_layout %{
  } #_map_layout

  @map_widget %{
  } #_map_widget

  @map_page %{
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
end

defmodule LAYOUT do
end #_LAYOUT

defmodule WIDGET do
end #_WIDGET

defmodule PAGE do
end #_PAGE
