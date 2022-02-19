# WebSocket-based Speedrun Objective Tracker

A generic tracker viewable in your browser that can be updated automatically by other tools!

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

An application that want to update the tracker automatically would need to use some HTTP binding for the language, such as [`requests`](https://docs.python-requests.org/en/latest/) for Python or [`System.Net.Http.HttpClient`](https://docs.microsoft.com/en-us/dotnet/api/system.net.http.httpclient) for C#.
