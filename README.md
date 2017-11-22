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

## Usage
To insert use it you just have to create a simple go project with a main package and a main function
that contains the following lines:

    import "github.com/alexgunkel/logbook"
    // ...
    LogApplication := logbook.Application()
    LogApplication.Run()

Then start your app and make sure that relevant services are able to reach you logbook on port 80.


## Logger API
Loggers should only submit log events for clients with "logbook"-cookie. This cookie contains the identifier by which
LogBook decides whom to send the log messages. Messages are sent as POST-requests in JSON to

   <domain>/logbook/<LogBook_id>/logs

Logs have to be sent as JSON-objects that contain:
* the *time* (integer, required),
* the *message* (string, required),
* the *severity* (string or integer, required),
* a *logger_name* (string, optional),
* a *context* (array, optional).

So, an example json might look like this:

    {
        "time": 1511390786,
        "message": "",
        "severity": 2,
        "logger_name": "example-logger"
    }

Information about the *application* that is logging the event(s) should be put into the header
information. The header key is _LogBook-App-Identifier_:

    LogBook-App-Identifier: MyMicroService

Logger packages are available for:
* [PHP](https://github.com/axel-kummer/logbook-php)

## Web Frontend
The web frontend is available under the url

    <domain>/logbook

The frontend is responsible for redirecting the client to the log-messages-page

    <domain>/logbook/<LogBook_id>/logs

The EcmaScript application then establishes a websocket-connection and displays the log messages related to the client.