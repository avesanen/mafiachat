'use strict';

angular.module('mafiachat.controllers').controller('GameCtrl', ['$rootScope', '$scope', '$location', 'WebSocket', 'ResponseHandler', function($rootScope, $scope, $location, WebSocket, ResponseHandler) {
    WebSocket.setScope($scope);

    //$scope.log = "<b>Welcome " + $scope.name + "!</b>";
    if (!$scope.gameInfo) {
        $location.path("/login");
    }
    
    if (!$scope.games) {
        $rootScope.games = [];
        $scope.games = [];
    }


    $scope.roleDescription = {
        mafia:"Mafia",
        villager:"Villager",
        doctor:"Doctor",
        cop:"Cop",
        dead:"Dead",
        new:"Unknown"
    };

    if (!$scope.gameInfo) $scope.gameInfo = {};
    if (!$scope.gameInfo.game) $scope.gameInfo.game = {name:"Game"};
    if (!$scope.gameInfo.game.players) {
        $scope.gameInfo.game.players = [
            {"name":$rootScope.name, "state":"mafia"},
            {"name":"Jakke", "state":"villager", votes:0},
            {"name":"Makke", "state":"doctor", votes:0},
            {"name":"Sakke", "state":"cop", votes:0},
            {"name":"Takke", "state":"dead", votes:0},
            {"name":"Nakke", "state":"new", votes:0},
            {"name":"Takke", "state":"dead", votes:0},
            {"name":"Takke", "state":"dead", votes:0},
            {"name":"Takke", "state":"dead", votes:0},
            {"name":"Takke", "state":"dead", votes:0},
            {"name":"Takke", "state":"dead", votes:0},
            {"name":"Nakke", "state":"new", votes:0},
            {"name":"Nakke", "state":"new", votes:0},
            {"name":"Nakke", "state":"new", votes:0},
            {"name":"Nakke", "state":"new", votes:0},
            {"name":"Nakke", "state":"new", votes:0},
            {"name":"Nakke", "state":"new", votes:0}
        ];
    }

    $scope.game = {
        "id":$scope.games.length,
        "name":"",
        "maxPlayers":8,
        "minPlayers":8,
        "minVillagers":2,
        "cops":1,
        "doctors":2,
        "mafiosi":3,
        "password":"",
        "state":"open",
        "players":[]
    };

    $scope.createGame = function() {
        $scope.game.needsPassword = $scope.game.password != "";
        $rootScope.games.push($scope.game);
        $location.path("/game");
    }

    $scope.startGame = function() {
        if ($scope.thisPlayer.admin) {
            var message = {data:{}};
            message.msgType = 'actionMessage';
            message.data.action = 'startGame';
            WebSocket.sendMsg(message);
        }
    }

    $scope.sendMsg = function() {
        if (!$scope.msg || !$scope.msgType) {
            //$scope.log += "<br /><span class='glyphicon glyphicon-exclamation-sign text-danger'></span><span class='text-danger'> Select chat room and enter a message</span>";
            return false;
        }

        sendChatMessage($scope.msgType, $scope.msg);
    }

    var sendChatMessage = function(type, msg) {
        var message = {data:{}};
        message.msgType = 'chatMessage';
        message.data.faction = type;
        message.data.message = msg;
        message.data.player = {};
        message.data.player.name = $rootScope.name;
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
    };

    $scope.kill = function(player) {
        sendActionMessage("kill", player.name);
    };

    $scope.identify = function(player) {
        sendActionMessage("identify", player.name);
    };

    $scope.heal = function(player) {
        sendActionMessage("heal", player.name);
    };
}]);

