'use strict';

angular.module('mafiachat.controllers', []).controller('MainCtrl', ['$rootScope', '$scope', '$location', '$timeout', 'WebSocket', 'GameService', function($rootScope, $scope, $location, $timeout, WebSocket, GameService) {
    if (!$rootScope.games) {
        $rootScope.games = [
            {"id":"23235-235234-23423", "name":"Jea tässä olis yks peli.", "state":"open", "cops":3, "mafiosi":3, "doctors": 2, "needsPassword":true, "maxPlayers":10, "players":[{"name":"juki"}, {"name":"antti"}]},
            {"id":"23235-235234-87345", "name":"Menossa oleva peli.", "state":"ongoing", "cops":3, "mafiosi":3, "doctors": 2, "needsPassword":false, "maxPlayers":5, "players":[{"name":"juki"}, {"name":"antti"}, {"name":"jaakko"}, {"name":"juuso"}, {"name":"jani"}]},
            {"id":"23235-235234-77543", "name":"Avoin peli 1.", "state":"open", "cops":3, "mafiosi":3, "doctors": 2, "needsPassword":true, "maxPlayers":100, "players":[]},
            {"id":"23235-235234-56477", "name":"Avoin peli 2.", "state":"open", "cops":3, "mafiosi":3, "doctors": 2, "needsPassword":false, "maxPlayers":5, "players":[]}
        ];
    }

    WebSocket.setScope($scope);
    $scope.invalidPass = true;

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

        if (!$rootScope.name && !sessionStorage.name) {
            GameService.setGameData(this.game.id);
            $location.path("/login");
            return;
        } else if (!$rootScope.name) {
            $rootScope.name = sessionStorage.name;
        }

        var message = {data:{}};
        message.msgType = 'login';
        message.data.player = $rootScope.name;
        message.data.gameId = this.game.id;

        WebSocket.sendDeferMsg(message).
            then(function(resp) {
                $rootScope.gameInfo = resp.data;
                $location.path("/game");
            }, function(error) {
                $scope.errorMsg = error;
            }
        );
    }

    $scope.newGame = function() {
        $location.path("/createGame");
    }
}]);

