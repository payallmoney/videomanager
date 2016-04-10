'use strict';

angular.module('videosystem.login', ['ngRoute'])
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
                url: '/admin/login',
                method: 'POST',
                //param: $.param($scope.data),
                data: ($scope.data),
                cache: false
            }).success(function (data, status, headers, config) {
                console.log(data);
                if (data.Success) {
                    Auth.setlogin(true);
                    Auth.setUserName(data.Data.name);
                    $location.path('/main');
                } else {
                    $scope.msg = data.Msg;
                }
            }).error(function (data, status, headers, config) {
                console.log(data);
                $scope.msg = "登录失败";
            });
        };

    });