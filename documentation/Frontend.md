# Web Frontend

A *logbook* is considered to be the whole information about a session.
It's like keeping a book about a journey. Every session consists of the
whole interaction between a client and an application. The logbook
collects the information that is considered relevant from a developer's
prospective.

This package is already shipped with a small and simple log viewer.
This web frontend will available under the url

    <domain>/logbook/app

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
	
When the session is established, a cookie is set that contains the logbook-id.
Furthermore, you can use some variables in your *Index.html* that are
automatically substituted by corresponding values. Those variables are:

    {{.BaseHref}} is substituted by the relative path where assets can be downloaded.
	{{.Uri}} contains the link to the websocket
	{{.Identifier}} contains the logbook-id
	