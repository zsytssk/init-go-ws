// ==UserScript==
// @name         youtube jump
// @namespace    http://tampermonkey.net/
// @version      2024-09-16
// @description  try to take over the world!
// @author       You
// @match        https://www.youtube.com/watch*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=youtube.com
// @grant        GM_setClipboard
// ==/UserScript==

function triggerKeyEvent(
  keyCode,
  opts = {
    altKey: false,
    ctrlKey: false,
    metaKey: false,
    shiftKey: false,
    repeat: false,
  }
) {
  var event = new KeyboardEvent("keydown", {
    keyCode: keyCode,
    which: keyCode,
    bubbles: true,
    cancelable: true,
    ...opts,
  });
  document.dispatchEvent(event);
}

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
    const [type, action, link, time] = evt.data.split("|");
    console.log(`test:>`, [type, action, link, time]);
    if (type !== "youtube" || !location.href.startsWith(link)) {
      return;
    }

    const video = document.querySelector("video");
    switch (action) {
      case "jump":
        video.currentTime = time;
        video.play();
        break;
      case "toggle_subtitle":
        triggerKeyEvent("C".charCodeAt(0));
        break;
      case "reduce_speed":
        triggerKeyEvent(188, { shiftKey: true });
        break;
      case "plus_speed":
        triggerKeyEvent(190, { shiftKey: true });
        break;
      case "play_back":
        triggerKeyEvent(37);
        break;
      case "play_forward":
        triggerKeyEvent(39);
        break;
      case "copy_time":
        GM_setClipboard(video.currentTime.toFixed(1), "text", () =>
          console.log("Clipboard set!")
        );
        break;
      case "toggle_pause":
        triggerKeyEvent(75);
        break;
    }
  };
  ws.onerror = function (evt) {
    console.log("ws:> error " + evt.data);
  };
})();
