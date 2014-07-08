'use strict';

angular.module('mafiachat.filters').
    filter('max', function() {
        return function(input, min) {
            return Math.min(input, min);
        };
    }).
    filter('min', function() {
        return function(input, min) {
            return Math.max(input, min);
        };
    });

