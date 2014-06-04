'use strict';

angular.module('mafiachat.controllers', []).controller('MainCtrl', ['$rootScope', '$scope', '$location', '$timeout', '$http', '$window', 'GameService', function($rootScope, $scope, $location, $timeout, $http, $window, GameService) {
    $scope.invalidPass = true;

    $http({method: 'GET', url: '/games.json'}).
        success(function(data, status, headers, config) {
            $scope.games = data;
        }).
        error(function(data, status, headers, config) {
            console.log("Unable to retrieve games. Error: ", data);
        });

    var timer = false;

    $scope.gameAuth = function() {
        var scope = this;

        if (timer) {
            $timeout.cancel(timer);
        }

        if (scope.gamePassword == '') {
            return;
        }

        timer = $timeout(function() {
            var message = {data:{}};
            message.msgType = 'gameAuth';
            message.data.name = scope.game.id;
            message.data.password = scope.gamePassword;
            WebSocket.sendDeferMsg(message).
                then(function(resp) {
                    if (resp.data.success) {
                        scope.errorMsg = "";
                        scope.invalidPass = false;
                    } else {
                        scope.errorMsg = "Invalid password";
                        scope.invalidPass = true;
                    }
                }, function(error) {
                    scope.invalidPass = true;
                }
            );
        }, 500);
    }

    $scope.joinGame = function() {
        console.log("Joining game: " + this.game.id);
        $window.location.href = "/g/"+this.game.id+"/#game";

    }

    $scope.newGame = function() {
        $window.location.href = "/g/";
    }

    $scope.alivePlayers = function() {
        var count = 0;
        for (var i in this.game.players) {
            if (this.game.players[i].faction != 'ghost') {
                count++;
            }
        }
        return count;
    }
}]);

