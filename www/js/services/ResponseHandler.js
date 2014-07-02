'use strict';

angular.module('mafiachat.services').factory('ResponseHandler', ['$q', '$rootScope', function($q, $rootScope) {
    // We return this object to anything injecting our service
    var Service = {};

    Service.handle = function($scope, msg) {
        if (!$scope.messageBuffer) {
            $scope.messageBuffer = [];
        }

        switch (msg.msgType) {
            case 'chatMessage':
                $scope.messageBuffer.push(msg.data);
                break;
            case 'loginFailed':
                break;
            default:
                $scope.game = msg;
                $scope.messageBuffer = msg.messages;
                calculateVoteLevels($scope.game);
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

        if ($scope) {
            $scope.$apply();
        }

        var objDiv = document.getElementById("log");
        if (objDiv) {
            objDiv.scrollTop = objDiv.scrollHeight;
        }
    };

    function calculateVoteLevels(game) {
        var highestVotedPlayer;
        for (var i = 0; i < game.players.length; i++) {
            var p = game.players[i];
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
            for (var i = 0; i < game.players.length; i++) {
                var p = game.players[i];
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

