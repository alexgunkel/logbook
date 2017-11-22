# LogBook

LogBook is a server-side application that lets you debug your web application by showing the log messages related
to your web request.

The LogBook listens for log messages that are sent via POST-requests. The default format is JSON. Messages are
identified by an identifier that is created by the web frontend application and stored in a cookie named "logbook".
When a client has a logbook-cookie then the logger should send messages to the logbook-server.

## Installation
You can easily install LogBook via Go:

    go get github.com/alexgunkel/logbook

## Usage
To insert use it you just have to create a simple go project with a main package and a main function
that contains the following lines:

    import "github.com/alexgunkel/logbook"
    // ...
    logBook := logbook.Application()
    logBook.Run()

## Logger API
Loggers should only submit log events for clients with "logbook"-cookie. This cookie contains the identifier by which
LogBook decides whom to send the log messages. Messages are sent as POST-requests in JSON to

   <domain>/logbook/<client_id>/logs

Logger packages are available for:
* [PHP](https://github.com/axel-kummer/logbook-php)

## Web Frontend
The web frontend is available under the url

    <domain>/logbook

The frontend is responsible for redirecting the client to the log-messages-page

    <domain>/logbook/<client_id>/logs

The EcmaScript application then establishes a websocket-connection and displays the log messages related to the client.