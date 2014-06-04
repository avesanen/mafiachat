'use strict';

var app = angular.module(
    'mafiachat',
    [
        'ngRoute',
        'ngSanitize',
        'ui.bootstrap',
        'mafiachat.services',
        'mafiachat.filters',
        'mafiachat.controllers'
    ]
);

function configConstants($rootScope) {
    $rootScope.minPlayers = 3;
}

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/login', {templateUrl: '/partials/login.html', controller: 'LoginCtrl'});
    $routeProvider.when('/game', {templateUrl: '/partials/lobby.html', controller: 'GameCtrl'});
    $routeProvider.when('/createGame', {templateUrl: '/partials/createGame.html', controller: 'GameCtrl'});
    $routeProvider.when('/games', {templateUrl: '/partials/games.html', controller: 'MainCtrl'});
    $routeProvider.otherwise({redirectTo: '/games'});
}])
.run(function($rootScope, $location, $window) {
    configConstants($rootScope);

    if ($location.absUrl().indexOf("/g/") > 0) {
        $location.path("/game");
    }

    $rootScope.home = function() {
        if ($location.path() != "/games") {
            $window.location.href = "/";
        }
    }

    $rootScope.$on('$routeChangeStart', function (event, next) {
        $rootScope.currentView = $location.path();
        // Check that login is done if any other page is requested than games listing
        if ($location.path() != "/games" && $location.path() != "/login" && !$rootScope.name) {
            // Check HTML5 storage
            if (sessionStorage.name) {
                console.log("Loading user from session storage: " + sessionStorage.name);
                $rootScope.name = sessionStorage.name;
            } else {
                $rootScope.requiredPath = $location.path();
                $location.path("/login");
            }
        }
    })
});

