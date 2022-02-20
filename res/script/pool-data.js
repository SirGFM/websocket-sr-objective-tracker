let pool = function() {
	/** Pooling function. */
    let _interval = null;

	/**
	 * "Handles" anything that goes wrong.
	 *
	 * @param{e} The event
	 */
	let _onError = function(e) {
		console.log(e);
	}

	/**
	 * Open a new connection. After this, the caller only has to call xhr.send(obj).
	 *
	 * @param{url} HTTP address of the connecting server.
	 * @param{mode} HTTP Method for the connection (e.g., PUT or GET).
	 * @param{onLoad} Callback executed when the request is fully sent/received.
	 *                onLoad must receive a single argument (the response).
	 * @return A newly set up XMLHttpRequest.
	 */
	let _openConn = function(url, mode, onLoad=null, async=true) {
		let xhr = new XMLHttpRequest();
		xhr.open(mode, url, async);

		if (onLoad) {
			let cb = function (e) {
				let res = e.target.response;
				onLoad(res);
			};
			xhr.addEventListener('loadend', cb);
		}
		xhr.addEventListener('error', _onError);

		xhr.setRequestHeader('Accept', 'application/json');

		return xhr;
	}

	/**
	 * Retrieve a JSON object from the server.
	 *
	 * @param{url} HTTP address of the connecting server.
	 * @param{onLoad} Callback executed when the request is fully received.
	 *                onLoad must receive a single argument (the response).
	 */
	function _getData(url, onLoad) {
		let xhr = _openConn(url, 'GET', onLoad);
		xhr.send(null);
	}

	/**
	 * Stop pooling and updating the current page.
	 */
	function _stopTracking() {
		if (_interval !== null) {
            clearInterval(_interval);
			_interval = null;
		}
	}

	/**
	 * Start pooling and updating the current page.
	 *
	 * @param{uri} HTTP address of the tracking resource.
	 * @param{fps} How often per second should the resource be pooled.
	 */
	function _startTracking(uri, fps) {
		_stopTracking();

		let fn = function() {
			_getData(uri, function(ev) {
				try {
					let res = JSON.parse(ev);

					tracker.setObject(res);
				} catch (e) {
					console.log(`Failed to parse "${uri}"`);
					console.log(e);
				}
			});
		}
		_interval = setInterval(fn, 1000 / fps);
	}

	return {
		'getData': _getData,
		'startTracking': _startTracking,
		'stopTracking': _stopTracking,
	}
}()
