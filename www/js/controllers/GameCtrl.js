'use strict';

angular.module('mafiachat.controllers').controller('GameCtrl', ['$rootScope', '$scope', '$location', 'WebSocket', 'ResponseHandler', function($rootScope, $scope, $location, WebSocket, ResponseHandler) {
    WebSocket.setScope($scope);

    $scope.log = "<b>Welcome " + $scope.name + "!</b>";

    $scope.createGame = function() {
        var needsPassword = $scope.gamePassword != '';
        var game = {
            "id":$scope.games.length,
            "name":$scope.gameName,
            "needsPassword":needsPassword,
            "maxPlayers":$scope.maxPlayers,
            "cops":$scope.cops,
            "doctors":$scope.doctors,
            "mafiosi":$scope.mafiosi,
            "players":[]
        };
        $rootScope.games.push(game);
        $location.path("/game");
    }

    $scope.sendMsg = function() {
        if (!$scope.msg || !$scope.msgType) {
            $scope.log += "<br /><span class='glyphicon glyphicon-exclamation-sign text-danger'></span><span class='text-danger'> Select chat room and enter a message</span>";
            return false;
        }

        var message = {data:{}};
        var type = $scope.msgType;

        if (type === 'villagerChat') {
            message.msgType = 'chatMessage';
            message.data.faction = 'villager';
            message.data.message = $scope.msg;
            message.data.player = {};
            message.data.player.name = $rootScope.name;
        }
        if (type === 'mafiaChat') {
            message.msgType = 'chatMessage';
            message.data.faction = 'mafia';
            message.data.message = $scope.msg;
            message.data.player = {};
            message.data.player.name = $rootScope.name;
        }

        WebSocket.sendMsg(message);
        $scope.msg = "";
    }
}]);

