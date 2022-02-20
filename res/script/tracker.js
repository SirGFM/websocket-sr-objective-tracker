let tracker = function() {
	/**
	 * Update the page, showing, hiding or doing custom stuff to objects
	 * in the page.
	 *
	 * @param{key} ID of the object that should be updated.
	 * @param{value} Value of the object.
	 */
	function _setValue(key, value) {
		let obj = document.getElementById(key);
		if (!obj) {
			console.log(`Couldn't find element "${key}"!`);
			return
		}

		if (value === true) {
			obj.style.visibility = 'visible';
		}
		else if (value === false) {
			obj.style.visibility = 'hidden';
		}
		else {
			/* TODO: Handle custom values */
			console.log(`Can't handle value "${value}" for object "${key}" (yet!)`);
		}
	}

	/**
	 * Set various parameters in a page from a given object. Each field in
	 * the object should be named after an ID in the page.
	 *
	 * @param{value} Value of the object.
	 */
	function _setObject(obj) {
		for (key in obj) {
			tracker.setValue(key, obj[key]);
		}
	}

	return {
		'setValue': _setValue,
		'setObject': _setObject,
	}
} ()
