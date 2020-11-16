// ==UserScript==
// @name         Browsing Reminder Client
// @namespace    JustSong
// @version      0.1
// @description  Remind you when you waste too much time on specified websites.
// @author       JustSong
// @match        http://*/*
// @match        https://*/*
// @grant        none
// @connect      localhost
// @connect      ping.iamazing.cn
// ==/UserScript==

(function () {
    'use strict';
    const serverUrl = "http://home.justsong.cn";
    const pingInterval = 1; // Unit is second.
    const maxMinutes = 1;
    window.pingServer = function () {
        fetch(serverUrl, {
            method: 'POST',
            mode: 'no-cors',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                Host: window.location.hostname
            })
        }).then(function (response) {
            return response.text().then(function (minutes) {
                let counter = parseInt(minutes)
                window.process(counter);
            })
        }).catch(function (reason) {
            console.log(reason)
        })
    };

    window.process = function (minutes) {
        if (minutes >= maxMinutes) {
            let choose = confirm(`You have wasted ${minutes} minutes in this site, would you like to close it?`);
            if (choose) {
                // window.close()
                window.location.href = "https://google.com"
            } else {
                fetch(`${serverUrl}/clear`, {
                    method: 'POST',
                    mode: 'no-cors',
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

    function main() {
        setInterval(window.pingServer, pingInterval * 1000)
    }

    main();
})();