"use strict";

    (function() {
        updateAllDevices();
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
            light.style.filter = "grayscale(" + (1 - this.value / 100) + ")";
            let room = this.closest('.roomContainer').querySelector('input[name="room"]').value;
            let device = this.closest('.device').querySelector('data').value;
            updateState(room, device, this.value);
        }
    });

    selectAll('.shutterSlider').forEach(element => {
        element.onchange = function() {
            let value = -48 * (this.value/100);
            let shutter = this.closest('.device').querySelector('.shutter');
            let room = this.closest('.roomContainer').querySelector('input[name="room"]').value;
            let device = this.closest('.device').querySelector('data').value;
            shutter.style.backgroundPositionY =  "top " + value + "px, bottom";
            updateState(room, device, this.value);
        }
    });

    // KAFFEEEEEEEEEEEEEEEE
    selectAll('.coffeeSelect').forEach(element => {
        element.onchange = function() {
            let coffee = this.closest('.device').querySelector('.coffee');

            if ( this.value !== 'default' ) {
                coffee.style.filter = "grayscale(0%)";
            } else {
                coffee.style.filter = "grayscale(100%)";
            }

            let room = this.closest('.roomContainer').querySelector('input[name="room"]').value;
            let device = this.closest('.device').querySelector('data').value;
            updateState(room, device, this.value);
        }
    });

    // show/hide corresponding input field
    selectAll('.timeButton').forEach(element => {
        element.onclick = function () {
            let sunTime = select('input[name="sunTime"]');
            let clockTime = select('input[name="clockTime"]');

            hide(sunTime);
            hide(clockTime);
            sunTime.disabled = true;
            clockTime.disabled = true;

            if ( this.id === 'selectSunrise' ) {
                show(sunTime);
                sunTime.disabled = false;
                sunTime.value = 'Sunrise';
            } else if ( this.id === 'selectSunset' ) {
                show(sunTime);
                sunTime.disabled = false;
                sunTime.value = 'Sunset';
            } else {
                show(clockTime);
                clockTime.disabled = false;
            }
        }
    });

    // show/hide offset input field
    select('#offsetToggle').onclick = function() {
        let offset = select('input[name="offset"]');
        if ( this.checked === true ) {
            show(offset);
            offset.disabled = false;
        } else {
            hide(offset);
            offset.disabled = true;
        }
    };

    // open overlay
    selectAll('.edit').forEach(element => {
        element.onclick = function() {
            show(this.closest('.roomContainer').querySelector('.overlay'));
        }
    });

    // close overlay if user clicks on it
    selectAll('.overlay').forEach(element => {
        element.onclick = function(e) {
            if ( e.target.className === "overlay" ) {
                hide(this.closest('.roomContainer').querySelector('.overlay'));
            }
        }
    });

    // selectAll('.dateData').forEach(element => {
    //     let date = element.value;
    //
    // });


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

    function updateAllDevices() {
        selectAll('.light').forEach(element => {
            let value = element.closest('.device').querySelector('input').value;
            element.style.filter = "grayscale(" + (1 - value / 100) + ")";
        });

        selectAll('.coffee').forEach(element => {
            let value = element.closest('.device').querySelector('select').value;
            value === "default" ? element.style.filter = "grayscale(100%)" : element.style.filter = "grayscale(0%)";
        });

        selectAll('.shutter').forEach(element => {
            let value = element.closest('.device').querySelector('input').value;
            value = -48 * (value/100);
            element.style.backgroundPositionY =  "top " + value + "px, bottom";
        });
    }

    function show(element) {
        element.style.display = "block";
    }

    function hide(element) {
        element.style.display = "none";
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