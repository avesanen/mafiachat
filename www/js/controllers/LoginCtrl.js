'use strict';

angular.module('mafiachat.controllers').controller('LoginCtrl', ['$rootScope', '$scope', '$location', '$routeParams', 'WebSocket', 'GameService', function($rootScope, $scope, $location, $routeParams, WebSocket, GameService) {
    WebSocket.setScope($rootScope);

    $scope.login = function() {

        if (!$scope.name || !$scope.password) {
            $scope.errorMsg = "Enter all data.";
            return;
        }

        var message = {data:{}};
        message.msgType = 'login';
        message.data.name = $scope.name;
        message.data.password = $scope.password;

        var gameData = GameService.getGameData();
        if (gameData.id) {
            message.data.gameId = gameData.id;
        }

        WebSocket.sendDeferMsg(message).
            then(function(resp) {
                sessionStorage.name = $scope.name;
                $rootScope.name = $scope.name;
                if ($rootScope.requiredPath) {
                    $location.path($rootScope.requiredPath);
                } else {
                    $location.$$search = {};
                    $location.path("/game");
                }
            }, function(error) {
                $scope.errorMsg = "Couldn't connect to backend :(";
            }
        );
    }

}]);

