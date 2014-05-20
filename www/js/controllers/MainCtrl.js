'use strict';

angular.module('mafiachat.controllers', []).controller('MainCtrl', ['$scope', 'WebSocket', 'MsgPublisher', function($scope, WebSocket, MsgPublisher) {
    $scope.log = "";
    $scope.players = [];
    WebSocket.setScope($scope);

    $scope.sendMsg = function() {
        if (!$scope.msg) {
            return false;
        }

        var message = {data:{}};
        var type = $scope.msgType;

        if (type === 'villagerChat') {
            message.msgType = 'chatMessage';
            message.data.faction = 'villager';
            message.data.message = $scope.msg;
        }
        if (type === 'mafiaChat') {
            message.msgType = 'chatMessage';
            message.data.faction = 'mafia';
            message.data.message = $scope.msg;
        }
        if (type === 'joinGame') {
            message.msgType = 'joinGame';
            message.data.name = $scope.msg;
        }

        WebSocket.sendMsg(message).
            then(function(resp) {
                MsgPublisher.publish($scope, resp);
            }, function(error) {
                console.log("Failed to send data...");
            }
        );
        $scope.msg = "";
    }
}]);

