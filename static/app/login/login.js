'use strict';

angular.module('videosystem.login', ['ngRoute'])

    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/login', {
            templateUrl: '/app/login/login-json.html',
            //templateUrl: 'app/login/login.html',
            controller: 'LoginCtrl'
        });
    }])

    .controller('LoginCtrl', function ($scope, Auth, $location, $http) {
        $scope.getValideCode = function () {
            var codeArray = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
                'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
                'W', 'X', 'Y', 'Z'];
            var code = '';
            for (var i = 0; i < 4; i++) {
                code = code + codeArray[Math.round(Math.random() * 25)];
            }
            $scope.valideCode = code;
        };
        $scope.getValideCode();
        $scope.data = {};
        $scope.msg = '';
        $scope.login = function () {
            console.log($scope.data);
            $http({
                url: '/videosystem/login',
                method: 'POST',
                data: $.param($scope.data),
                responseType: 'text',
                cache: false,
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8'

                },
                transformResponse: appendTransform($http.defaults.transformResponse, function (value) {
                    if (value === "登录成功") {
                        return {"success": true, "msg": "登录成功!"}
                    } else {
                        return {"success": false, "msg": "登录失败!"}
                    }
                    return value;
                })
            }).success(function (data, status, headers, config) {
                console.log(data);

                if (data.success) {
                    Auth.setlogin(true);
                    $location.path('/main');
                } else {
                    $scope.msg = data.msg;
                }
            }).error(function (data, status, headers, config) {
                $scope.msg = "登录失败";
            });
        }
        function appendTransform(defaults, transform) {

            // We can't guarantee that the default transformation is an array
            defaults = angular.isArray(defaults) ? defaults : [defaults];

            // Append the new transformation to the defaults
            return defaults.concat(transform);
        }

    });