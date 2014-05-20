'use strict';

angular.module('mafiachat.services', []).factory('WebSocket', ['$q', '$rootScope', 'MsgPublisher', function($q, $rootScope, MsgPublisher) {
    // We return this object to anything injecting our service
    var Service = {};

    // Create our websocket object with the address to the websocket
    var ws = new WebSocket("ws://"+window.location.host+window.location.pathname+"ws/");

    var scope = {};

    ws.onopen = function(){
        console.log("Socket has been opened!");
    };

    ws.onmessage = function(message) {
        listener(JSON.parse(message.data));
    };

    Service.setScope = function($scope) {
        scope = $scope;
    }

    Service.sendMsg = function(request) {
        var defer = $q.defer();
        console.log('Sending request', request);
        ws.send(JSON.stringify(request));
        return defer.promise;
    };

    function listener(data) {
        console.log("Received data from websocket: ", data);
        // A public msg for everyone
        MsgPublisher.publish(scope, data);
    }

    // Define a "getter" for getting customer data
    Service.getCustomers = function() {
        var request = {
            type: "get_customers"
        }
        // Storing in a variable for clarity on what sendRequest returns
        return sendRequest(request);
    }

    return Service;
}]);
