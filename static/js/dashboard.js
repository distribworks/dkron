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

dcron.controller('IndexCtrl', function ($scope, $http, $interval, $element) {
  $scope.options = {
    renderer: 'line',
    interpolation: 'linear'
  };

  $scope.series = [{
      name: 'Success count',
      color: 'darkgreen',
      data: [{x: 0, y: 0}]
  },{
      name: 'Error count',
      color: 'red',
      data: [{x: 0, y: 0}]
  }];
  $scope.features = {
      hover: {
          xFormatter: function(x) {
              return x;
          },
          yFormatter: function(y) {
              return y;
          }
      },
      legend: {
        toggle: false,
        highlight: true
      },
      yAxis: {
        tickFormat: 'formatKMBT'
      },
      xAxis: {
        tickFormat: 'formatKMBT',
        timeUnit: 'hour'
      }
  };

  $interval(function() {
    var response = $http.get('/jobs/');
    response.success(function(data, status, headers, config) {
      $scope.updateGraph(data);
    });

    response.error(function(data, status, headers, config) {
      alert("Error getting data");
    });
  }, 2000);

  $scope.success_count = 0;
  $scope.error_count = 0;

  $scope.updateGraph = function(data) {
    var gdata = $scope.series[0].data;
    var name = $scope.series[0].name;
    var color = $scope.series[0].color;
    var success_count = 0;
    var diff = 0;

    for(i=0; i < data.length; i++) {
      success_count = success_count + data[i].success_count;
    }
    if($scope.success_count != 0) {
      diff = success_count - $scope.success_count;
    }
    $scope.success_count = success_count;

    gdata.push({x: gdata[gdata.length - 1].x + 1, y: diff})

    $scope.series[0] = {
      name: name,
      color: color,
      data: gdata
    };

    gdata = $scope.series[1].data;
    name = $scope.series[1].name;
    color = $scope.series[1].color;
    var error_count = 0;

    for(i=0; i < data.length; i++) {
      error_count = error_count + data[i].error_count;
    }
    if($scope.error_count != 0) {
      diff = error_count - $scope.error_count;
    }
    $scope.error_count = error_count;

    gdata.push({x: gdata[gdata.length - 1].x + 1, y: diff})

    $scope.series[1] = {
      name: name,
      color: color,
      data: gdata
    };
  }
});
