'use strict';

angular.module('mafiachat.controllers', []).controller('MainCtrl', ['$rootScope', '$scope', '$location', 'WebSocket', 'ResponseHandler', function($rootScope, $scope, $location, WebSocket, ResponseHandler) {
    if (!$rootScope.games) {
        $rootScope.games = [
            {"id":1, "name":"Jea tässä olis yks peli.", "state":"open", "cops":3, "mafiosi":3, "doctors": 2, "needsPassword":"yes", "maxPlayers":10, "players":[{"name":"juki"}, {"name":"antti"}]},
            {"id":2, "name":"Menossa oleva peli.", "state":"ongoing", "cops":3, "mafiosi":3, "doctors": 2, "needsPassword":"no", "maxPlayers":5, "players":[{"name":"juki"}, {"name":"antti"}, {"name":"jaakko"}, {"name":"juuso"}, {"name":"jani"}]},
            {"id":3, "name":"Avoin peli 1.", "state":"open", "cops":3, "mafiosi":3, "doctors": 2, "needsPassword":"yes", "maxPlayers":100, "players":[]},
            {"id":4, "name":"Avoin peli 2.", "state":"open", "cops":3, "mafiosi":3, "doctors": 2, "needsPassword":"no", "maxPlayers":5, "players":[]}
        ];
    }

    WebSocket.setScope($scope);

    $scope.joinGame = function(gameId, gamePassword) {
        console.log("Joining game: " + gameId);
        if (!$rootScope.name) {
            $rootScope.gameId = gameId;
            $location.path("/login");
        }

        var message = {data:{}};
        message.msgType = 'login';
        message.data.player = $rootScope.name;
        message.data.gameId = gameId;
        message.data.password = gamePassword;

        WebSocket.sendDeferMsg(message).
            then(function(resp) {
                $rootScope.gameInfo = resp.data;
                $location.path("/game");
            }, function(error) {
                $scope.errorMsg = "Couldn't connect to backend :(";
            }
        );
    }

    $scope.newGame = function() {
        $location.path("/createGame");
    }
}]);

