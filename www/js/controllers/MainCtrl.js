'use strict';

angular.module('mafiachat.controllers', []).controller('MainCtrl', ['$rootScope', '$scope', '$location', 'WebSocket', 'ResponseHandler', function($rootScope, $scope, $location, WebSocket, ResponseHandler) {
    $scope.log = "<b>Welcome " + $scope.name + "!</b>";

    if (!$rootScope.games) {
        $rootScope.games = [{"id":1, "name":"Jea tässä olis yks peli."}];
    }

    WebSocket.setScope($scope);

    if ($rootScope.name) {
        var message = {data:{}};
        message.msgType = "gameInfo";
        WebSocket.sendMsg(message);
    }

    $scope.sendMsg = function() {
        if (!$scope.msg || !$scope.msgType) {
            return false;
        }

        var message = {data:{}};
        var type = $scope.msgType;

        if (type === 'villagerChat') {
            message.msgType = 'chatMessage';
            message.data.faction = 'villager';
            message.data.message = $scope.msg;
        }
        if (type === 'mafiaChat') {
            message.msgType = 'chatMessage';
            message.data.faction = 'mafia';
            message.data.message = $scope.msg;
        }
        if (type === 'joinGame') {
            message.msgType = 'joinGame';
            message.data.name = $scope.msg;
        }

        WebSocket.sendMsg(message);
        $scope.msg = "";
    }

    $scope.joinGame = function(gameId) {
        var message = {data:{}};
        message.msgType = 'joinGame';
        message.data.gameId = gameId;
        WebSocket.sendDeferMsg(message).
            then(function(resp) {
                $location.path("/game");
            }, function(error) {
                $scope.errorMsg = "Couldn't connect to backend :(";
            }
        );
    }

    $scope.newGame = function() {
        $location.path("/createGame");
    }
}]);

