// TODO: Remove once Safari supports fetch
function fetchJSONFile(path, callback) {
	var httpRequest = new XMLHttpRequest();
	httpRequest.onreadystatechange = function() {
		if (httpRequest.readyState === 4) {
			if (httpRequest.status === 200) {
				var data = JSON.parse(httpRequest.responseText);
				if (callback) callback(data);
			}
		}
	};
	httpRequest.open('GET', path);
	httpRequest.send();
}

fetchJSONFile('/nsg.php', function(json){
	fetchJSONFile('http://data.hazewatchapp.com/index_v2.json', function(json2){
			var ractive = new Ractive({
				el: '#container',
				template: '#template',
				data: { items: json, haze: json2 }
			});
			});
});


var helpers = Ractive.defaults.data;
helpers.fromNow = function(timeString){
    return moment(timeString).fromNow()
}
helpers.formatTime = function(timeString){
    return moment(timeString).format("ddd, h:mmA");
}
helpers.humanizeTime = function(timeString){
    return moment.duration(timeString).humanize();
}
