'use strict';

angular.module('mafiachat.services', []).factory('WebSocket', ['$q', '$rootScope', 'ResponseHandler', function($q, $rootScope, ResponseHandler) {
    // We return this object to anything injecting our service
    var Service = {};

    // Create our websocket object with the address to the websocket
    var ws = new WebSocket("ws://"+window.location.host+window.location.pathname+"ws/");

    var scope = {};
    var defer;
    var msgQueue = [];
    var connectionOpened = false;

    ws.onopen = function(){
        console.log("Socket has been opened!");
        connectionOpened = true;
        handleQueue();
    };

    ws.onmessage = function(message) {
        listener(JSON.parse(message.data));
    };

    ws.onerror = function(error) {
        console.log("WebSocket error: ", error);
    }

    Service.isConnected = function() {
        return connectionOpened;
    }

    Service.setScope = function($scope) {
        scope = $scope;
    }

    Service.sendDeferMsg = function(request) {
        defer = $q.defer();
        console.log('Sending defer request', request);
        ws.send(JSON.stringify(request));
        return defer.promise;
    };

    Service.sendMsg = function(request) {
        if (!connectionOpened) {
            msgQueue.push(request);
        } else {
            defer = undefined;
            console.log('Sending request', request);
            ws.send(JSON.stringify(request));
        }
    };

    function handleQueue() {
        for (var i = 0; i < msgQueue.length; i++) {
            Service.sendMsg(msgQueue[i]);
        }
    }

    function listener(data) {
        console.log("Received data from websocket: ", data);

        if (defer) {
            $rootScope.$apply(defer.resolve(data));
        }
        // Message should be visible for user
        ResponseHandler.handle(scope, data);


    }

    return Service;
}]);
