'use strict';

angular.module('mafiachat.controllers').controller('GameCtrl', ['$rootScope', '$scope', '$location', 'WebSocket', 'ResponseHandler', function($rootScope, $scope, $location, WebSocket, ResponseHandler) {
    $scope.log = "<b>Welcome " + $scope.name + "!</b>";

    $scope.createGame = function() {
        console.log("Create game: ", $scope.gameName);
        var game = {"id":$scope.games.length, "name":$scope.gameName};
        $rootScope.games.push(game);
        $location.path("/lobby");
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
}]);

