window.addEventListener("paste", function(e) {
	for (var i = 0 ; i < e.clipboardData.items.length ; i++) {
		var clipboardItem = e.clipboardData.items[i];
		var type = clipboardItem.type;

		if (type.indexOf("image") != -1) {
			blob = clipboardItem.getAsFile();
			var img = new Image();
			img.src = window.webkitURL.createObjectURL(blob);
			document.querySelector('#screenshot').appendChild(img);
			document.querySelector('#instructions').style.display = 'none';
			getUploadURL(blob);
		}
	}
});

function getUploadURL(blob) {
	var xhr = new XMLHttpRequest();
	xhr.open('GET', 'http://localhost:8080/api/getuploadurl', true);

	xhr.onload = function(e) {
		if (this.status == 200) {
			upload(blob, this.responseText);
		}
	};

	xhr.send();
}

function upload(blob, uploadURL) {
	var formData = new FormData();
	var now = new Date();
	var fileName = moment().format() + '.png';
	formData.append('file', blob, fileName);

	var xhr = new XMLHttpRequest();
	xhr.open('POST', uploadURL, true);

	xhr.onload = function(e) {
		console.log(this, e);
		if (this.status == 200) {
			// TODO
		}
	};

	xhr.send(formData);
}