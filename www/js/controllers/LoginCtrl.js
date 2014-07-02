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

                if (resp.msgType == "loginFailed") {
                    if (resp.reason == "wrongPassword") {
                        $scope.errorMsg = "User already exists in this game with give name but password was wrong.";
                        return;
                    } else if (resp.reason == "alreadyLoggedIn") {
                        $scope.errorMsg = "User already logged in and online.";
                        return;
                    }
                }

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

