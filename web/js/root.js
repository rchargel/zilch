angular.module("zilch", []).controller("RootController", function($http, $scope) {
	$scope.countries = [];
	$scope.total = 0;

	$scope.init = function() {
		$http.get("/countries.json").success(function(map) {
			var list = [], t = 0;
			for (var country in map) {
				t += map[country];
				if (map[country] >= 100) {
					list.push(country);
				} 
			}
			list.sort();
			$scope.countries = list;
			$scope.total = t.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
		});
	};

	$scope.init();
});
