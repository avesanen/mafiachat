'use strict';

angular.module('mafiachat.directives', []).directive('message', ['$rootScope', function($rootScope) {
   return {
       restrict: 'E',
       scope: {
           text: '@'
       },
       templateUrl: '/partials/message.html',
       link: function(scope, element, attrs, tabsCtrl) {
           if(/^http:\/\/.+\.(gif|png|jpg|jpeg)$/i.test(scope.text)) {
               scope.text = "<a target=\"_blank\" href=\""+scope.text+"\"><img src=\"" + scope.text + "\" width=\"100\" /></a>";
               element.html(scope.text);
           }

       }
   };
}]);