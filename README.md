# LogBook

*LogBook* is a server-side application that lets you debug your web application by showing
the log messages related to your web request. It stands in the tradition of FirePHP and
ChromeLogger but extracts the core functionality into a server-side log-dispatcher.
The reason is that for larger amounts of data, header fields are an inappropriate way to send
all the logs from a web-application to the client. When you have a small application with
few log events, it is quite easy to send those data with your HTTP-header. But their are hard
limits that prevent you from sending as many log-data as a larger application will provide. *LogBook*
collects all those data in a server-side application just to provide this information to the
developer in an JavaScript-based frontend-application.

A *logbook* is considered to be the whole information about a session. It's like keeping a book about
a journey. Every session consists of the whole interaction between a client and an application.
The logbook collects the information that is considered relevant from a developer's prospective.

The LogBook listens for log messages that are sent via POST-requests. The default format is JSON. Messages are
identified by an identifier that is created by the web frontend application and stored in a cookie named "logbook".
When a client has a logbook-cookie then the logger should send messages to the logbook-server.

The logging requires the following information:
* The *logbook-identifier* identifies the client session.
* A *timestamp* which indicates the exact time when the event was created.
* The *severity-level* indiciates the kind of event, e.g. whether it's a fatal error, an information
or just a debug-message.
* A *log message* which contains the basic information of the event.
* A *context* might contain some further information about the logged event.
* information about the part of the application that created the log event. This is especially
relevant when it comes to micro-services, where it might be of great intereset to know which
part of your application is responsible for something. Therefore, some *application-identifier"
is required.

## Installation
You can easily install LogBook via Go CLI:

    go get github.com/alexgunkel/logbook
    cd ${GOPATH}/src/github.com/alexgunkel/logbook
    go install

## Usage
Start LogBook by typing

    logbook

*LogBook* by default listens on port 8080.

## Logger API
### Path
Loggers should only submit log events for clients with "logbook"-cookie. This cookie contains the identifier by which
LogBook decides whom to send the log messages. Messages are sent as POST-requests in JSON to

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

### Header
Information about the *application* that is logging the event(s) should be put into the header
information. The header keys have *LogBook* as prefix. At the moment there are three header
fields:

    LogBook-App-Identifier: MyMicroService
    LogBook-Logger-Name: MyLoggerInstance
    LogBook-Request-URI: https://my.lovely.app

These informations are not required but highly recommended. You should use them to make your
logs more helpfull

### Existing packages

Logger packages are available for:
* *PHP*: [logbook-php](https://github.com/axel-kummer/logbook-php)

## Web Frontend
This package is already shipped with a small and simple log viewer.
This web frontend will available under the url

    <domain>/logbook

For a more comfortable presentation of your LogBook you might want to
install the AngularJS-application
[logbook-frontend](https://github.com/XenosEleatikos/logbook-frontend)

The frontend is responsible for redirecting the client to the log-messages-page

    <domain>/logbook/<LogBook_id>/logs

The *EcmaScript* application then establishes a websocket-connection and displays the log messages related to the client.

## Orchestration
We recommend orchestration via Nginx-Proxy.
