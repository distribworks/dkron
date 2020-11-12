var dkron = angular.module('dkron', ['angular-rickshaw', 'ui.codemirror']);

dkron.filter('statusClass', function () {
  var friendlyStatusFilter = function (job) {
    if (job.disabled) {
      return 'text-muted glyphicon-ban-circle'
    }
    switch (job.status) {
      case 'success':
        return 'status-success glyphicon-ok-sign'
      case 'failed':
        return 'status-danger glyphicon-remove-sign'
      case 'partially_failed':
        return 'status-warning glyphicon-exclamation-sign'
      case 'running':
        return 'status-running glyphicon-play-circle'
      default:
        return 'glyphicon-question-sign'
    }
    return input;
  };
  return friendlyStatusFilter;
});

dkron.constant('hideDelay', 2000);

dkron.filter('startFrom', function () {
  return function (input, start) {
      if (input) {
          start = +start; //parse to int
          return input.slice(start);
      }
      return [];
  }
});

dkron.run(function($rootScope) {
  $rootScope.alertOptions = function (type) {
    return {
      type: type,
      placement: {
        from: "top",
        align: "center"
      },
      offset: 40,
      animate: {
        enter: 'animated fadeInDown',
        exit: 'animated fadeOutUp'
      }
    }
  };
});

dkron.controller('JobListCtrl', function ($scope, $location, $http, $interval, hideDelay, $rootScope) {
  $scope.searchJob = $location.search()['filter']

  $scope.filter = function (filter) {
    $scope.searchJob = filter;
  };

  // pretty json func
  $scope.toJson = function (str) {
    return angular.toJson(str, true);
  }
  $scope.jobTemplate = {
    name: "job",
    displayname: "",
    schedule: "",
    owner: "",
    owner_email: "",
    disabled: false,
    tags: {},
    retries: 0,
    processors: null,
    concurrency: "allow",
    executor: "shell",
    executor_config: {
      command: "/bin/true"
    }
  }
  $scope.jobTemplateJson = $scope.toJson($scope.jobTemplate);

  $scope.editorOptions = {
    lineWrapping: true,
    autoCloseBrackets: true,
    autoRefresh: true,
    height: "auto",
    mode: "application/json",
  };

  $scope.runJob = function (jobName) {
    let i = $scope.pagedItems[$scope.currentPage].findIndex(j => j.name === jobName);
    $scope["running_" + i] = true;
    $http.post(DKRON_API_PATH + '/jobs/' + jobName).
      then(function (response) {
        jQuery.notify({
          message: 'Success running job ' + jobName
        },$rootScope.alertOptions('success'));
        updateView();
        $scope["running_" + i] = false;
    }, function (response) {
      jQuery.notify({
        message: 'Error running job ' + jobName
      },$rootScope.alertOptions('danger'));
    });
  };

  $scope.createJob = function (jobTemplateJson) {
    try {
      job = angular.fromJson(jobTemplateJson);
    } catch (err) {
      window.alert('Json Format Error');
      return
    }
    
    $http.post(DKRON_API_PATH + '/jobs', job).
      then(function (response) {
        jQuery.notify({
          message: 'Success created job ' + job.name
        },$rootScope.alertOptions('success'));
        updateView();
    }, function (response) {
      jQuery.notify({
        message: 'Error creating job ' + job.name
      },$rootScope.alertOptions('danger'));
    });
  };

  $scope.updateJob = function (jobJson) {
    try {
      job = angular.fromJson(jobJson);
    } catch (err) {
      window.alert('Json Format Error');
      return
    }
    
    $http.post(DKRON_API_PATH + '/jobs', job).
      then(function (response) {
        jQuery.notify({
          message: 'Success updating job ' + jobName
        },$rootScope.alertOptions('success'));
        updateView();
    }, function (response) {
      jQuery.notify({
        message: 'Error updating job ' + job.name
      },$rootScope.alertOptions('danger'));
    });
  };

  $scope.deleteJob = function (jobName) {
    let i = $scope.pagedItems[$scope.currentPage].findIndex(j => j.name === jobName);
    $scope["deleting_" + i] = true;
    
    $http.delete(DKRON_API_PATH + '/jobs/' + jobName).
      then(function (response) {
        jQuery.notify({
          message: 'Successfully removed job ' + jobName
        },$rootScope.alertOptions('success'));
        updateView();
    }, function (response) {
      jQuery.notify({
        message: 'Error removing job ' + jobName
      },$rootScope.alertOptions('danger'));
    });
  };

  $scope.toggleJob = function (jobName) {
    $http.post(DKRON_API_PATH + '/jobs/' + jobName + '/toggle').
      then(function (response) {
        jQuery.notify({
          message: 'Successfully toggled job ' + jobName
        },$rootScope.alertOptions('success'));
        updateView();
    }, function (response) {
      jQuery.notify({
        message: 'Error toggle job ' + jobName
      },$rootScope.alertOptions('danger'));
    });
  }

  var updateView = function () {
    $http.get(DKRON_API_PATH + '/jobs').
      then(function (response) {
        $scope.updateStatus(response.data);

        $("#conn-error").delay(hideDelay).slideUp(200, function () {
          $("#conn-error").alert('close');
        });
    }, function (response) {
      jQuery.notify({
        message: 'Error loading data'
      },$rootScope.alertOptions('success'));
    });
  }

  // calculate page in place
  $scope.groupToPages = function () {
    $scope.gap = Math.round($scope.jobs.length / $scope.itemsPerPage);
    $scope.gap = Math.min($scope.gap, 10);
    
    $scope.pagedItems = [];

    for (var i = 0; i < $scope.jobs.length; i++) {
      if (i % $scope.itemsPerPage === 0) {
        $scope.pagedItems[Math.floor(i / $scope.itemsPerPage)] = [$scope.jobs[i]];
      } else {
        $scope.pagedItems[Math.floor(i / $scope.itemsPerPage)].push($scope.jobs[i]);
      }
    }
  };

  $scope.range = function (size, start, end) {
    var ret = [];

    if (size < end) {
      end = size;
      start = size - $scope.gap;
    }
    for (var i = start; i < end; i++) {
      ret.push(i);
    }
    return ret;
  };

  $scope.firstPage = function () {
    if ($scope.currentPage > 0) {
      $scope.currentPage = 0;
    }
  };

  $scope.prevPage = function () {
    if ($scope.currentPage > 0) {
      $scope.currentPage--;
    }
  };

  $scope.nextPage = function () {
    if ($scope.currentPage < $scope.pagedItems.length - 1) {
      $scope.currentPage++;
    }
  };

  $scope.lastPage = function () {
    if ($scope.currentPage < $scope.pagedItems.length - 1) {
      $scope.currentPage = Math.ceil($scope.jobs.length / $scope.itemsPerPage) - 1;
    }
  };

  $scope.setPage = function () {
    $scope.currentPage = this.n;
  };

  // init
  $scope.gap = 0;
  $scope.groupedItems = [];
  $scope.itemsPerPage = 15;
  $scope.pagedItems = [];
  $scope.currentPage = 0;

  $scope.success_count = 0;
  $scope.error_count = 0;
  $scope.failed_jobs = 0;
  $scope.successful_jobs = 0;
  $scope.jobs = [];

  $scope.updateStatus = function (data) {
    var success_count = 0;
    var error_count = 0;
    $scope.jobs = data;

    $scope.failed_jobs = 0;
    $scope.successful_jobs = 0;

    // functions have been describe process the data for display
    $scope.groupToPages();

    for (i = 0; i < data.length; i++) {
      success_count = success_count + data[i].success_count;
      error_count = error_count + data[i].error_count;

      // Compute last...Dates: they're either a date or null
      var lastSuccessDate = data[i].last_success && new Date(Date.parse(data[i].last_success));
      var lastErrorDate = data[i].last_error && new Date(Date.parse(data[i].last_error));
      if ((lastSuccessDate !== null && lastErrorDate === null) || lastSuccessDate > lastErrorDate) {
        $scope.successful_jobs = $scope.successful_jobs + 1;
      } else if ((lastSuccessDate === null && lastErrorDate !== null) || lastSuccessDate < lastErrorDate) {
        $scope.failed_jobs = $scope.failed_jobs + 1;
      }
    }

    $scope.success_count = success_count;
    $scope.error_count = error_count;
  }

  updateView();
});

dkron.controller('ExecutionsCtrl', function ($scope, $http, $interval, hideDelay, $rootScope) {
  $scope.runJob = function (jobName) {
    $scope["running_job"] = true;
    $http.post(DKRON_API_PATH + '/jobs/' + jobName).
      then(function (response) {
        jQuery.notify({
          message: 'Success running job ' + jobName
        },$rootScope.alertOptions('success'));
        $scope["running_job"] = false;
    }, function (response) {
      jQuery.notify({
        message: 'Error running job ' + jobName
      },$rootScope.alertOptions('danger'));
    });
  };
});

dkron.controller('IndexCtrl', function ($scope, $http, $timeout, $element, $rootScope) {
  $scope.options = {
    renderer: 'line',
    interpolation: 'linear'
  };

  $scope.features = {
    hover: {
      xFormatter: function (x) {
        return new Date(x*1000);
      },
      yFormatter: function (y) {
        return y;
      }
    },
    legend: {
      toggle: false,
      highlight: true
    },
    yAxis: {
      tickFormat: Rickshaw.Fixtures.Number.formatKMBT
    }
  };

  $scope.series = new Rickshaw.Series.FixedDuration([{name: 'Success', color: '#006f68'}, {name: 'Failed', color: '#B1003E'}], undefined, {
    timeInterval: 2000,
    maxDataPoints: 100
  });

  updateView = function () {
    $http.get(DKRON_API_PATH + '/jobs').
      then(function (response) {
        $scope.updateGraph(response.data);

        $timeout(function () {
          updateView();
        }, 2000);
    }, function (response) {
      jQuery.notify({
        message: 'Error loading data'
      },$rootScope.alertOptions('danger'));
    });

    $http.get(DKRON_API_PATH + '/members').
      then(function (response) {
        var data = response.data;
        angular.forEach(data, function (val, key) {
          switch (val.Status) {
            case 0:
              data[key].Status = "none";
              break;
            case 1:
              data[key].Status = "alive";
              break;
            case 2:
              data[key].Status = "leaving";
              break;
            case 3:
              data[key].Status = "left";
              break;
            case 4:
              data[key].Status = "failed";
              break;
          }
        });
        $scope.members = data;
    }, function (response) {
      jQuery.notify({
        message: 'Error loading data'
      },$rootScope.alertOptions('danger'));
    });
  }

  // Init values
  $scope.success_count = 0;
  $scope.error_count = 0;
  $scope.failed_jobs = 0;
  $scope.successful_jobs = 0;
  $scope.running = 0;
  $scope.jobs = [];
  $scope.i = 0;

  $scope.updateGraph = function (data) {
    var success_count = 0;
    var error_count = 0;
    var running = 0;

    $scope.jobs = data;
    $scope.failed_jobs = 0;
    $scope.successful_jobs = 0;

    for (i = 0; i < data.length; i++) {
      success_count = success_count + data[i].success_count;
      error_count = error_count + data[i].error_count;

      // Compute last...Dates: they're either a date or null
      var lastSuccessDate = data[i].last_success && new Date(Date.parse(data[i].last_success));
      var lastErrorDate = data[i].last_error && new Date(Date.parse(data[i].last_error));
      if ((lastSuccessDate !== null && lastErrorDate === null) || lastSuccessDate > lastErrorDate) {
        $scope.successful_jobs = $scope.successful_jobs + 1;
      } else if ((lastSuccessDate === null && lastErrorDate !== null) || lastSuccessDate < lastErrorDate) {
        $scope.failed_jobs = $scope.failed_jobs + 1;
      }

      if (data[i].status == 'running') {
        running = running + 1;
      }
    }

    // Store the previous data
    if ($scope.i === 0) {
      var successData = success_count;
      var failedData = error_count;
    } else {
      var successData = $scope.success_count;
      var failedData = $scope.error_count;
    }

    // Update panel stats
    $scope.success_count = success_count;
    $scope.error_count = error_count;

    // Rickshaw graph update
    dataPoint = {}
    dataPoint['Success'] = success_count - successData;
    dataPoint['Failed'] = error_count - failedData;

    $scope.series.addData(dataPoint);

    // Broadcast a fake resize event to force render
    $scope.$broadcast('rickshaw::resize');
    $scope.i++;
  }

  updateView();
});

dkron.controller('BusyCtrl', function ($scope, $http, $timeout, $element, $rootScope) {
  busy = function () {
    $http.get(DKRON_API_PATH + '/busy').then(function (response) {
      $scope.executions = response.data

      $timeout(function () {
        busy();
      }, 500);
    }, function (response) {
      jQuery.notify({
        message: 'Error loading data'
      },$rootScope.alertOptions('danger'));
    });
  };
  busy();
});
