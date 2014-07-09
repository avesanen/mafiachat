'use strict';

angular.module('mafiachat.filters', []).
    filter('nospace', function () {
        return function (value) {
            return (!value) ? '' : value.replace(/ /g, '');
        };
    }
);