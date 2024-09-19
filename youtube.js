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
    console.log("ws:>  ", evt);
    const [type, action, link, time] = evt.data.split("|");
    if (type !== "youtube" || !location.href.startsWith(link)) {
      return;
    }

    const video = document.querySelector("video");
    const subtitles = document.querySelector(".ytp-subtitles-button");
    switch (action) {
      case "jump":
        video.currentTime = time;
        video.play();
        break;
      case "toggle_subtitle":
        subtitles.click();
        break;
      case "reduce_speed":
        video.playbackRate *= 0.75;
        break;
      case "plus_speed":
        video.playbackRate /= 0.75;
        break;
      case "play_back":
        video.currentTime -= 5;
        break;
      case "play_forward":
        video.currentTime += 5;
        break;
      case "toggle_pause":
        if (video.paused) {
          video.play();
        } else {
          video.pause();
        }
        break;
    }
  };
  ws.onerror = function (evt) {
    console.log("ws:> error " + evt.data);
  };
})();
