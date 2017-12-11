"use strict";

    (function() {
        updateAllLights();
    })();

    /****************************** Aliases ******************************/
    function select(element) {
        return document.querySelector(element);
    }

    function selectAll(element) {
        return document.querySelectorAll(element);
    }

    /****************************** Events ******************************/

    // adjust grayscale of the light image
    selectAll('.lightSlider').forEach(element => {
        element.onchange = function() {
            let light = this.closest('.device').querySelector('.light');
            // light.style.opacity = (this.value / 100).toString();
            light.style.filter = "grayscale(" + (1 - this.value / 100) + ")";

            let room = this.closest('.roomContainer').querySelector('input[name="room"]').value;
            let device = this.closest('.device').querySelector('data').value;
            updateState(room, device, this.value);
        }
    });

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
                // let str = String.fromCharCode.apply(String, this.responseText);
                console.log(typeof(this.responseText));
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

    function updateState(room, device, state) {
        let data = {
            room:   room,
            device: device,
            state:  state,
            action: "updateState"
        };

        ajax({
            method:     'POST',
            data:       data,
            success:    function (resp) {

            }
        })
    }

    function updateAllLights() {
        selectAll('.light').forEach(element => {
            let light = element.closest('.device').querySelector('.light');
            let value = element.closest('.device').querySelector('input').value;
            light.style.filter = "grayscale(" + (1 - value / 100) + ")";
        });
    }

    // Adds text to an element and removes it after a short time
    function showInfo(element, text) {
        element.innerText = text;
        setTimeout( _ => {
            element.innerText = "";
        }, 3000);
    }

    function updateRooms() {
        console.log("Updating rooms");
        ajax({
            method: 'POST',
            uri:    '',
            data:   {action: "getRooms"},
            success: function(rooms) {
                let list = select('#roomList');
                removeChildren(list);
                rooms.forEach(element => {
                    let tag = '<li><input data-value="' + element.id + '" value="' + element.name + '" readonly="" type="text"></li>';
                    let test = getNewElement(tag);
                    console.log(test);
                    list.appendChild(test);
                });
            }
        });
    }

    // removes all children from an element
    function removeChildren(parent) {
        while (parent.firstChild) {
            parent.removeChild(parent.firstChild);
        }
    }

    // Returns a HTML element created from a string
    // https://stackoverflow.com/questions/494143/creating-a-new-dom-element-from-an-html-string-using-built-in-dom-methods-or-pro
    function getNewElement(string) {
        let div = document.createElement('div');
        div.innerHTML = string;
        return div.firstChild;
    }