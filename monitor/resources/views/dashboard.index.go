package views

import "github.com/totoval/framework/view"

func init() {
	view.AddView("totoval_dashboard.index", `
{{define "totoval_dashboard.index"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Totoval Dashboard</title>
</head>
<body>
<div id="main" style="width: 1960px;height:1080px;">
	<div><span id="flows">0</span></div>
</div>


<script>
	var wsUrl = 'ws://127.0.0.1:8080/monitor/dashboard/ws'; //{{ .url }}
    var webSocket = new WebSocket(wsUrl);
	webSocket.onmessage = function (event) {
      data = JSON.parse(event.data) 
      document.getElementById("flows").textContent = data.flows;
	}
</script>


</body>
</html>
{{ end }}
`)
}
