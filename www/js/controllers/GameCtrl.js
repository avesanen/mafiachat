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
        console.log("PASS:" , $scope.game.password, $scope.game.needsPassword);
        $rootScope.games.push($scope.game);
        $location.path("/game");
    }

    $scope.sendMsg = function() {
        if (!$scope.msg || !$scope.msgType) {
            //$scope.log += "<br /><span class='glyphicon glyphicon-exclamation-sign text-danger'></span><span class='text-danger'> Select chat room and enter a message</span>";
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

    $scope.vote = function(player) {
        var message = {data:{}};
        message.msgType = 'actionMessage';
        message.data.action = 'vote';
        message.data.target = player.name;
        WebSocket.sendMsg(message);
        /*
        // TODO: Move this logic in backend
        $scope.log += "<br><b>*** " + $rootScope.name + " voted for player " + player.name + "!</b>";
        player.votes++;
        var highestVotedPlayer = player;
        for (var i = 0; i < $scope.gameInfo.game.players.length; i++) {
            var p = $scope.gameInfo.game.players[i];
            if (p.votes > 0) {
                p.voteLevel = 'warning';
                if (p.votes > highestVotedPlayer.votes) {
                    highestVotedPlayer = p;
                }
            }
        }

        highestVotedPlayer.voteLevel = 'danger';
        // Highlight all players with same vote count
        for (var i = 0; i < $scope.gameInfo.game.players.length; i++) {
            var p = $scope.gameInfo.game.players[i];
            if (p.votes == highestVotedPlayer.votes) {
                p.voteLevel = 'danger';
            }
        }
        */
    };

    $scope.kill = function(player) {
        //$scope.log += "<br><b>*** " + $rootScope.name + " kills player " + player.name + "!</b>";
    };

    $scope.identify = function(player) {
        var isMafioso = true; // TODO: from backend
        if (isMafioso) {
            //$scope.log += "<br><b>*** " + player.name + " is a mafioso!</b>";
        } else {
            //$scope.log += "<br><b>*** " + player.name + " ain't a mafioso.</b>";
        }
    };

    $scope.heal = function(player) {
        //$scope.log += "<br><b>*** " + $rootScope.name + " heals player " + player.name + "!</b>";
    };
}]);

