(function () {
  "use strict";
  const url = "ws://localhost:60829/ws";
  var ws = new WebSocket(url);

  ws.onopen = function (evt) {
    console.log("ws:> onopen", url);
  };
  ws.onclose = function (evt) {
    console.log("ws:> close");
    ws = null;
  };
  ws.onmessage = function (evt) {
    const inputElement = document.getElementById("search_input");
    inputElement.value = evt.data;
    const inputEvent = new Event("input");
    inputElement.dispatchEvent(inputEvent);

    document.querySelector(".translate_btn").click();
  };
  ws.onerror = function (evt) {
    console.log("ws:> error " + evt.data);
  };
})();
