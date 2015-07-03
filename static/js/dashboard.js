var dcron = angular.module('dcron', []);

dcron.controller('JobListCtrl', function ($scope, $http) {
  $scope.click = function(jobName) {
    var response = $http.put('/jobs/' + jobName);
    response.success(function(data, status, headers, config) {
      alert("Success running job " + jobName);
    });

    response.error(function(data, status, headers, config) {
      alert("Error running job " + jobName);
    });
  };
});
