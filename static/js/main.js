"use strict";

(function() {

        if ( select('body').id === "controlBody" ) {
            updateAllDevices();
            updateSceneDates();

            // Update control UI every second
            setInterval( _ => {
                ajax({
                    method: "POST",
                    url: "",
                    data: {action: "getDevices"},
                    success: function(resp) {
                        if ( resp.Msg ) {
                            updateControlUI(resp);
                        }
                    }
                })
            }, 1000)
        }

        if ( select('body').id === "simBody" ) {
            // Update Sim UI every second if running
            setInterval( _ => {
                if ( select('#simStateData').value === "running" ) {
                    ajax({
                        method: "POST",
                        url: "",
                        data: {action: "getSim"},
                        success: function(resp) {
                            updateSimUI(resp);
                        }
                    })
                    hide(select('#xmlContainer'));
                } else {
                    show(select('#xmlContainer'));
                }
            }, 1000)
        }
    })();

    /****************************** Aliases ******************************/

    function select(element) {
        return document.querySelector(element);
    }

    function selectAll(element) {
        return document.querySelectorAll(element);
    }

    /****************************** Events ******************************/

    // Events for the control UI
    if ( select('body').id === "controlBody" ) {
        // adjust grayscale of the light image
        selectAll('input[name^="lightState"]').forEach(element => {
            element.onchange = function() {
                let light = this.closest('.device').querySelector('.light');
                light.style.filter = this.value === "on" ? "grayscale(0%)" : "grayscale(100%)";
                let room = this.closest('.roomContainer').querySelector('input[name="room"]').value;
                let device = this.closest('.device').querySelector('data').value;
                updateState(room, device, this.value);
            }
        });

        // update graphic onchange
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

        // update graphic onchange
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
    }

    // Events for the sim UI
    if ( select('body').id === "simBody" ) {

        // Change the simulation state
        selectAll('.simButton').forEach(element => {
            element.onclick = function() {
                changeSimState(this.value);
            }
        });
    }


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
            let value = element.closest('.device').querySelector('input[name^="lightState"]:checked').value;
            element.style.filter = value === "on" ? "grayscale(0%)" : "grayscale(100%)";
        });

        selectAll('.shutter').forEach(element => {
            let value = element.closest('.device').querySelector('input').value;
            value = -48 * (value/100);
            element.style.backgroundPositionY =  "top " + value + "px, bottom";
        });

        selectAll('.coffee').forEach(element => {
            let value = element.closest('.device').querySelector('select').value;
            value === "default" ? element.style.filter = "grayscale(100%)" : element.style.filter = "grayscale(0%)";
        });
    }

    // Format and show date
    function updateSceneDates() {
        selectAll('.sceneContainer').forEach(element => {
            let date = element.querySelector('data').value.split(" ");
            let dateString = date[0] + "T" + date[1];
            element.querySelector('.dateContainer').innerHTML = formatDate(dateString);
        });
    }

    function show(element) {
        element.style.display = "block";
    }

    function hide(element) {
        element.style.display = "none";
    }

    function changeSimState(state) {
        let data = {};

        if ( state === "Start" ) {
            data = {
                fff: select('#fastForward').value,
                date: select('#simDate').value,
                time: select('#simTime').value,
                action: "start"
            }
        } else if ( state === "Toggle" ) {
            data = {
                action: "toggle"
            }
        }

        ajax({
            method: 'POST',
            uri:    '',
            data:  data,
            success: function(resp) {
                select('#simStateData').value = resp.msg;
                select('#simStateContainer').innerText = "Simulator is " + resp.msg;
            }
        })
    }

    function updateSimUI(data) {
        let date = formatDate(data.simTime);
        select('#simTimeContainer').innerText = date;
        select('#simStateContainer').innerText = "Simulator is " + data.State;
        select('#simStateData').value = data.State;
    }

    function updateControlUI(data) {
        let devices = data.Devices;
        devices.forEach(device => {
            let id = device.Id;
            let state = device.State;
            let type = device.Type;

            let deviceContainer = select('data[value="' + id + '"]').closest('div.device');

            if ( type === "Light" ) {
                deviceContainer.querySelector('input[value="' + state + '"]').checked = true;
            } else if ( type === "Shutter" ) {
                deviceContainer.querySelector('input').value = state;
            } else if ( type === "Coffee machine" ) {
                deviceContainer.querySelector('select').value = state;
            }
        });
        updateAllDevices();
    }

    function formatDate(dateString) {
        let date = new Date(dateString);
        let day = date.getDate();
        let month = date.getMonth() + 1;
        let year = date.getFullYear();
        let hours = date.getHours().toString();
        let minutes = date.getMinutes().toString();
        let secs = date.getSeconds().toString();

        hours = hours.length === 1 ? "0" + hours : hours;
        minutes = minutes.length === 1 ? "0" + minutes : minutes;
        secs = secs.length === 1 ? "0" + secs : secs;

        return day + "/" + month + "/" + year + " " + hours + ":" + minutes + ":" + secs;
    }
