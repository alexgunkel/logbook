# LogBook

1. [Installation](documentation/Install.md)
2. [Web-Frontend](documentation/Frontend.md)
3. [Logger-Implementation](documentation/LoggerApi.md)


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

![LogBook: logging screen](documentation/img/logbook-screenshot.png?raw=true "LogBook: logging screen")