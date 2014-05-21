'use strict';

angular.module('mafiachat.controllers').controller('GameCtrl', ['$rootScope', '$scope', '$location', 'WebSocket', 'ResponseHandler', function($rootScope, $scope, $location, WebSocket, ResponseHandler) {

    $scope.createGame = function() {
        console.log("Create game: ", $scope.gameName);
        var game = {"id":$scope.games.length, "name":$scope.gameName};
        $rootScope.games.push(game);
        $location.path("/lobby");
    }

}]);

