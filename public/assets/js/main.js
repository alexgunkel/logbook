$(function() {
    var utility = {},
        elementCount = 0,
        port = $('body').data('port'),
        endPoint = $('body').data('endpoint'),
        $loader = $('.loader');

    var handleOutput = function (el) {

        $(el).find('.js-toggle').on( "click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            $(this).closest('.panel').find('.panel-body').slideToggle(180);
        });
        $(el).find('.btn-copy').on( "click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            var $el = $(this),
                $panel = $el.closest('.panel')
                copyText = $panel.find('.full-message').text() + '\n\n' + $panel.find('.context').text(),
                $temp = $("<textarea>");
            $("body").append($temp);
            $temp.val(copyText).select();
            document.execCommand("copy");
            $el.addClass('active');
            $temp.remove();

            setTimeout(function(){
                $el.removeClass('active');
            }, 800);

        });
        return el;

    };
    // update page title
    document.title = "LogBook » " + window.location.protocol + '//' + window.location.hostname;

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
            panelBody = '',
            requestUri = window.location.protocol + '//' + window.location.hostname + data.request_uri,
            requestLink = '',
            requestLinkText = '',
            row;

        elementCount++;

        if(showContent) {
            toggleLink = ' <a class="js-toggle" href="#element-' + elementCount + '"><span class="glyphicon glyphicon-zoom-in" title="show more"></span></a>';

            if(typeof data.request_uri != 'undefined' && data.request_uri.length) {
                if(data.request_uri.length > 130) {
                    requestLinkText = window.location.hostname + data.request_uri.substring(0,130) + '...';
                } else if(data.request_uri === "/") {
                    requestLinkText = window.location.hostname;
                } else {
                    requestLinkText = window.location.hostname + data.request_uri
                }
                requestLink = '<div><a href="' + requestUri + '" title="' + requestUri + '">' + requestLinkText + '</a></div>';
            }

            panelBody = '<div class="panel-body" id="element-' + elementCount + '"><button class="btn-copy" title="Copy to clipboard">Copy</button>' +
                            '<div>' + data.severity_text + '</div>' +
                            '<div class="card-subtitle text-muted">' + data.time + ' - ' + data.application + ' </div>' +
                            '<div> ' + requestLink + '</div>' +
                            '<div class="full-message">' + data.message + '</div>' +
                            '<div class="context">' + JSON.stringify(data.context) + '</div>' +
                        '</div>';
        }


        if(typeof data.request_uri != 'undefined' && data.request_uri.length) {
            if(data.request_uri.length > 130) {
                requestLinkText = window.location.hostname + data.request_uri.substring(0,130) + '...';
            } else if(data.request_uri === "/") {
                requestLinkText = window.location.hostname;
            } else {
                requestLinkText = window.location.hostname + data.request_uri
            }
            requestLink = '<div><a href="' + requestUri + '" title="' + requestUri + '">' + requestLinkText + '</a></div>';
        }
        if(data.severity == '10') {
            row = '<div class="panel panel-default">' +
                '<div class="panel-heading  js-toggle">' +
                '<div class="panel-title"><b>Ready to log</b></div>' +
                '</div>' +
                '</div>';
        } else {
            row = '<div class="panel panel-default">' +
                '<div class="panel-heading  js-toggle">' +
                '<div class="panel-title"><b>' + data.logger + '</b></div>' +
                '<div class="data-message">' + data.message.slice(0,130) + '</div>' +
                '<div>' + toggleLink + '</div>' +
                '</div>' +
                panelBody +
                '</div>';
        }

        utility.print(row, data.severity);
    };

    window.addEventListener("load", function(evt) {
        var output = document.getElementById("output"),
            ws,
            host = window.location.hostname;


        ws = new WebSocket('ws://' + host + ':' + port + endPoint);
        ws.onopen = function(evt) {
            var data = {
                "severity" : 10
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
            utility.printLog(data, true);
        };
        ws.onerror = function(evt) {
            var data = {
                "message" : "ERROR: " + evt.data,
                "severity" : 2
            };

            utility.printLog(data);
        };
    });

    window.addEventListener("beforeunload", function(evt) {
        // delete the cookie
        document.cookie = "logbook=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    });

    $('#btn-up').on('click', function(){
        $('html, body').animate({scrollTop:0});
        return false;
    });

    $('#btn-down').on('click', function(){
        $('html, body').animate({scrollTop:$(document).height()});
        return false;
    });
});


