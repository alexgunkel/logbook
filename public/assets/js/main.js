$(function() {
    var utility = {},
        elementCount = 0,
        elementCounter = 0,
        timer = null,
        $body = $('body'),
        $loader = $('.loader'),
        $toggleAutoScroll = $body.find('.js-toggle-autoscroll'),
        $toggleSeverityFilter = $body.find('#nav-filter .js-toggleSeverityFilter'),
        port = $body.data('port'),
        endPoint = $body.data('endpoint');

    // update page title
    document.title = "LogBook » " + window.location.protocol + '//' + window.location.hostname;

    // Utility variables
    utility.requestGroups = [];
    utility.lastLogger = "";
    utility.lastLogLevel = 100;

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

    function checkAutoScroll() {
        return getCookie("autoscroll");
    }

    if(checkAutoScroll() === 'true') {
        $toggleAutoScroll.addClass('active');
    } else {
        $toggleAutoScroll.removeClass('active');
    }

    function initSeverityFilters($levels) {
        $.each($levels , function(index, hideLogLevel) {

            if(getCookie('hideLogLevel-' + hideLogLevel) === 'true') {
                $body.addClass("hide-log-level-" + hideLogLevel);
                $('#nav-filter').find('[data-severity="' + hideLogLevel + '"]').removeClass('active');
            }

        });
    }

    // Check what filters are set via cookie
    initSeverityFilters(['info','warning','error']);

    $toggleAutoScroll.on('click', function() {
        if($(this).hasClass('active')) {
            setCookie("autoscroll", false);
            $(this).removeClass('active');
        } else {
            setCookie("autoscroll", true);
            $(this).addClass('active');
        }
        return false;
    });

    $toggleSeverityFilter.on('click', function() {
        var hideLogLevel = $(this).data('severity');

        if($(this).hasClass('active')) {
            setCookie("hideLogLevel-" + hideLogLevel, true);
            $(this).removeClass('active');
            $body.addClass("hide-log-level-" + hideLogLevel);
        } else {
            setCookie("hideLogLevel-" + hideLogLevel, false);
            $(this).addClass('active');
            $body.removeClass("hide-log-level-" + hideLogLevel);
        }


        return false;
    });

    var handleOutput = function () {

        $(document).on("click", '.js-toggle', function(event) {
            event.preventDefault();
            event.stopPropagation();

            var button = $(this),
                bodyID = button.attr('href');
            button.find('.glyphicon').toggleClass('glyphicon-zoom-out glyphicon-zoom-in');
            button.parents('.panel').find('.panel-body').toggleClass('hidden');
        });

        $(document).on("click", '.btn-copy', function(event) {
            event.preventDefault();
            event.stopPropagation();

            var button = $(this),
                text   = button.next('.full-message');
                input  = $('<input>');

            button.parent().append(input);
            input.val(text.text()).select();
            document.execCommand('copy');
            input.remove();
        });

        $(document).on("click", '.js-toggle-context', function(event) {
            event.preventDefault();
            event.stopPropagation();

            var button = $(this);

            button.find('span').toggleClass('glyphicon-chevron-down').toggleClass('glyphicon-chevron-up');
            button.next().next('.context').toggleClass('hidden');
        });

        $(document).on("click", '.js-toggle-body', function(event) {
            event.preventDefault();
            event.stopPropagation();

            var button = $(this);

            button.find('span').toggleClass('glyphicon-chevron-down').toggleClass('glyphicon-chevron-up');
            button.next('.full-message').toggleClass('content-hidden');
        });
    };

    function LogEntry(data) {

        if(Object.keys(data).length > 0) {
            this.logger = data.logger;
            this.severity = (data.severity) ? data.severity : '7';
            this.message = data.message;
            this.time = data.time;
            this.application = data.application;
            this.request_uri = data.request_uri;
            elementCount = ++elementCounter;
            toggleContent = '';
            hideContent = '';
            if(this.message.length > 1000) {
                toggleContent = '<span class="js-toggle-body pull-right"><span class="glyphicon glyphicon-chevron-down"></span></span>';
                hideContent = ' content-hidden';
            }

            this.getBody = function () {
                return '<div class="panel-body-inner severity-' + this.severity + ' logentry" id="entry-' + elementCount + '">' +
                    '<span class="loglevel"></span>' +
                    '<div class="loglevel-text"><small><b>' + (data.time - utility.requestGroups[data.request_id].startTime) + 's</b> ' + data.severity_text + '</small></div>' +
                    '<button class="btn-copy" title="Copy to clipboard">Copy</button>' + toggleContent +
                    '<div class="full-message' + hideContent + '">' + this.message + '</div>' + utility.getContextAsString(data.context) +
                    '</div>';
            };
        }
    }

    utility.getContextAsString = function(context)
    {
        var output = "";

        if (context.data.length !== 0) {
            output = '<span class="js-toggle-context pull-right"><span class="glyphicon glyphicon-chevron-down"></span>Show Context</span>' +
                '<div class="clearfix"></div>' +
                '<div class="context hidden"><table>';
            $.each(context.data, function (index, value) {
                if (index === 'backtrace') {
                    value = utility.getBackTraceAsString(value);
                }

                output += '<tr><td class="table-label">' + index + '</td><td>' + JSON.stringify(value) + '</td></tr>'
            });
            output += '</table></div>';
        }

        return output;

    };

    utility.getBackTraceAsString = function(backtrace)
    {
        var string = "<table>";

        $.each(backtrace, function (index, value) {
            string += "<tr><td>";
            string += (value.file) ? value.file + ":" : "";
            string += (value.line) ? value.line + " " : "";
            string += "<br>";
            string += (value.class) ? value.class : "method";
            string += (value.type) ? value.type  : "::";
            string += (value.function) ? value.function : "";
            string += "</tr></td>";
        });

        string += "</table>";

        return string;
    };

    utility.showSlider = function () {
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
        if(checkAutoScroll() === 'true') {
            $("html, body").stop().animate({ scrollTop: $(document).height() }, 20);
        }
    };

    utility.updateLogLevel = function (currentLogLevel, lastPanel) {
        // Update main Loglevel for panel (set to most serious one)
        if (typeof this.lastLogLevel != 'undefined' && typeof currentLogLevel != 'undefined' && this.lastLogLevel > currentLogLevel) {
            this.lastLogLevel = currentLogLevel;
            this.printLogLevel(lastPanel);
        }
    };

    utility.printLogLevel = function ($lastPanel) {
        // Update main Loglevel for panel (set to most serious one)
        console.log('panel updated', $lastPanel);
        $lastPanel.attr('data-loglevel', this.lastLogLevel);
    };

    utility.printEntry = function (logEntry, data, output) {
        var requestGroup = this.getRequestGroup(data, output),
            $lastPanel;
        if (this.lastLogger !== data.logger) {
            this.lastLogLevel = 99;
            requestGroup.find('.panel-body').append('<div class="panel-body-wrap" data-loglevel=""><div class="loggername panel-body-inner-heading">' + data.logger + '</div></div>')
        }
        $lastPanel = requestGroup.find('.panel-body .panel-body-wrap').last();
        $lastPanel.append(logEntry.getBody());

        this.updateLogLevel(data.severity, $lastPanel);
        this.lastLogger = data.logger;
        this.handleScrolling();
        this.showSlider();
    };

    utility.getRequestGroup = function(data, output) {
        if (this.requestGroups[data.request_id]) {
            return this.requestGroups[data.request_id];
        }

        if(typeof data.request_uri != 'undefined' && data.request_uri.length) {
            // cut long request linktexts
            if(data.request_uri === "/") {
                requestLinkText = window.location.hostname;
            } else if(data.request_uri.length > 130) {
                requestLinkText = window.location.hostname + data.request_uri.substring(0,130) + '...';
            } else {
                requestLinkText = window.location.hostname + data.request_uri
            }
            var requestLink = '<a href="' + data.request_uri + '" title="' + data.request_uri + '">' + requestLinkText + '</a>';
        }

        this.requestGroups[data.request_id] = $('<div class="panel-wrapper"><div class="panel panel-default" id="' + data.request_id + '">' +
            '<div class="panel-heading js-toggle"><b>' + requestLink + '</b><br>' +
                '<span class="small">' + data.request_id + '</span>' +
                '<a class="js-toggle" href="#entry-' + data.request_id + '">' +
                    '<span class="glyphicon glyphicon-zoom-out" title="show more"></span>' +
                '</a>' +
            '</div>' +
            '<div class="panel-body" id="#entry-' + data.request_id + '"></div> ' +
            '</div></div>');

        this.requestGroups[data.request_id].startTime = data.time;

        $(output).append(this.requestGroups[data.request_id]);

        return this.requestGroups[data.request_id];
    };

    utility.print = function(message, logger) {
        var d = $('<div class="panel-wrapper"></div>')
                    .append(message);

        $(output).append(d);

        utility.lastLogger = logger;

        this.handleScrolling();
        this.showSlider();
    };

    utility.printStartMessage = function() {
        // default startmessage (onopen)
        function getStartText() {
            var values = [
                "Ready to log!",
                "Jauchzet frohlogget!",
                "Log 'n' lol!",
                "Have a log!",
                "Log the house, dudes!",
                "The log is ticking",
                "Log it out!",
                "Oh my log!",
                "Log me Amadeus!",
                "Ready to log!",
                "Login' on heaven's door...",
                "Log at that:",
                "You log nice today!",
                "Log out now!",
                "Log n’feel.",
                "Log it up!",
                "Get your log post!",
                "Log’a’rhythm.",
                "To be logged!",
                "Die or log!",
                "Fifty Shades of Log",
                "We're gonna log",
                "Logs slay time",
                "Anti logwise"
            ];
            return values[Math.floor(Math.random() * values.length)];
        }

        var row = '<div class="panel panel-default">' +
            '<div class="panel-heading  js-toggle severity-10">' +
            '<div class="panel-title"><span class="loglevel"></span><b>' + getStartText() + '</b></div>' +
            '</div>' +
            '</div>';

        utility.print(row, "");
    };

    $body.find('#btn-up').on('click', function(){
        $('html, body').animate({scrollTop:0});
        return false;
    });

    $body.find('#btn-down').on('click', function(){
        $('html, body').animate({scrollTop:$(document).height()});
        return false;
    });

    window.addEventListener("load", function(evt) {
        var output = document.getElementById("output"),
            ws,
            host = window.location.hostname;

        ws = new WebSocket('ws://' + host + ':' + port + endPoint);

        ws.onopen = function(event) {
            utility.printStartMessage();
        };
        ws.onclose = function(event) {
            var data = {
                    "message" : "Connection closed",
                    "severity" : 7
                },
                entry = new LogEntry(data);
            utility.printEntry(entry, data.logger);
            ws = null;
        };
        ws.onmessage = function(event) {
            var data = JSON.parse(event.data),
                entry = new LogEntry(data);
            $loader.addClass('active');
            utility.printEntry(entry, data, output);
        };
        ws.onerror = function(evt) {
            var data = {
                "message" : "ERROR: " + evt.data,
                "severity" : 2
            };

            var entry = new LogEntry(data);
            utility.printEntry(entry, data.logger);
        };
    });

    window.addEventListener("beforeunload", function(evt) {
        // delete the cookie
        document.cookie = "logbook=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    });

    handleOutput();
});
