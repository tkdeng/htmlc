Mix.install([:jason])

defmodule App do
  @map_layout %{
  } #_MAP_LAYOUT

  @map_widget %{
  } #_MAP_WIDGET

  @map_page %{
  } #_MAP_PAGE

  def listen() do
    input = IO.read(:line) |> String.trim() |> String.split(",")

    cond do
      length(input) == 1 ->
        [arg1] = input
        if arg1 == "stop" do
          System.halt()
        end
      length(input) == 2 ->
        [arg1, arg2] = input
        if arg1 == ":widget" do
          IO.puts apply(App, :widget, [arg2, %{}])
        else if arg1 == "widget" do
          IO.puts apply(App, :widget, [arg2, %{}]) |> String.replace(~r/<<(.*?)>>/, fn(c) ->
            String.slice(c, 2..-3) |> Base.decode64!
          end)
        else if arg1 == ":layout" do
          IO.puts apply(App, :layout, [arg2, %{}, %{}])
        else if arg1 == "layout" do
          IO.puts apply(App, :layout, [arg2, %{}, %{}]) |> String.replace(~r/<<(.*?)>>/, fn(c) ->
            String.slice(c, 2..-3) |> Base.decode64!
          end)
        else if arg1 == ":render" do
          IO.puts apply(App, :render, [arg2, "layout", %{}])
        else
          IO.puts apply(App, :render, [arg2, "layout", %{}]) |> String.replace(~r/<<(.*?)>>/, fn(c) ->
            String.slice(c, 2..-3) |> Base.decode64!
          end)
        end end end end end
      length(input) == 3 ->
        [arg1, arg2, arg3] = input
        if arg1 == ":widget" do
          with {:ok, json} <- Base.decode64!(arg3) |> Jason.decode() do
            IO.puts apply(App, :widget, [arg2, json])
          end
        else if arg1 == "widget" do
          with {:ok, json} <- Base.decode64!(arg3) |> Jason.decode() do
            IO.puts apply(App, :widget, [arg2, json]) |> String.replace(~r/<<(.*?)>>/, fn(c) ->
              String.slice(c, 2..-3) |> Base.decode64!
            end)
          end
        else if arg1 == ":layout" do
          with {:ok, json} <- Base.decode64!(arg3) |> Jason.decode() do
            IO.puts apply(App, :layout, [arg2, json, %{}])
          end
        else if arg1 == "layout" do
          with {:ok, json} <- Base.decode64!(arg3) |> Jason.decode() do
            IO.puts apply(App, :layout, [arg2, json, %{}]) |> String.replace(~r/<<(.*?)>>/, fn(c) ->
              String.slice(c, 2..-3) |> Base.decode64!
            end)
          end
        else if arg1 == ":render" do
          IO.puts apply(App, :render, [arg2, arg3, %{}])
        else
          IO.puts apply(App, :render, [arg2, arg3, %{}]) |> String.replace(~r/<<(.*?)>>/, fn(c) ->
            String.slice(c, 2..-3) |> Base.decode64!
          end)
        end end end end end
      length(input) == 4 ->
        [arg1, arg2, arg3, arg4] = input
        if arg1 == ":layout" do
          with {:ok, json} <- Base.decode64!(arg3) |> Jason.decode(), {:ok, cont} <- Base.decode64!(arg4) |> Jason.decode() do
            IO.puts apply(App, :layout, [arg2, json, cont])
          end
        else if arg1 == "layout" do
          with {:ok, json} <- Base.decode64!(arg3) |> Jason.decode(), {:ok, cont} <- Base.decode64!(arg4) |> Jason.decode() do
            IO.puts apply(App, :layout, [arg2, json, cont]) |> String.replace(~r/<<(.*?)>>/, fn(c) ->
              String.slice(c, 2..-3) |> Base.decode64!
            end)
          end
        else if arg1 == ":render" do
          with {:ok, json} <- Base.decode64!(arg4) |> Jason.decode() do
            IO.puts apply(App, :render, [arg2, arg3, json])
          end
        else
          with {:ok, json} <- Base.decode64!(arg4) |> Jason.decode() do
            IO.puts apply(App, :render, [arg2, arg3, json]) |> String.replace(~r/<<(.*?)>>/, fn(c) ->
              String.slice(c, 2..-3) |> Base.decode64!
            end)
          end
        end end end
      true ->
        IO.puts "<h1>Error 500</h1><h2>Internal Server Error!</h2>"
    end

    App.listen()
  end

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
end #_LAYOUT

defmodule WIDGET do
end #_WIDGET

defmodule PAGE do
end #_PAGE

App.listen()
