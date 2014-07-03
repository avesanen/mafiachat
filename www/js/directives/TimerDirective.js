'use strict';

angular.module('mafiachat.directives').directive('timer', ['$rootScope', '$timeout', function($rootScope, $timeout) {
    return {
        restrict: 'E',
        scope: {
            time: '@',
            countdown: '@'
        },
        templateUrl: '/partials/timer.html',
        link: function(scope, element, attrs) {
            scope.onTimeout = function() {
                scope.timeSet = true;
                var seconds = parseInt(scope.time);
                if (scope.time <= 0) {
                    scope.seconds = 0;
                } else {
                    scope.minutes = Math.floor(seconds / 60);
                    scope.hours = Math.floor(scope.minutes / 60);
                    scope.minutes = scope.minutes % 60;
                    scope.seconds = seconds % 60;
                    scope.time = seconds + ((scope.countdown) ? -1 : 1);
                }
                mytimeout = $timeout(scope.onTimeout,1000);
            }
            var mytimeout = $timeout(scope.onTimeout,1000);
        }
    };
}]);

