// ==UserScript==
// @name         youtube jump
// @namespace    http://tampermonkey.net/
// @version      2024-09-16
// @description  try to take over the world!
// @author       You
// @match        https://www.youtube.com/watch*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=youtube.com
// @grant        none
// ==/UserScript==

(function () {
  "use strict";
  const url = "http://localhost:60829/ws";
  var ws = new WebSocket(url);

  ws.onopen = function (evt) {
    console.log("ws:> onopen", url);
  };
  ws.onclose = function (evt) {
    console.log("ws:> close");
    ws = null;
  };
  ws.onmessage = function (evt) {
    const [type, time, link] = evt.data.split("|");
    if (type !== "youtube" || !location.href.startsWith(link)) {
      return;
    }
    const video = document.querySelector("video");
    video.currentTime = time;
    video.play();
  };
  ws.onerror = function (evt) {
    console.log("ws:> error " + evt.data);
  };
})();
