'use strict';

angular.module('myApp.version.version-directive', [])

.directive('appVersion', ['version', function(version) {
  return function(scope, elm, attrs) {
    elm.text(version);
  };
}]).directive('qrcode', ['$http','Auth', function ($http,Auth) {

    return {
        restrict: 'E',
        template: '<div class="qr-mini"><a style="cursor: pointer;text-decoration: none;color: #fff;">显示二维码</a></div><div class="qr-big"></div>',
        link: link
    };

    function link(scope, element, attrs) {
        $http.get(attrs['url']).then(function(res){
            var data = res.data.Data;
            Auth.setlogin(true);
            Auth.setUserName(data.name);
            $(element).find(".qr-big").qrcode({"size":200,text:"user:"+data.userid});
        });

        element.on('$destroy', function () {
        });
    }
}]);
