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

dkron.controller('JobListCtrl', function ($scope, $location, $http, $interval, hideDelay) {
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
    var response = $http.post(DKRON_API_PATH + '/jobs/' + jobName);
    response.success(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-success fade in">Success running job ' + jobName + '</div>');
      updateView();
      $scope["running_" + i] = false;

      $(".alert-success").delay(hideDelay).slideUp(200, function () {
        $(".alert").alert('close');
      });
    });

    response.error(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-danger fade in"><button type="button" class="close close-alert" data-dismiss="alert" aria-hidden="true">x</button>Error running job ' + jobName + '</div>');
    });
  };

  $scope.createJob = function (jobTemplateJson) {
    try {
      job = angular.fromJson(jobTemplateJson);
    } catch (err) {
      window.alert('Json Format Error');
      return
    }
    var response = $http.post(DKRON_API_PATH + '/jobs', job);
    response.success(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-success fade in">Success created job ' + job.name + '</div>');
      updateView();

      $(".alert-success").delay(hideDelay).slideUp(200, function () {
        $(".alert").alert('close');
      });
    });

    response.error(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-danger fade in"><button type="button" class="close close-alert" data-dismiss="alert" aria-hidden="true">x</button>Error creating job ' + job.name + '</div>');
    });
  };

  $scope.updateJob = function (jobJson) {
    try {
      job = angular.fromJson(jobJson);
    } catch (err) {
      window.alert('Json Format Error');
      return
    }
    var response = $http.post(DKRON_API_PATH + '/jobs', job);
    response.success(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-success fade in">Success updating job ' + job.name + '</div>');
      updateView();

      $(".alert-success").delay(hideDelay).slideUp(200, function () {
        $(".alert").alert('close');
      });
    });

    response.error(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-danger fade in"><button type="button" class="close close-alert" data-dismiss="alert" aria-hidden="true">x</button>Error updating job ' + job.name + '</div>');
    });
  };

  $scope.deleteJob = function (jobName) {
    let i = $scope.pagedItems[$scope.currentPage].findIndex(j => j.name === jobName);
    $scope["deleting_" + i] = true;
    var response = $http.delete(DKRON_API_PATH + '/jobs/' + jobName);
    response.success(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-success fade in">Successfully removed job ' + jobName + '</div>');
      updateView();

      $(".alert-success").delay(hideDelay).slideUp(200, function () {
        $(".alert").alert('close');
      });
    });

    response.error(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-danger fade in"><button type="button" class="close close-alert" data-dismiss="alert" aria-hidden="true">x</button>Error removing job ' + jobName + '</div>');
    });
  };

  $scope.toggleJob = function (jobName) {
    var response = $http.post(DKRON_API_PATH + '/jobs/' + jobName + '/toggle');
    response.success(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-success fade in">Successfully toggled job ' + jobName + '</div>');
      updateView();

      $(".alert-success").delay(hideDelay).slideUp(200, function () {
        $(".alert").alert('close');
      });
    });

    response.error(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-danger fade in"><button type="button" class="close close-alert" data-dismiss="alert" aria-hidden="true">x</button>Error toggle job ' + jobName + '</div>');
    });
  }

  var updateView = function () {
    var response = $http.get(DKRON_API_PATH + '/jobs');
    response.success(function (data, status, headers, config) {
      $scope.updateStatus(data);

      $("#conn-error").delay(hideDelay).slideUp(200, function () {
        $("#conn-error").alert('close');
      });
    });

    response.error(function (data, status, headers, config) {
      $('#message').html('<div id="conn-error" class="alert alert-danger fade in"><button type="button" class="close close-alert" data-dismiss="alert" aria-hidden="true">x</button>Error getting data</div>');
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

dkron.controller('ExecutionsCtrl', function ($scope, $http, $interval, hideDelay) {
  $scope.runJob = function (jobName) {
    $scope["running_job"] = true;
    var response = $http.post(DKRON_API_PATH + '/jobs/' + jobName);
    response.success(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-success fade in">Success running job ' + jobName + '</div>');
      $scope["running_job"] = false;

      $(".alert-success").delay(hideDelay).slideUp(200, function () {
        $(".alert").alert('close');
        window.location.reload();
      });
    });

    response.error(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-danger fade in"><button type="button" class="close close-alert" data-dismiss="alert" aria-hidden="true">x</button>Error running job ' + jobName + '</div>');
    });
  };
});

dkron.controller('IndexCtrl', function ($scope, $http, $timeout, $element) {
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
    var response = $http.get(DKRON_API_PATH + '/jobs');
    response.success(function (data, status, headers, config) {
      $scope.updateGraph(data);

      $timeout(function () {
        updateView();
      }, 2000);
    });

    response.error(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-danger fade in"><button type="button" class="close close-alert" data-dismiss="alert" aria-hidden="true">x</button>Error getting data</div>');
    });

    var mq = $http.get(DKRON_API_PATH + '/members');
    mq.success(function (data, status, headers, config) {
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
    });

    mq.error(function (data, status, headers, config) {
      $('#message').html('<div class="alert alert-danger fade in"><button type="button" class="close close-alert" data-dismiss="alert" aria-hidden="true">x</button>Error getting data</div>');
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
