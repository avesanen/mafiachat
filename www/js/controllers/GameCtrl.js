'use strict';

angular.module('mafiachat.controllers').controller('GameCtrl', ['$rootScope', '$scope', '$location', 'WebSocket', 'ResponseHandler', function($rootScope, $scope, $location, WebSocket, ResponseHandler) {
    WebSocket.setScope($scope);

    if (!$rootScope.name) {
        $location.path("/login");
    } else {
        var message = {data:{}};
        message.msgType = 'loginMessage';
        message.data.name = $rootScope.name;
        message.data.password = sessionStorage.pass;

        WebSocket.sendMsg(message);
    }

    $scope.roleDescription = {
        mafia:"Mafia",
        villager:"Villager",
        doctor:"Doctor",
        cop:"Cop",
        dead:"Dead",
        new:"Unknown"
    };

    $scope.startGame = function() {
        // Commented out to fasten testing ;) TODO: uncomment
        //if ($scope.thisPlayer.admin && $scope.gameInfo.game.players.length >= $rootScope.minPlayers) {
            var message = {data:{}};
            message.msgType = 'actionMessage';
            message.data.action = 'startGame';
            WebSocket.sendMsg(message);
        //}
    }

    $scope.sendMsg = function() {
        if (!$scope.msg || !$scope.msgType) {
            return false;
        }

        sendChatMessage($scope.msgType, $scope.msg);
    }

    var sendChatMessage = function(type, msg) {
        var message = {data:{}};
        message.msgType = 'chatMessage';
        message.data.faction = type;
        message.data.message = msg;
        message.data.player = $rootScope.name;
        WebSocket.sendMsg(message);
        $scope.msg = "";
    }

    var sendActionMessage = function(action, playerName) {
        var message = {data:{}};
        message.msgType = 'actionMessage';
        message.data.action = action;
        message.data.target = playerName;
        WebSocket.sendMsg(message);
    }

    $scope.vote = function(player) {
        sendActionMessage("vote", player.name);
        $scope.togglePlayerList();
    };

    $scope.kill = function(player) {
        sendActionMessage("kill", player.name);
        $scope.togglePlayerList();
    };

    $scope.identify = function(player) {
        sendActionMessage("identify", player.name);
        $scope.togglePlayerList();
    };

    $scope.heal = function(player) {
        sendActionMessage("heal", player.name);
        $scope.togglePlayerList();
    };

    $scope.contextMenuAvailable = function() {
        // Doctors can heal themselves
        if (this.player.name == $scope.game.myPlayer.name && $scope.game.myPlayer.faction != 'doctor') {
            return false;
        }

        // Ghosts can't do anything
        if ($scope.game.myPlayer.faction == 'ghost' || this.player.faction == 'ghost') {
            return false;
        }

        // If game is in lobby or debrief state
        if ($scope.game.state == 'lobby' || $scope.game.state == 'debrief') {
            return false;
        }

        // If myPlayer is cop and player is already identified
        if ($scope.game.state == 'night' && $scope.game.myPlayer.faction == 'cop' && this.player.faction != 'unknown') {
            return false;
        }

        return true;
    }

    $scope.togglePlayerList = function() {
        $scope.playersListToggle = $scope.playersListToggle == 'expanded' ? '' : 'expanded';
    }
}]);

