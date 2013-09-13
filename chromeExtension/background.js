function showNotification() {
	chrome.notifications.clear('straightshot', function() {
		chrome.notifications.create('straightshot', {
			type: 'basic',
			title: 'Straightshot',
			message: 'Pics or it didn\'t happen!',
			iconUrl: 'logo.png',
			buttons: [{ title: 'OK!' }, { title: 'GTFO' }]
		}, function(id) {});
	});
}

chrome.notifications.onButtonClicked.addListener(function(id, index) {
	if (id == 'straightshot') {
		chrome.notifications.clear(id, function() {});
		if (index == 0) {
			chrome.app.window.create('index.html', {
				id: 'straightshot',
				bounds: {
					width: 700,
					height: 600
				}
			});
		}
	}
});

chrome.alarms.create('straightshot', { periodInMinutes: 1 });

chrome.alarms.onAlarm.addListener(function (alarm) {
	if (alarm.name == 'straightshot') {
		showNotification();
	}
});

//chrome.app.runtime.onInstalled.addListener(function() {
//chrome.app.runtime.onLaunched.addListener(function() {
showNotification();
//});
