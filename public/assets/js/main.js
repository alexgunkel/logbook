$(function() {
    var utility = {},
        elementCount = 0,
        elementCounter = 0,
        $body = $('body'),
        port = $body.data('port'),
        endPoint = $body.data('endpoint'),
        $loader = $('.loader'),
        timer = null;

    utility.lastLogger = "";
    utility.lastLogLevel = 0;

    function setCookie(cname, cvalue, exdays) {
        var d = new Date();
        d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
        var expires = "expires="+d.toUTCString();
        document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
    }

    function getCookie(cname) {
        var name = cname + "=";
        var ca = document.cookie.split(';');
        for(var i = 0; i < ca.length; i++) {
            var c = ca[i];
            while (c.charAt(0) == ' ') {
                c = c.substring(1);
            }
            if (c.indexOf(name) == 0) {
                return c.substring(name.length, c.length);
            }
        }
        return "";
    }

    function checkAutoscroll() {
        return getCookie("autoscroll");
    }

    function activateCopyButton(el) {
        $(el).on( "click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            var $el = $(this),
                $panel = $el.closest('.panel-body-inner'),
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
    }

    var handleOutput = function (el) {

        $(el).find('.js-toggle').on( "click", function(e) {
            e.preventDefault();
            e.stopPropagation();
            $(this).closest('.panel').find('.panel-body').slideToggle(180);
        });
        activateCopyButton($(el).find('.btn-copy'));
        return el;

    };

    // update page title
    document.title = "LogBook » " + window.location.protocol + '//' + window.location.hostname;

    function LogEntry(data) {
        this.logger = data.logger
        this.severity = (data.severity) ? data.severity : '7';
        this.message = data.message;
        this.time = data.time;
        this.application = data.application;
        this.request_uri = data.request_uri;
        this.elementCount = ++elementCounter;
        this.getRequestUri = function () {
            if(typeof this.request_uri != 'undefined' && this.request_uri.length) {
                if(this.request_uri.length > 130) {
                    requestLinkText = window.location.hostname + this.request_uri.substring(0,130) + '...';
                } else if(this.request_uri === "/") {
                    requestLinkText = window.location.hostname;
                } else {
                    requestLinkText = window.location.hostname + this.request_uri
                }
                return '<a href="' + this.request_uri + '" title="' + this.request_uri + '">' + requestLinkText + '</a>';
            }
        }
        this.getToggleLink = function () {
            return ' <a class="js-toggle" href="#entry-' + this.elementCount +
                '"><span class="glyphicon glyphicon-zoom-in" title="show more"></span></a>';
        }
        this.getBody = function () {
            return '<div class="panel-body-inner severity-' + this.severity + '" id="entry-\' + this.elementCount + \'">' +
                       '<div class="full-message">' + this.message + '</div>' +
                       '<button class="btn-copy" title="Copy to clipboard">Copy</button>' +
                       '<div class="card-subtitle text-muted">' + this.time + ' - ' + this.application + ' - ' + this.getRequestUri() + '</div>' +
                    '</div>';
        }
        this.getHeader = function () {
            return '<div class="panel-title"><b>' + this.logger + '</b></div>' +
                '<div class="data-message"><span class="loglevel"></span>' + this.message.slice(0,130) + '</div>' +
                '<div>' + this.getToggleLink() + '</div>';
        }
        this.getRowAsHtml = function () {
            return '<div class="panel panel-default" data-loglevel="' + this.severity + '">' +
                       '<div class="panel-heading js-toggle severity-' + this.severity + '">' +
                           this.getHeader() +
                       '</div>' +
                       '<div class="panel-body">' +
                           this.getBody() +
                       '</div>' +
                   '</div>';
        }
    }

    utility.showSlider = function () {
        // timer to set actions after request
        if (timer) {
            clearTimeout(timer); //cancel the previous timer.
            timer = null;
        }
        timer = setTimeout(function() {
            $loader.removeClass('active');
        }, 800);
    }

    utility.handleScrolling = function () {
        if(checkAutoscroll() === 'true') {
            $("html, body").stop().animate({ scrollTop: $(document).height() }, 20);
        }
    }
    
    utility.printEntry = function (logEntry) {
        var d = $('<div class="panel-wrapper"></div>')
                .append(logEntry.getRowAsHtml()),
            lastPanel = $('.panel').last();

        d = handleOutput($(d));

        // if new logger is similar to last one, group within one panel
        if (this.lastLogger.length && logEntry.logger.length && this.lastLogger == logEntry.logger) {
            var newPanel = $(d).find('.panel');

            this.lastLogLevel = lastPanel.data('loglevel');

            lastPanel.find('.panel-body').append('<br>' + logEntry.getBody());
            activateCopyButton(lastPanel.find('.btn-copy').last());

            // Update main Loglevel for panel (set to most serious one)
            if (typeof this.lastLogLevel != 'undefined' && typeof logEntry.severity != 'undefined' && this.lastLogLevel > logEntry.severity) {

                lastPanel.data('loglevel', logEntry.severity);
                lastPanel.find('.panel-heading')
                    .removeClass('severity-' + this.lastLogLevel)
                    .addClass('severity-' + logEntry.severity);

                lastPanel.find('.panel-heading').addClass('my-severity-XXXX');
                lastPanel.data('loglevel', '4');
            }
        } else {
            $(output).append(d);
        }
        this.lastLogger = logEntry.logger;

        this.handleScrolling();
        this.showSlider();
    }

    utility.print = function(message, logger) {
        var d = $('<div class="panel-wrapper"></div>')
                    .append(message);

        d = handleOutput($(d));
        $(output).append(d);
        lastLogger = logger;

        this.handleScrolling();
        this.showSlider()
    };

    utility.printLog = function(data, showContent) {
        var entry = null,
            row;

        // default startmessage (onopen) OR logmessage with panel body
        if(data.severity == '10') {

            function getStartText() {
                var values = [
                        "Ready to log",
                        "Jauchzet frohlogget!",
                        "Let's lock 'n' lol",
                        "Let’s log, dudes!",
                        "Let’s log the house, dudes!",
                        "Oh my log!",
                        "Log Me Amadeus!",
                        "Ready to log"
                ];
                return values[Math.floor(Math.random() * values.length)];
            }

            row = '<div class="panel panel-default">' +
                '<div class="panel-heading  js-toggle severity-10">' +
                '<div class="panel-title"><span class="loglevel"></span><b>' + getStartText() + '</b></div>' +
                '</div>' +
                '</div>';

            utility.print(row, data.logger);
        } else {
            entry = new LogEntry(data);
            utility.print(entry.getRowAsHtml(), data.logger);
        }
    };

    window.addEventListener("load", function(evt) {
        var output = document.getElementById("output"),
            ws,
            host = window.location.hostname;


        ws = new WebSocket('ws://' + host + ':' + port + endPoint);
        ws.onopen = function(evt) {
            var data = {
                "severity" : 10,
                "logger": ""
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
            var entry = new LogEntry(data);
            utility.printEntry(entry, data.logger);
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


    if(checkAutoscroll() === 'true') {
        $('.js-toggle-autoscroll').addClass('active');
    } else {
        $('.js-toggle-autoscroll').removeClass('active');
    }

    $('.js-toggle-autoscroll').on('click', function() {
        if($(this).hasClass('active')) {
            setCookie("autoscroll", false);
            $(this).removeClass('active');
        } else {
            setCookie("autoscroll", true);
            $(this).addClass('active');
        }
        return false;
    });
});
