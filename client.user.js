// ==UserScript==
// @name         Browsing Reminder Client
// @namespace    JustSong
// @version      0.1
// @description  Remind you when you waste too much time on specified websites.
// @author       JustSong
// @match        https://www.zhihu.com/*
// @match        https://github.com/*
// @match        https://*.youtube.com/*
// @match        https://*.bilibili.com/*
// ==/UserScript==

(function () {
    'use strict';
    const serverUrl = "https://ping.iamazing.cn";
    const pingInterval = 60; // Unit is second.
    const maxMinutes = 10;
    const needConfirm = false;
    window.pingServer = function () {
        console.log("ping~");
        fetch(serverUrl, {
            method: 'POST',
            mode: 'cors',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                Host: window.location.hostname
            })
        }).then(function (response) {
            return response.text().then(function (minutes) {
                let counter = parseInt(minutes);
                window.processRequest(counter);
            })
        }).catch(function (reason) {
            console.log(reason)
        })
    };

    window.processRequest = function (minutes) {
        console.log(`You have wasted ${minutes} minutes.`);
        if (minutes >= maxMinutes) {
            let choose = true;
            if (needConfirm) {
                choose = confirm(`You have wasted ${minutes} minutes in this site, would you like to close it?`);
            }
            if (choose) {
                // window.close()
                window.location.href = "https://google.com"
            } else {
                fetch(`${serverUrl}/clear`, {
                    method: 'POST',
                    mode: 'cors',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        Host: window.location.hostname
                    })
                }).finally()
            }
        }
    };
    window.pingServer();
    let timer = setInterval(window.pingServer, pingInterval * 1000);
    document.addEventListener("visibilitychange", function() {
        if (document.visibilityState === 'visible') {
            window.pingServer();
            timer = setInterval(window.pingServer, pingInterval * 1000);
        } else {
            clearInterval(timer)
        }
    });
})();