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

app.config(['$routeProvider', '$locationProvider', function($routeProvider) {
    $routeProvider.when('/login', {templateUrl: '/partials/login.html', controller: 'LoginCtrl'});
    $routeProvider.when('/game', {templateUrl: '/partials/game.html', controller: 'GameCtrl'});
    $routeProvider.when('/createGame', {templateUrl: '/partials/createGame.html', controller: 'GameCtrl'});
    $routeProvider.when('/lobby', {templateUrl: '/partials/lobby.html', controller: 'MainCtrl'});
    $routeProvider.otherwise({redirectTo: '/lobby'});
}])
.run(function($rootScope, $location) {
    $rootScope.$on('$routeChangeSuccess', function () {
        if (!$rootScope.name) {
            $location.path("/login");
        }
    })
});

