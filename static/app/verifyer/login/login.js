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
        $scope.next = next;
        $scope.login = login;
        $scope.enterLogin = enterLogin;
        //$scope.register = register;
        function login() {
            if(!valid()){
                return;
            }
            $http({
                url: '/verifyer/login',
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

        }



        function next(){
            if (event.keyCode == 13) {
                $(event.target).parent().parent().next().find("input:visible")[0].focus();
                event.preventDefault();
            }
        }
        function enterLogin(){
            if (event.keyCode == 13) {
                login()
            }
        }


        function valid(){
            var data = $scope.data;

            if(!data.userid){
                $scope.msg = "帐号不能为空!";
                return false;
            }

            if(!data.password){
                $scope.msg = "密码不能为空!";
                return false;
            }

            if(data.code.toLowerCase() !== $scope.valideCode.toLowerCase()){
                $scope.msg = "验证码输入错误!";
                return false;
            }
            return true;
        }

    });