$(function() {
    var utility = {},
        elementCount = 0,
        elementCounter = 0,
        timer = null,
        $body = $('body'),
        $loader = $('.loader'),
        $toggleAutoscroll = $body.find('.js-toggle-autoscroll'),
        port = $body.data('port'),
        endPoint = $body.data('endpoint');
    utility.lastLogger = "";
    utility.lastLogLevel = 0;

    function setCookie(cname, cvalue, exdays) {
        var d = new Date(),
            days = (exdays) ? exdays : 7;
        d.setTime(d.getTime() + (days * 24 * 60 * 60 * 1000));
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

        // $(el).find('.js-toggle').on( "click", function(e) {
        //     e.preventDefault();
        //     e.stopPropagation();
        //     $(this).closest('.panel').find('.panel-body').slideToggle(180);
        // });
        // activateCopyButton($(el).find('.btn-copy'));
        return el;

    };

    // update page title
    document.title = "LogBook » " + window.location.protocol + '//' + window.location.hostname;

//     function LogEntry(data) {
//         this.logger = data.logger;
//         this.severity = (data.severity) ? data.severity : '7';
//         this.message = data.message;
//         this.time = data.time;
//         this.application = data.application;
//         this.request_uri = data.request_uri;
//         this.elementCount = ++elementCounter;
//         this.getRequestLink = function () {
//             if(typeof this.request_uri != 'undefined' && this.request_uri.length) {
//
//                 // cut long request linktexts
//                 if(this.request_uri === "/") {
//                     requestLinkText = window.location.hostname;
//                 } else if(this.request_uri.length > 130) {
//                     requestLinkText = window.location.hostname + this.request_uri.substring(0,130) + '...';
//                 } else {
//                     requestLinkText = window.location.hostname + this.request_uri
//                 }
//                 return '<a href="' + this.request_uri + '" title="' + this.request_uri + '">' + requestLinkText + '</a>';
//             }
//         }
//         this.getToggleLink = function () {
//             return '<a class="js-toggle" href="#entry-' + this.elementCount + '"><span class="glyphicon glyphicon-zoom-in" title="show more"></span></a>';
//         }
//         this.getHeader = function () {
//             return '<div class="panel-title"><b>' + this.logger + '</b></div>' +
//                    '<div class="data-message"><span class="loglevel"></span>' + this.message.slice(0,130) + '</div>' +
//                    '<div class="app-info text-muted">' + this.getRequestLink() + '</div>' +
//                    '<div>' +
//                    '</div>';
//         }
//         this.getBody = function () {
//             return '<div class="panel panel-sub panel-default" data-loglevel="' + this.severity + '">' +
//                        '<div class="panel-heading js-toggle severity-' + this.severity + '">' +
//                 this.getHeader() +
//                 '</div>' +
//                 '<div class="panel-body">' +
//                 '<div class="panel-body-inner severity-' + this.severity + '" id="entry-' + this.elementCount + '">' +
//                 '<span class="loglevel"></span>' +
//                 '<button class="btn-copy" title="Copy to clipboard">Copy</button>' +
//                 '<div class="app-info text-muted">' + this.time + ' - ' + this.application + '</div>' +
//                 '<div class="full-message">' + this.message + '</div>' +
//                 '</div>';
//                 '</div>' +
//                 '</div>';
//
// /*
//             '<div class="panel-body-inner severity-' + this.severity + '" id="entry-' + this.elementCount + '">' +
//                        '<span class="loglevel"></span>' +
//                        '<button class="btn-copy" title="Copy to clipboard">Copy</button>' +
//                        '<div class="app-info text-muted">' + this.time + ' - ' + this.application + '</div>' +
//                        '<div class="full-message">' + this.message + '</div>' +
//                    '</div>';
// */
//         }
//         this.getRowAsHtml = function () {
//             return '<div class="panel panel-top panel-default">' +
//                        '<div class="panel-heading">' + this.getHeader() + '</div>' +
//                        '<div class="panel-body">' + this.getBody() + '</div>' +
//                    '</div>';
//
// /*
//             '<div class="panel panel-default" data-loglevel="' + this.severity + '">' +
//                        '<div class="panel-heading js-toggle severity-' + this.severity + '">' +
//                            this.getHeader() +
//                        '</div>' +
//                        '<div class="panel-body">' +
//                            this.getBody() +
//                        '</div>' +
//                    '</div>';
// */
//         }
//     }

    utility.showLoader = function () {
        // timer to set actions after request
        if (timer) {
            clearTimeout(timer); //cancel the previous timer.
            timer = null;
        }
        timer = setTimeout(function() {
            $loader.removeClass('active');
        }, 800);
    };

    utility.handleScrolling = function () {
        if(checkAutoscroll() === 'true') {
            //$("html, body").stop().animate({ scrollTop: $(document).height() }, 20);
            window.scrollTo(0,document.body.scrollHeight);
        }
    };
    
    // utility.printEntry = function (logEntry) {
    //     var d = $('<div class="panel-wrapper"></div>')
    //             .append(logEntry.getRowAsHtml()),
    //         lastPanel = $('.panel-top').last();
    //
    //     d = handleOutput($(d));
    //
    //     // if new logger is similar to last one, group within one panel
    //     if (this.lastLogger.length && logEntry.logger.length && this.lastLogger == logEntry.logger) {
    //         var newPanel = $(d).find('.panel');
    //         //console.log(newPanel);
    //
    //         lastPanel.find('.panel-body').append(logEntry.getBody());
    //         activateCopyButton(lastPanel.find('.btn-copy').last());
    //
    //         // Update main Loglevel for panel (set to most serious one)
    //         if (typeof this.lastLogLevel != 'undefined' && typeof logEntry.severity != 'undefined' && this.lastLogLevel > logEntry.severity) {
    //             this.lastLogLevel = logEntry.severity;
    //
    //             lastPanel.data('loglevel', logEntry.severity);
    //             lastPanel.find('.panel-heading')
    //                 .removeClass('severity-' + this.lastLogLevel)
    //                 .addClass('severity-' + logEntry.severity);
    //         }
    //     } else {
    //         this.lastLogLevel = logEntry.severity;
    //         $(output).append(d);
    //     }
    //     this.lastLogger = logEntry.logger;
    //
    //     this.handleScrolling();
    //     this.showLoader();
    // };

    // utility.print = function(message, logger) {
    //     var d = $('<div class="panel-wrapper"></div>')
    //                 .append(message);
    //
    //     d = handleOutput($(d));
    //     $(output).append(d);
    //     lastLogger = logger;
    //
    //     this.handleScrolling();
    //     this.showLoader()
    // };

    utility.printStartMessage = function() {
        // default startmessage (onopen)
        function getStartText() {
            var values = [
                "Ready to log!",
                "Jauchzet frohlogget!",
                "Let's log 'n' lol!",
                "Let’s log!",
                "Log the house, dudes!",
                "Log it out!",
                "Oh my log!",
                "Log me Amadeus!",
                "Ready to log!",
                "Log log login' on heaven's door...",
                "Log at that:",
                "You log nice today!",
                "Log out now!",
                "Log n'feel.",
                "Log it up!",
                "Get your log post!",
                "Log a'rhythm.",
                "To be logged!",
                "Die or log!",
                "Fifty Shades of Log"
            ];
            return values[Math.floor(Math.random() * values.length)];
        }

        var row = '<div class="panel-heading severity-10">' +
            '<div class="panel-title"><span class="loglevel"></span><b>' + getStartText() + '</b></div>' +
            '</div>';

        //utility.print(row, "");

        var item = document.createElement("div");
        item.className = 'panel panel-default';
        item.innerHTML = row;
        document.getElementById("output").appendChild(item);
    };


    // test native js print to document without function
    utility.printToDocument = function (message) {
        var text =
                '<div class="panel-inner">' +
                    '<div class="panel-heading severity-' + message.severity + '"><span class="loglevel"></span>' +
                        message.logger + '<br>' +
                        message.application + '<br>' +
                    '</div><div class="panel-body">' +
                        'request_uri: ' + message.request_uri + '<br>' +
                        'request_id: ' + message.request_id + '<br>' +
                        'time: ' + message.time + '<br>' +
                        'severity: ' + message.severity + '<br>' +
                        'severity_text: ' + message.severity_text + '<br>' +
                        'context: ' +  JSON.stringify(message.context) +
                        '<br><br>' + message.message +
                    '</div>' +
                '</div>',
            item = document.createElement("div");
        item.className = 'panel panel-default';
        item.innerHTML = text;
        document.getElementById("output").appendChild(item);
    };

    window.addEventListener("load", function(evt) {
        var output = document.getElementById("output"),
            ws,
            host = window.location.hostname;


        ws = new WebSocket('ws://' + host + ':' + port + endPoint);
        ws.onopen = function(evt) {
            utility.printStartMessage();
        };
        ws.onclose = function(evt) {
            var data = {
                    "message" : "Connection closed",
                    "severity" : 7
                };
                //entry = new LogEntry(data);
            //utility.printEntry(entry, data.logger);
            utility.printToDocument(data);
            utility.handleScrolling();
            ws = null;
        };
        ws.onmessage = function(evt) {
            var data = JSON.parse(evt.data);
                //entry = new LogEntry(data);
            $loader.addClass('active');
            //utility.printEntry(entry, data.logger);
            utility.printToDocument(data);
            utility.showLoader();
            utility.handleScrolling();
        };
        ws.onerror = function(evt) {
            var data = {
                "message" : "ERROR: " + evt.data,
                "severity" : 2
            };

            //var entry = new LogEntry(data);
            //utility.printEntry(entry, data.logger);
            utility.printToDocument(data);
            utility.handleScrolling();

        };
    });

    window.addEventListener("beforeunload", function(evt) {
        // delete the cookie
        document.cookie = "logbook=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    });

    $body.find('#btn-up').on('click', function(){
        $('html, body').animate({scrollTop:0});
        return false;
    });

    $body.find('#btn-down').on('click', function(){
        $('html, body').animate({scrollTop:$(document).height()});
        return false;
    });


    if(checkAutoscroll() === 'true') {
        $toggleAutoscroll.addClass('active');
    } else {
        $toggleAutoscroll.removeClass('active');
    }

    $toggleAutoscroll.on('click', function() {
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
