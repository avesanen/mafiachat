'use strict';

angular.module('mafiachat.services').factory('ResponseHandler', ['$q', '$rootScope', function($q, $rootScope) {
    // We return this object to anything injecting our service
    var Service = {};

    Service.handle = function($scope, msg) {
        switch (msg.msgType) {
            case 'loginFailed':
                break;
            default:
                $scope.game = msg;
                $scope.messageBuffer = msg.messages;
                calculateVoteLevels($scope.game);

                if ($scope.game.state == 'villager-win' || $scope.game.state == 'mafia-win') {
                    $scope.factionHidden = undefined;
                } else {
                    $scope.factionHidden = "hidden";
                }
        }

        if ($scope) {
            $scope.$apply();

            if ($scope.contextMenuForPlayer) {
                angular.element(document.getElementById("player-"+$scope.contextMenuForPlayer)).triggerHandler('click');
            }
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

    return Service;
}]);

