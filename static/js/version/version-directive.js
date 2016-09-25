'use strict';

angular.module('myApp.version.version-directive', [])

.directive('appVersion', ['version', function(version) {
  return function(scope, elm, attrs) {
    elm.text(version);
  };
}]).directive('qrcode', ['$http','Auth', function ($http,Auth) {

    return {
        restrict: 'E',
        template: '<div class="qr-mini"><a class="qr"><span class="qr-small"></span>二维码</a></div><div class="qr-big"></div>',
        link: link
    };

    function link(scope, element, attrs) {
        $http.get(attrs['url']).then(function(res){
            var data = res.data.Data;
            Auth.setlogin(true);
            Auth.setUserName(data.name);
            $(element).find(".qr-big").qrcode({"size":200,text:"user:"+data.userid});
            $(element).find(".qr-small").qrcode({"size":20,text:"user:"+data.userid});
        });

        element.on('$destroy', function () {
        });
    }
}]).directive('headBanner', ['$location','$route','Auth', function ($location,$route,Auth) {

    return {
        restrict: 'E',
        template: `
        <div class="navbar">
            <div class="navbar-inner">
                <div class="navbar-container">
                    <!-- Navbar Barnd -->
                    <div class="navbar-header pull-left">
                        <a href="#" class="navbar-brand">
                            <small>
                                <img src="img/logo.jpg" alt=""/>
                            </small>
                        </a>
                    </div>
                    <div class="pull-left"  style="width:100px;">
                        &nbsp;
                    </div>
                    <div class="link-header pull-left activeparent" >
                        <a href="#/about" class="navbar-brand">
                            关于傲路
                        </a>
                        <a href="#/solution" class="navbar-brand">
                            解决方案
                        </a>
                        <a href="#/clients" class="navbar-brand">
                            客户分布
                        </a>
                        <a href="#/link" class="navbar-brand">
                            联系方式
                        </a>
                    </div>

                    <div class="right-header activeparent" >
                        <a href="#/upload"  class="upload">上传</a>/<a href="#/manager"  class="">管理</a>
                    </div>
                    <div class="right-login activeparent" >
                        <i class="fa  fa-user "></i><a href="#/login" class="login">登录</a>/<a href="#/register"  class="register">注册</a>
                    </div>
                </div>
            </div>
        </div>
        <div style="clear: both;width:100%;position: relative;z-index: 10;">
            <img style="width:100%;" src="/img/head-bg.jpg">
        </div>
        `,
        link: link
    };

    function link(scope, element, attrs) {
        console.log($route.originalPath);
        console.log($route.current);
        console.log($location.path());
        var path = $location.path();
        $(element).find("div.activeparent >a").each(function(idx,elem){
            var url = $(elem).prop("href");
            var subpath = url.substring(url.indexOf("#")+1);
            console.log(path,subpath,subpath ===path);
            if(subpath ===path){
                $(elem).addClass("active");
            }else{
                $(elem).removeClass("active");
            }
        });
        element.on('$destroy', function () {
        });
    }
}]);
