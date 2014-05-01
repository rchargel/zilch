angular.module("zilch", []).controller("RootController", function($http, $scope) {
	$scope.countries = [];
	$scope.total = 0;

	$scope.init = function() {
		$http.get("/countries.json").success(function(countries) {
			var list = [], t = 0;
			for (var i = 0; i < countries.length; i++) {
				list.push(countries[i]); 
				for (var j = 0; j < countries[i].States.length; j++) {
					t += countries[i].States[j].ZipCodes
				}
			}
			$scope.countries = list;
			$scope.total = t.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		});
	};

	$scope.init();
});
