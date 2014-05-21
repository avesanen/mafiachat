'use strict';

angular.module('mafiachat.controllers').controller('LoginCtrl', ['$rootScope', '$scope', '$location', 'WebSocket', 'ResponseHandler', function($rootScope, $scope, $location, WebSocket, ResponseHandler) {

    $scope.login = function() {

        if (!$scope.name || !$scope.password) {
            $scope.errorMsg = "Enter all data.";
            return;
        }

        var message = {data:{}};
        message.msgType = 'login';
        message.data.name = $scope.name;
        message.data.password = $scope.password;

        WebSocket.sendDeferMsg(message).
            then(function(resp) {
                $rootScope.name = $scope.name;
                $location.path("/lobby");
            }, function(error) {
                $scope.errorMsg = "Couldn't connect to backend :(";
            }
        );
    }

}]);

