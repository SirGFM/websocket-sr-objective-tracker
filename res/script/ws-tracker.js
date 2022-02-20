let ws_tracker = function() {
	/** The current WebSocket connection. */
	let _ws = null;

	/** Reconnection function, in case the connection fails. */
	let _reconnectFn = null;

	/**
	 * Handle receiving messages from the server.
	 */
	function _wsRecv(ev) {
		let res = JSON.parse(ev.data);

		try {
			id = res['id'];
			value = res['value'];

			tracker.setValue(id, value);
		} catch (e) {
			console.log(`Failed to parse the WebSocket response"`);
			console.log(e);
		}
	}

	/**
	 * Stop pooling and updating the current page.
	 */
	function _stopTracking() {
		if (_ws !== null) {
			_reconnectFn = null;

            _ws.close();
			_ws = null;
		}
	}

	/**
	 * Start pooling and updating the current page.
	 *
	 * @param{uri} Address of the WebSocket tracker.
	 */
	function _startTracking(uri) {
		_stopTracking();

		_reconnectFn = function() {
			let retry = true;
			let proto = 'ws';
			if (window.location.protocol.startsWith('https')) {
				proto = 'wss';
			}

			let reconnect = function() {
				/* Retry after some time. */
				if (retry && _reconnectFn != null) {
					console.log(`Failed to establish/maintain a WebSocket connection to "${uri}". Retrying...`);
					setTimeout(_reconnectFn, 1000);
					retry = false;
				}
			}

			let ws = new WebSocket(proto + '://' + window.location.host + uri);
			ws.addEventListener('open', function() {
				console.log(`New WebSocket connection to "${uri}"`)
				_ws = ws;
			});
			ws.addEventListener('message', _wsRecv);
			ws.addEventListener('close', reconnect);
			ws.addEventListener('error', reconnect);
		}

		_reconnectFn();
	}

	return {
		'startTracking': _startTracking,
		'stopTracking': _stopTracking,
	}
}()
