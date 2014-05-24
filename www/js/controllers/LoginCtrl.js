'use strict';

angular.module('mafiachat.controllers').controller('LoginCtrl', ['$rootScope', '$scope', '$location', 'WebSocket', 'ResponseHandler', function($rootScope, $scope, $location, WebSocket, ResponseHandler) {
    WebSocket.setScope($rootScope);

    $scope.login = function() {

        if (!$scope.name || !$scope.password) {
            $scope.errorMsg = "Enter all data.";
            return;
        }

        var message = {data:{}};
        message.msgType = 'login';
        message.data.name = $scope.name;
        message.data.password = $scope.password;
        if ($rootScope.gameId) {
            message.data.gameId = $rootScope.gameId;
        }

        WebSocket.sendDeferMsg(message).
            then(function(resp) {
                $rootScope.name = $scope.name;
                if ($rootScope.requiredPath) {
                    $location.path($rootScope.requiredPath);
                } else {
                    $location.path("/game");
                }
            }, function(error) {
                $scope.errorMsg = "Couldn't connect to backend :(";
            }
        );
    }

}]);

