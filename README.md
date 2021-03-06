# WebSocket-based Speedrun Objective Tracker

A generic tracker viewable in your browser that can be updated automatically by other tools!

## Quick start

* Requirements: [Golang](https://go.dev/dl/)
	* This was built and tested with Go 1.17.6, but it should work with any more recent version.

To cross-compile for Windows from Linux, run:

```bash
go get -d github.com/gorilla/websocket
GOOS=windows go build .
```

To build natively, simply run

```bash
go get github.com/gorilla/websocket
go build .
```

## Core concept

This tracker has a few components:

* **HTTP server** 
	* Allows updating data with `POST`/`DELETE` requests
	* Allows viewing the server's state in a web page
* **Web Socket**
	* (this is actually part of the http server)
	* Allows a simple, bi-directional and continuous communication between the page and the server

A game should have its own custom HTML page with all the resources needed to track the objectives in the game. Each resource SHALL have a unique ID, which will be used to show/highlight the resources that were obtained. Something like:

```html
<head>
	<!-- ... -->

	<!-- Script that should be generic enough to work for any page/game. -->
	<script type="text/javascript" src="/script/ws-tracker.js"></script>
</head>

<!-- This mess uses a image as the background, so the img elements may be
  - placed on top of that... Hack-y? Yes... but it works. ¯\_(ツ)_/¯
  -
  - Also, good luck with the CSS in this... 😬 -->
<body style='background-image: url("/img/bg.png");'>
	<img style='visibility: hidden;' id='key-item-1' src='/img/key-item-1.png' />
	<img style='visibility: hidden;' id='key-item-2' src='/img/key-item-2.png' />

	<img style='visibility: hidden;' id='optional-item' src='/img/optional-item.png' />

	<label id='textual-data'> 0 </label>
</body>
```

A application wanting to update a value would then make a HTTP request with the game as the URI and the resource (and its value) in the payload. Using `curl` to send the commands, it would be something like:

```bash
# Enable key-item-1
curl -X POST -H 'Content-type: application/json' --data '{ "id": "key-item-1", "value": true}' http://localhost:8000/tracker/dummy-game

# Enable key-item-2
curl -X POST -H 'Content-type: application/json' --data '{ "id": "key-item-2", "value": true}' http://localhost:8000/tracker/dummy-game

# Disable key-item-1
curl -X POST -H 'Content-type: application/json' --data '{ "id": "key-item-1", "value": true}' http://localhost:8000/tracker/dummy-game

# Set the textual data to "some-text"
curl -X POST -H 'Content-type: application/json' --data '{ "id": "textual-data", "value": "some-text"}' http://localhost:8000/tracker/dummy-game
```

Alternatively, the ID may be used directly for simple boolean values. In this case, a `POST` enables the resource and a `DELETE` disables it:

```bash
# Enable key-item-3
curl -X POST http://localhost:8000/tracker/dummy-game/key-item-3

# Disable key-item-3
curl -X DELETE http://localhost:8000/tracker/dummy-game/key-item-3
```

Lastly, sending a `DELETE` without an specific resource clears every resource in that game:

```bash
# Clear everything in the game
curl -X DELETE http://localhost:8000/tracker/dummy-game
```

An application that want to update the tracker automatically would need to use some HTTP binding for the language, such as [`requests`](https://docs.python-requests.org/en/latest/) for Python or [`System.Net.Http.HttpClient`](https://docs.microsoft.com/en-us/dotnet/api/system.net.http.httpclient) for C#.

## Example

To exemplify how this tracker may be used, see `res/example/index.html`, a "tracker" for rainbow-colored orbs (taken from my [LD#32 entry](https://github.com/sirgfm/ld32)).

Start by compiling and launching the server:

```
go get github.com/gorilla/websocket
go build .
./websocket-sr-objective-tracker
```

This tracker has 7 objectives:

* `red`
* `orange`
* `yellow`
* `green`
* `cyan`
* `blue`
* `purple`

You may use the simple python script `simple-request.py` to send requests to the server:

```bash
# set the red orb objective
python3 simple-request.py localhost:8000 /tracker/example/red
# remove the red orb objective
python3 simple-request.py localhost:8000 /tracker/example/red DELETE
```

To visualize the changes live, simply access its page (http://localhost:8000/example) on a browser.

## WebSocket connection

A tracker may either send GET requests to update itself manually or it may establish a WebSocket connection to get updated by the server whenever a new value is set.

The WebSocket may be accessed in `ws://<server-url>/ws-tracker/<game>`. For example, to establish a WebSocket connection to a game "example" registered to a local server running in port 8000, the WebSocket connection should be established to `ws://localhost:8000/ws-tracker/example`.

Every WebSocket message send the following object as its payload:

```
{
	"cmd": "command-as-a-string",
	"id": "resource-id-as-a-string",
	"value": any-value
}
```

Whenever a value is updated, regardless of whether its was "set" or "cleared", the server sends a payload with `cmd` set to `SET`. `id` is set to whichever resource was modified in the game and `value` is set to the value that was set. A tracker may support whichever type it desires, but the currently implemented WebSocket client only supports booleans, for setting and releasing a resource, and text, for setting textual values.

If a `DELETE` is sent to a game, clearing up every value in the server, the server sends a WebSocket message with `cmd` set to `CLEAR`. In this case, `id` and `value` should be ignored. The WebSocket client is able to handle these messages, but the tracker page must configure a clean up function by calling `ws_tracker.setClearFunction`.
