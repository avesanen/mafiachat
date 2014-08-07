'use strict';

angular.module('mafiachat.services').factory('ResponseHandler', ['$q', '$rootScope', function($q, $rootScope) {
    // We return this object to anything injecting our service
    var Service = {};

    Service.handle = function($scope, msg) {
        var openContextMenu = angular.element(document.getElementById("playerListItem-"+$scope.contextMenuForPlayer)).hasClass("open");
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

                $rootScope.title = 'Mafioso - ' + $scope.game.state;
                if ($scope.currentState != $scope.game.state) {

                    var msg = " ";
                    if ($scope.game.state == 'night') {
                        msg = "Night has fallen!";
                    } else if ($scope.game.state == 'day') {
                        msg = "A new day has begun!";
                    } else if ($scope.game.state == 'villager-win' || $scope.game.state == 'mafia-win') {
                        msg = "Game has ended!";
                    }

                    var newExcitingAlerts = (function () {
                    var timeoutId;
                    var blink = function() { document.title = document.title == $rootScope.title ? msg : $rootScope.title; };
                    var clear = function() {
                        clearInterval(timeoutId);
                        document.title = $rootScope.title;
                        window.onmousemove = null;
                        timeoutId = null;
                    };
                    return function () {
                        if (!timeoutId) {
                            timeoutId = setInterval(blink, 1000);
                            window.onmousemove = clear;
                        }
                    };
                    }());
                    newExcitingAlerts();
                    $scope.currentState = $scope.game.state;
                }
        }

        if ($scope) {
            $scope.$apply();

            if (openContextMenu) {
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

