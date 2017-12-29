# Logger API

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

## Existing packages

Logger packages are available for:
* *PHP*: [logbook-php](https://github.com/axel-kummer/logbook-php)

## Path
Loggers should only submit log events for clients with "logbook"-cookie.
This cookie contains the identifier by which LogBook decides whom to
send the log messages. Messages are sent as POST-requests in JSON to

    <domain>/api/v1/logbooks/<LogBook_id>/logs

That is the same path as where the web frontend will request the logs. You
can change this path by setting the ``API_ROOT_PATH``-variable:

    export API_ROOT_PATH=/alternative/root/path

## Body
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
| 6      | info          | informational messages           |
| 7      | debug         | debug-level messages             |

At the moment, string recognition is case-sensitive and only lower case
strings are recognized.

*LogBook* will accept your logs and send back a *200 Status OK* response
as soon as it can recognize the body of the request as valid JSON. That
does *not* mean that it can interpret it. Your message might still be
invalid.

## Header
Information about the *application* that is logging the event(s) should
be put into the header information. The header keys have *LogBook*
as prefix. At the moment there are three header fields:

    LogBook-App-Identifier: MyMicroService
    LogBook-Logger-Name: MyLoggerInstance
    LogBook-Request-URI: https://my.lovely.app

These informations are not required but highly recommended. You should
use them to make your logs more helpful.
