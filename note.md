`go build -ldflags "-w -s" -o transws`

- @ques 怎么用 curl 请求一个 ws 信息
- @ques 如果已经启动程序直接关闭进程

```
curl -X POST -d "msg=value1" "127.0.0.1:60829/send"

curl -X POST -H "Content-Type: application/json" -d '{"time_list": [[416.6, 419.3], [446.6, 450.2], [474.2, 476], [557
, 560.3], [611.2, 615.6], [625.7, 629.3], [672, 674.6], [754.1, 757.3], [790.8, 793.7], [819.5, 821.6]], "action": "li
st_loop", "type": "youtube", "count": 3, "link": "https://www.youtube.com/watch?v=rn9dkV4sVYQ", "cur_index": 0}' "127.
0.0.1:60829/send"
```

```ts
{
  type: "youtube" | "xxx";
  link: string;
  action: string
  count: number,
  time: [number,number]
  timeList: [number,number][],
  cur_index: number,
  count: number,
}
```

```
var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server,
"Send" to send a message to the server and "Close" to close the connection.
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
```
