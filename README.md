# LogBook

*LogBook* is a server-side application that lets you debug your web
application by showing the log messages related to your web request.
It stands in the tradition of
[FirePHP](https://addons.mozilla.org/de/firefox/addon/firephp/) and
[ChromeLogger](https://craig.is/writing/chrome-logger)
but extracts the core functionality into a server-side log-dispatcher.
The reason is that for larger amounts of data, header fields are an
inappropriate way to send all the logs from a web-application to the client.
When you have a small application with few log events, it is quite easy
to send those data with your HTTP-header. But their are hard limits
that prevent you from sending as many log-data as a larger application
will provide. *LogBook* collects all those data in a server-side
application just to provide this information to the developer in an
JavaScript-based frontend-application.

A *logbook* is considered to be the whole information about a session.
It's like keeping a book about a journey. Every session consists of the
whole interaction between a client and an application. The logbook
collects the information that is considered relevant from a developer's
prospective.

The LogBook listens for log messages that are sent via POST-requests.
The default format is JSON. Messages are identified by an identifier
that is created by the web frontend application and stored in a cookie
named "logbook". When a client has a logbook-cookie then the logger
should send messages to the logbook-server.

The logging requires the following information:
* The *logbook-identifier* identifies the client session.
* A *timestamp* which indicates the exact time when the event was created.
* The *severity-level* indiciates the kind of event, e.g. whether it's a
fatal error, an information or just a debug-message.
* A *log message* which contains the basic information of the event.
* A *context* might contain some further information about the logged event.
* information about the part of the application that created the log event.
This is especially relevant when it comes to micro-services, where it
might be of great intereset to know which part of your application is
responsible for something. Therefore, some *application-identifier"
is required.

## Installation
You can easily install LogBook via Go CLI:

    go get github.com/alexgunkel/logbook
    cd ${GOPATH}/src/github.com/alexgunkel/logbook
    go install

## Requirements
This program requires Golang at least in version 1.6 (the least version
that's tested by me. For the tests to run you need at least version 1.7.

## Usage
Start LogBook by typing

    logbook

*LogBook* by default listens on port 8080. To change this just set the
environmental variable PORT to any other value. Likewise you can configure
the hostname:

    PORT=80 HOST=127.0.0.1 logbook

will change the port to 80 and the hostname to 127.0.0.1.

## Logger API
### Path
Loggers should only submit log events for clients with "logbook"-cookie.
This cookie contains the identifier by which LogBook decides whom to
send the log messages. Messages are sent as POST-requests in JSON to

    <domain>/logbook/<LogBook_id>/logs

That is the same path as where the web frontend will request the logs.

### Body
Logs have to be sent as JSON-objects that contain:
* the *time* (integer, required),
* the *message* (string, required),
* the *severity* (string or integer, required),
* a *context* (array, optional).

So, an example json might look like this:

    {
        "time": 1511390786,
        "message": "",
        "severity": 2,
        "logger_name": "example-logger"
    }

The severity has to be a string or a number in accordance to
[RFC 5424](https://tools.ietf.org/html/rfc5424).
As a number, it must be one of the
values between *0* (highest severity, "*emergeny*") and *7* (lowest severity,
"debug"). The relqtionship between numbers and strings is as follows:


| number | string        | explanation                      |
|--------|---------------|----------------------------------|
| 0      | emergency     | system is unusable               |
| 1      | alert         | action must be taken immediately |
| 2      | critical      | critical conditions              |
| 3      | error         | error conditions                 |
| 4      | warning       | warning conditions               |
| 5      | notice        | normal but significant condition |
| 6      | informational | informational messages           |
| 7      | debug         | debug-level messages             |

At the moment, string recognition is case-sensitive and only lower case
strings are recognized.

*LogBook* will accept your logs and send back a *200 Status OK* response
as soon as it can recognize the body of the request as valid JSON. That
does *not* mean that it can interpret it. Your message might still be
invalid.

### Header
Information about the *application* that is logging the event(s) should
be put into the header information. The header keys have *LogBook*
as prefix. At the moment there are three header fields:

    LogBook-App-Identifier: MyMicroService
    LogBook-Logger-Name: MyLoggerInstance
    LogBook-Request-URI: https://my.lovely.app

These informations are not required but highly recommended. You should
use them to make your logs more helpful.

### Existing packages

Logger packages are available for:
* *PHP*: [logbook-php](https://github.com/axel-kummer/logbook-php)

## Web Frontend
This package is already shipped with a small and simple log viewer.
This web frontend will available under the url

    <domain>/logbook

To serve some static files instead, just start LogBook with the respective
directory path set as environmental variable STATIC_APP, e.g.:

    export STATIC_APP=$(pwd)/html
    logbook

For a more comfortable presentation of your LogBook you might want to
install the AngularJS-application
[logbook-frontend](https://github.com/XenosEleatikos/logbook-frontend)

The frontend receives its information from a websocket that is established
under the URL

    <domain>/logbook/api/<LogBook_id>/logs

## Docker and Orchestration
There will be a Docker image available which you can get with

    docker pull alexandergunkel/logbook

For further information visit the
[public repository](https://hub.docker.com/r/alexandergunkel/logbook/)
at [Docker Hub](https://hub.docker.com)
