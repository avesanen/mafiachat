'use strict';

angular.module('mafiachat.directives', []).directive('message', ['$rootScope', function($rootScope) {
   return {
       restrict: 'E',
       scope: {
           text: '@'
       },
       templateUrl: '/partials/message.html',
       link: function(scope, element, attrs) {
           if(/^http(s?):\/\/.+\.(gif|png|jpg|jpeg)$/i.test(scope.text)) {
               scope.text = "<a class='imgThumb' target='_blank' href='"+scope.text+"'><img src='" + scope.text + "' height='100' /></a>";
               element.html(scope.text);
           } else if (/^http(s?):\/\/.+\.(webm)$/i.test(scope.text)) {
               scope.text = '<video height="240" autoplay><source src="'+scope.text+'" type="video/webm"></video>';
               element.html(scope.text);
           } else if (/(\b(https?):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/i.test(scope.text)) {
               var replacePattern = /(\b(https?):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/gim;
               element.html(scope.text.replace(replacePattern, '<a href="$1" target="_blank">$1</a>'));
           }
       }
   };
}]);
