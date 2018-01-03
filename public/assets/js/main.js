$(function() {
    var utility = {},
        elementCount = 0,
        port = $('body').data('port'),
        endPoint = $('body').data('endpoint'),
        $loader = $('.loader');

    var handleOutput = function (el) {

        $(el).find('.js-toggle').on( "click", function(e) {
            e.preventDefault();
            $(this).closest('.panel').find('.panel-body').slideToggle();
        });

        return el;

    };

    utility.print = function(message, severity) {
        var d =  $("<div></div>")
            .addClass((severity) ? 'severity-' + severity : 'severity-7')
            .append(message);
        d = handleOutput($(d));

        $(output).append(d);

        var timer = null;
        if (timer) {
            clearTimeout(timer); //cancel the previous timer.
            timer = null;
        }
        timer = setTimeout(function() {
            $loader.removeClass('active');
        }, 10000);
    };



    utility.printLog = function(data, showContent) {
        showContent = showContent || false;
        var toggleLink = '',
            cardtext = '',
            row;

        elementCount++;

        if(showContent) {
            toggleLink = ' <a class="js-toggle" href="#element-' + elementCount + '"><span class="glyphicon glyphicon-zoom-in" title="show more"></span></a>';
            cardtext = '<div class="panel-body" id="element-' + elementCount + '">' +
                '<div class="full-message">' + data.message + '</div>' +
                '<div class="context">' + JSON.stringify(data.context) + '</div>' +
                '</div>';
        }

        row = '<div class="panel panel-default">' +
            '<div class="panel-heading">' +
            '<div class="panel-title"><h5>' + data.logger + '</h5><br>' + data.message.slice(0,130) + '<br>' + data.severity_text + '</div>' +
            '<h6 class="card-subtitle text-muted">' + data.time + ' - ' + data.application + ' - ' + ' <br> ' + data.request_uri + '</h6>' + toggleLink +
            '</div>' + cardtext + '</div>';
        utility.print(row, data.severity);
    };

    window.addEventListener("load", function(evt) {
        var output = document.getElementById("output"),
            ws,
            host = window.location.hostname;


        ws = new WebSocket('ws://' + host + ':' + port + endPoint);
        ws.onopen = function(evt) {
            var data = {
                "message" : "Ready to read",
                "severity" : 7
            };
            utility.printLog(data);
        };
        ws.onclose = function(evt) {
            var data = {
                "message" : "Connection closed",
                "severity" : 7
            };

            utility.printLog(data);
            ws = null;
        };
        ws.onmessage = function(evt) {
            $loader.addClass('active');
            var data = JSON.parse(evt.data);
            console.log(data);
            utility.printLog(data, true);
        };
        ws.onerror = function(evt) {
            var data = {
                "message" : "ERROR: " + evt.data,
                "severity" : 2
            };

            utility.printLog(data);
        }
    });

    window.addEventListener("beforeunload", function(evt) {
        // delete the cookie
        document.cookie = "logbook=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    });
});


