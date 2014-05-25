'use strict';

angular.module('mafiachat.services').factory('GameService', ['$q', '$rootScope', function($q, $rootScope) {
    // We return this object to anything injecting our service
    var Service = {};
    var gameId;
    var gamePassword;

    Service.setGameData = function(id, password) {
        gameId = id
        gamePassword = password;
    }

    Service.getGameData = function() {
        return {"id":gameId, "password": gamePassword};
    }

    return Service;
}]);

