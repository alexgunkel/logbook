package main

// we should test that logger and frontend can be initiated
// together without conflicts. conflicts can arise when there
// are conflicting routes defined, e.g. frontend tries to
// serve the same route like the websocket-service (which starts
// as GET-request as well.
