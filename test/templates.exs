Mix.install([:jason])

defmodule App do
  @map_layout %{
		"layout" => :_layout_aafya8VUtWAGvCSI,
  } #_MAP_LAYOUT

  @map_widget %{
		"app" => :_app_aafyabTNCOuCvQC4,
		"md:text" => :md_text_aafyaMOLDfvmAEUH,
  } #_MAP_WIDGET

  @map_page %{
		"404" => :_404_aafyaKe4kY0cAMyl,
		"error" => :_error_aafya6XJpMiubaNU,
		"index" => :_index_aafya588VnM2N94u,
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
  end #_LISTEN

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
	def _layout_aafya8VUtWAGvCSI(args, cont) do
		"<<PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9ImVuIj4KPGhlYWQ+CiAgPG1ldGEgY2hhcnNldD0iVVRGLTgiLz4KICA8bWV0YSBuYW1lPSJ2aWV3cG9ydCIgY29udGVudD0id2lkdGg9ZGV2aWNlLXdpZHRoLCBoZWlnaHQ9ZGV2aWNlLWhlaWdodCwgaW5pdGlhbC1zY2FsZT0xLjAsIG1pbmltdW0tc2NhbGU9MS4wIi8+CiAgPG1ldGEgbmFtZT0iZGVzY3JpcHRpb24iIGNvbnRlbnQ9Ig==>>#{App.escapeArg args[:desc]}<<Ii8+CiAgPHRpdGxlPg==>>#{args[:title]}<<PC90aXRsZT4KICA=>>#{cont[:head]}<<CjwvaGVhZD4KPGJvZHk+CiAg>>#{cont[:body]}<<CjwvYm9keT4KPC9odG1sPg==>>"
	end
end #_LAYOUT

defmodule WIDGET do
	def _app_aafyabTNCOuCvQC4(args) do
		"<<PGRpdiBjbGFzcz0id2lkZ2V0Ij4KICA=>>#{App.escapeHTML args[:n]}<<ICogMiA9IDwlCiAgICBhcmdzLm4gKiAyCiAgJT4KPC9kaXY+>>"
	end
	def md_text_aafyaMOLDfvmAEUH(args) do
		"<<IyBIZWxsbywgTWFya2Rvd24=>>"
	end
end #_WIDGET

defmodule PAGE do
	def _404_aafyaKe4kY0cAMyl(layout, args) do
		App.layout layout, args, %{
			body: "<<CiAgPGgxPkVycm9yIDQwNDwvaDE+CiAgPGgyPlBhZ2UgTm90IEZvdW5kITwvaDI+>>",
		}
	end
	def _error_aafya6XJpMiubaNU(layout, args) do
		App.layout layout, args, %{
			body: "<<CiAgPGgxPkVycm9yIA==>>#{App.escapeHTML args[:status]}<<PC9oMT4KICA8aDI+>>#{App.escapeHTML args[:error]}<<PC9oMj4=>>",
		}
	end
	def _index_aafya588VnM2N94u(layout, args) do
		App.layout layout, args, %{
			head: "<<CiAgPGxpbmsgcmVsPSJzdHlsZXNoZWV0IiBocmVmPSIvc3R5bGUuY3NzIj4=>>",
			body: "<<CiAgPGgxPkhlbGxvLCBXb3JsZCE8L2gxPgogIA==>>#{App.widget "app", Map.merge(args, %{
	body: "<<CiAgICB3aWRnZXQgYm9keQog>> ",
	n: 2,
})}<<CgogIDxtYWluIGlmPSJtYWluIj4KICAgIA==>>#{args[:main]}<<CiAgPC9tYWluPgoKICA=>>#{App.widget "md:text", args}<<CgogIDxtZD4KICAgIFttZCBsaW5rXSgvKQogIDwvbWQ+CgogIDxzY3JpcHQ+CiAgICBjb25zb2xlLmxvZygiYXJncy5tc2ciLCBhcmdzLm1zZykKICA8L3NjcmlwdD4KCiAgPHN0eWxlPgogICAgYm9keSB7CiAgICAgIGNvbG9yOiB2YXIoLS10ZXh0LCBibGFjayk7CiAgICB9CiAgPC9zdHlsZT4KCiAgPCEtLSA8dWwgZWFjaD0ibWVudSI+CiAgICA8bGk+PGEgaHJlZj0ieyN1cmx9Ij57I25hbWV9PC9hPjwvbGk+CiAgPC91bD4gLS0+>>",
		}
	end
end #_PAGE

App.listen()
