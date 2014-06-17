'use strict';

angular.module('mafiachat.services').factory('ResponseHandler', ['$q', '$rootScope', function($q, $rootScope) {
    // We return this object to anything injecting our service
    var Service = {};

    Service.handle = function($scope, msg) {
        switch (msg.msgType) {
            case 'messageBuffer':
                $scope.messageBuffer = msg.data;
                break;
            case 'chatMessage':
                $scope.messageBuffer.push(msg.data);
                break;
            case 'gameInfo':
                if (!$scope.gameInfo) {
                    $scope.gameInfo = {};
                }
                if (!$scope.gameInfo.game) {
                    $scope.gameInfo.game = {};
                }
                if (!$scope.gameInfo.game.players) {
                    $scope.gameInfo.game.players = [];
                }
                if (!$scope.messageBuffer) {
                    $scope.messageBuffer = [];
                }

                $scope.gameName = msg.data.game.name;

                //var oldPlayers = $scope.gameInfo.game.players;
                $scope.gameInfo = msg.data;
                calculateVoteLevels($scope.gameInfo);
                for (var i = 0; i < $scope.gameInfo.game.players.length; i++) {
                    var p = $scope.gameInfo.game.players[i];
                    if (p.name == $rootScope.name) {
                        $scope.thisPlayer = p;
                        break;
                    }
                }
                //var newPlayers = $scope.gameInfo.game.players;

                // check if someone joined or leaved the game and who that was (oldPlayers != gameInfo.game.players)
                /*
                if (newPlayers.length > oldPlayers.length) {
                    $scope.log += getJoinedOrPartedPlayer(newPlayers, oldPlayers, true) + " joined the game";
                } else {
                    $scope.log += getJoinedOrPartedPlayer(oldPlayers, newPlayers, false) + " left the game";
                }
                */

                break;
        }

        if ($scope) {
            $scope.$apply();
        }

        var objDiv = document.getElementById("log");
        if (objDiv) {
            objDiv.scrollTop = objDiv.scrollHeight;
        }
    };

    function calculateVoteLevels(gameInfo) {
        var highestVotedPlayer;
        for (var i = 0; i < gameInfo.game.players.length; i++) {
            var p = gameInfo.game.players[i];
            if (p.votes > 0) {
                p.voteLevel = 'warning';
                if (!highestVotedPlayer || p.votes > highestVotedPlayer.votes) {
                    highestVotedPlayer = p;
                }
            }
        }

        if (highestVotedPlayer) {
            highestVotedPlayer.voteLevel = 'danger';
            // Highlight all players with same vote count
            for (var i = 0; i < gameInfo.game.players.length; i++) {
                var p = gameInfo.game.players[i];
                if (p.votes == highestVotedPlayer.votes) {
                    p.voteLevel = 'danger';
                }
            }
        }
    }

    function getJoinedOrPartedPlayer(players1, players2, joins) {
        for (var i = 0; i < players1.length; i++) {
            var found = false;
            for (var j = 0; j < players2.length; j++) {
                if (players1[i].name == players2[j].name) {
                    found = true;
                }
            }
            if (joins && !found && players1[i].name != "") {
                return players1[i].name;

            } else if (!joins && !found && players1[i].name != "" ) {
                return players1[i].name;
            }
        }
        return "WTF!?";
    }

    return Service;
}]);

