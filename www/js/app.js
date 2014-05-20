'use strict';

var app = angular.module(
    'mafiachat',
    [
        'ngRoute',
        'ngSanitize',
        //'mafiachat.filters',
         'mafiachat.services',
         'mafiachat.controllers'
    ]
);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/game', {templateUrl: '/partials/game.html', controller: 'GameCtrl'});
    $routeProvider.when('/lobby', {templateUrl: '/partials/lobby.html', controller: 'MainCtrl'});
    $routeProvider.otherwise({redirectTo: '/lobby'});
}]);

