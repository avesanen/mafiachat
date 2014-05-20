'use strict';

angular.module('mafiachat.services').factory('MsgPublisher', ['$q', '$rootScope', function($q, $rootScope) {
    // We return this object to anything injecting our service
    var Service = {};

    Service.publish = function($scope, msg) {
        if ($scope.log != "") {
            $scope.log += "<br />";
        }

        var now = new Date();
        $scope.log += ('0'  + now.getHours()).slice(-2)+':'+('0' + now.getMinutes()).slice(-2) + " ";

        switch (msg.msgType) {
            case 'joinGame':
                $scope.players.push(msg.data.name);
                $scope.log += "<b>" + msg.data.name + " joined the game!</b>";
                break;
            case 'chatMessage':
                $scope.log += "<span class='"+msg.data.faction+"Message'>" + msg.data.message + "</span>";
                break;
        }
        $scope.$apply();

        var objDiv = document.getElementById("log");
        objDiv.scrollTop = objDiv.scrollHeight;
    };

    return Service;
}]);

