"use strict";

    /****************************** Aliases ******************************/
    function select(element) {
        return document.querySelector(element);
    }

    function selectAll(element) {
        return document.querySelectorAll(element);
    }

    /****************************** Events ******************************/

    select('#createEntity').onclick = function () {
        var type = select('#entityType').value;
        var room = select('#entityRoom').value;
        var name = select('#entityName').value;

        if ( type != 'deafult' && room != 'default' && name != '' ) {
            var data = {
                type: type,
                room: room,
                name: name,
                action: "insertRoom"
            };

            ajax({
                method: 'POST',
                uri:    '',
                data:   data,
                success: function(model) {
                    console.log(model);
                }
            });
        } else {
            // Fehlermeldung
            return;
        }
    };

    /****************************** Functions ******************************/

    function ajax(config) {
        // http://stackoverflow.com/a/15096979/570336
        // Input: { a = "foo", b = 123 }
        // Output: a=foo&b=123
        var serialize = function (obj) {
            var str = [];
            for (var p in obj) {
                if (obj.hasOwnProperty(p)) {
                    str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
                }
            }
            return str.join("&");
        }

        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            if (this.readyState === 4) {
                var model = JSON.parse(this.responseText);
                if (this.status === 200) {
                    config.success(model);
                } else {
                    config.failure(model);
                }
            }
        };

        xhttp.open(config.method || "GET", config.uri + "?" + serialize(config.params), true);
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send(serialize(config.data));
    }
