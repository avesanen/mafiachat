'use strict';

angular.module('mafiachat.controllers').controller('LoginCtrl', ['$rootScope', '$scope', '$location', '$routeParams', 'WebSocket', function($rootScope, $scope, $location, $routeParams, WebSocket) {
    WebSocket.setScope($rootScope);

    $scope.login = function() {

        if (!$scope.name || !$scope.password) {
            $scope.errorMsg = "Enter all data.";
            return;
        }

        var message = {data:{}};
        message.msgType = 'loginMessage';
        message.data.name = $scope.name;
        message.data.password = $scope.password;

        WebSocket.sendDeferMsg(message).
            then(function(resp) {
                sessionStorage.name = $scope.name;
                sessionStorage.pass = $scope.password;
                $rootScope.name = $scope.name;
                if ($rootScope.requiredPath) {
                    $location.path($rootScope.requiredPath);
                } else {
                    $location.$$search = {};
                    $location.path("/game");
                }
            }, function(error) {
                $scope.errorMsg = "Couldn't connect to backend :(";
            }
        );
    }

}]);

