angular.module('zilch', []).controller('ZilchController', function($http, $scope) {
	$scope.url = '/query.js?callback=JSON_CALLBACK';
	$scope.entries = [];
	$scope.entry = null;
	$scope.errorMessage = null;

	$scope.$watch('zipCode', function(value) {
		$scope.findEntries();
	});

	$scope.clear = function() {
		$scope.entries = [];
		$scope.entry = null;
		$scope.errorMessage = null;
	};

	$scope.assignCity = function(entry) {
		$scope.clear();
		$scope.entry = entry;
	};
	
	$scope.findEntries = function() {
		if ($scope.zipCode && $scope.zipCode.length >= 3) {
			var url = $scope.url + '&ZipCode=' + $scope.zipCode;
			$http.jsonp(url).success(function(response) {
				$scope.clear();
				var selectEntries = $scope.toSelectEntries(response.ZipCodeEntries);
				if (selectEntries.length == 1) {
					$scope.entry = selectEntries[0];
				} else if (selectEntries.length > 0 && selectEntries.length < 15) {
					$scope.entries = selectEntries;
				} else if (selectEntries.length >= 10) {
					$scope.errorMessage = "There are too many responses for this query: " + selectEntries.length;
				}
			});
		}
	};

	$scope.toSelectEntries = function(zipEntries) {
		var entries = [], i, j;
		for (i = 0; i < zipEntries.length; i++) {
			entries.push({
				city: zipEntries[i].City,
				state: zipEntries[i].State,
				country: zipEntries[i].Country,
				latitude: zipEntries[i].Latitude,
				longitude: zipEntries[i].Longitude
			});
			for (j = 0; j < zipEntries[i].AcceptableCities.length; j++) {
				entries.push({
					city: zipEntries[i].AcceptableCities[j],
					state: zipEntries[i].State,
					country: zipEntries[i].Country,
					latitude: zipEntries[i].Latitude,
					longitude: zipEntries[i].Longitude
				});
			}
		}
		return entries;
	};
});
