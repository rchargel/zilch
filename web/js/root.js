angular.module("zilch", []).controller("RootController", function($http, $scope) {
	$scope.countries = [];
	$scope.total = 0;

	$scope.init = function() {
		$http.get("/countries.json").success(function(countries) {
			var list = [], t = 0, countryTotal = 0;
			for (var i = 0; i < countries.length; i++) {
				countryTotal = 0;
				for (var j = 0; j < countries[i].States.length; j++) {
					t += countries[i].States[j].ZipCodes
					countryTotal += countries[i].States[j].ZipCodes
				}
				countries[i].Total = countryTotal.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
				list.push(countries[i]); 
			}
			$scope.countries = list;
			$scope.total = t.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		});
	};

	$scope.init();
});
