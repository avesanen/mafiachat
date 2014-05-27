'use strict';

angular.module('mafiachat.services').factory('ResponseHandler', ['$q', '$rootScope', function($q, $rootScope) {
    // We return this object to anything injecting our service
    var Service = {};

    Service.handle = function($scope, msg) {
        if ($scope.log != "") {
            $scope.log += "<br />";
        }

        var now = new Date();
        $scope.log += ('0'  + now.getHours()).slice(-2)+':'+('0' + now.getMinutes()).slice(-2) + " ";
        switch (msg.msgType) {
            case 'joinGame':
                if (!$scope.players) $scope.players = [];
                $scope.players.push(msg.data.name);
                $scope.log += "<b>" + msg.data.name + " joined the game!</b>";
                break;
            case 'chatMessage':
                $scope.log += "<b>"+msg.data.player.name+":</b>&nbsp;<span class='"+msg.data.faction+"Message'>" + msg.data.message + "</span>";
                break;
            case 'login':
                return;
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

                $scope.gameName = msg.data.game.name;

                var oldPlayers = $scope.gameInfo.game.players;
                $scope.gameInfo = msg.data;
                var newPlayers = $scope.gameInfo.game.players;

                // check if someone joined or leaved the game and who that was (oldPlayers != gameInfo.game.players)
                if (newPlayers.length > oldPlayers.length) {
                    $scope.log += getJoinedOrPartedPlayer(newPlayers, oldPlayers, true) + " joined the game";
                } else {
                    $scope.log += getJoinedOrPartedPlayer(oldPlayers, newPlayers, false) + " left the game";
                }

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

