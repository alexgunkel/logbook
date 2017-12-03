# Web Frontend
A *logbook* is considered to be the information about a session.
It's like keeping a book about a journey. The logbook collects
the information that is considered relevant from a developer's
prospective. To view the stream of information you need a web-frontend
application.

## Built-in Frontend-Application
This package is already shipped with a small and simple log viewer.
This web frontend listens on port `:8080` by default. It is important,
that *LogBook* has the same domain as the application you are
monitoring. Otherwise your browser won't handle the cookie correctly.

## LogBook-Frontend
For a more comfortable presentation of your LogBook you might want to
install the AngularJS-application
[logbook-frontend](https://github.com/XenosEleatikos/logbook-frontend)

## Web-Frontend API
You are free to use some custom frontent-application instead. This application
can be served by the built-in web-server. It must constist of a mandatory
`Index.html`-file and may include further files (e.g. for css and JavaScript).
The `Index.html`-file will be treated as a template, which meens that you can
use some variables that are automatically replaced by the corresponding values.
Those variables are:

| variable name     | content |
|-------------------|---------|
| `{{.PathToStatic}}`   | is substituted by the relative path from where static files can be downloaded. |
| `{{.Uri}}`        | contains the link to the websocket |
| `{{.Identifier}}` | contains the logbook-id |

    

To use your custom frontend-application, just start LogBook with the respective
directory path set as environmental variable `STATIC_APP`, e.g.:

    export STATIC_APP=$(pwd)/html
    logbook

The frontend receives its information from a websocket that is established
under the URL

    <domain>/api/<api_version>/logbooks/<LogBook_id>/logs
	
The path is available in your `Index.html`file with the template variable `{{.Uri}}`.

When the session is established, a cookie is set that contains the logbook-id.
Furthermore, you can use some variables in your *Index.html* that are
automatically substituted by corresponding values. Those variables are:

	