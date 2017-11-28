"use strict";

    /****************************** Aliases ******************************/
    function select(element) {
        return document.querySelector(element);
    }

    function selectAll(element) {
        return document.querySelectorAll(element);
    }

    /****************************** Events ******************************/

    select('#createRoom').onclick = function() {
        let room = select('#roomName').value;
        let infoSpan = select('#roomInfo');

        if ( room !== '') {
            let data = {
                room: room,
                action: 'insertRoom'
            };

            ajax({
                method:     'POST',
                data:       data,
                success:    function (model) {
                    showInfo(infoSpan, 'Room inserted');
                }
            })
        } else {
            showInfo(infoSpan, 'Please enter a room');
        }
    };

    select('#createDevice').onclick = function() {
        let type = select('#deviceType').value;
        let name = select('#deviceName').value;
        let infoSpan = select('#deviceInfo');

        if ( type !== 'default' && name !== '' ) {
            let data = {
                type: type,
                name: name,
                action: "insertDevice"
            };

            ajax({
                method: 'POST',
                uri:    '',
                data:   data,
                success: function(model) {
                    let text = "Insertion was successfull";
                    showInfo(infoSpan, text);
                }
            });
        } else if ( type === 'default' ) {
            showInfo(infoSpan, 'Please select a device');
        } else if ( name === '' ) {
            showInfo(infoSpan, 'Please input a name');
        }
    };

    /****************************** Functions ******************************/

    function ajax(config) {
        // http://stackoverflow.com/a/15096979/570336
        // Input: { a = "foo", b = 123 }
        // Output: a=foo&b=123
        let serialize = function (obj) {
            let str = [];
            for (let p in obj) {
                if (obj.hasOwnProperty(p)) {
                    str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
                }
            }
            return str.join("&");
        };

        let xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState === 4) {
                let model = JSON.parse(this.responseText);
                if (this.status === 200) {
                    config.success(model);
                } else {
                    config.failure(model);
                }
            }
        };

        if ( config.uri === undefined) {
            config.uri = '';
        }

        xhttp.open(config.method || "GET", config.uri + "?" + serialize(config.params), true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send(serialize(config.data));
    }

    // Adds text to an element and removes it after a short time
    function showInfo(element, text) {
        element.innerText = text;
        setTimeout( _ => {
            element.innerText = "";
        }, 3000);
    }