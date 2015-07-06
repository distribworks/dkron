var dcron = angular.module('dcron', ['angular-rickshaw']);

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

dcron.controller('IndexCtrl', function ($scope, $http, $interval) {
  $scope.options = {
    renderer: 'line'
  };

  $scope.series = [{
      name: 'Series 1',
      color: 'steelblue',
      data: [{x: 0, y: 0}]
  }];
  $scope.features = {
      hover: {
          xFormatter: function(x) {
              return 't=' + x;
          },
          yFormatter: function(y) {
              return '$' + y;
          }
      },
      legend: {
        toggle: false,
        highlight: true
      }
  };
  $interval(function() {
    data = $scope.series[0].data;
    data.push({x: data[data.length - 1].x + 10, y: 60 * Math.random()});

    $scope.series = [{
      name: 'Series 1',
      color: 'steelblue',
      data: data
    }];
  }, 1000);
});
