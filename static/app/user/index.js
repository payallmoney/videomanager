'use strict';

// Declare app level module which depends on views, and components
var app =
    angular.module('videosystem', [
        'ngRoute',
        'ui.grid',
        'ui.grid.autoResize',
        'ui.grid.resizeColumns',
        'ui.grid.pagination',
        'ui.grid.selection',
        'ui.utils',
        'cfp.loadingBar',
        'videosystem.login',
        'videosystem.register',
        'videosystem.main',
        'videosystem.filters',
        'myApp.version.version-directive'
    ]).config(['$routeProvider', '$sceProvider', '$httpProvider', function ($routeProvider, $sceProvider, $httpProvider) {
        $httpProvider.interceptors.push(function ($timeout, $q, $injector) {
            var Auth, $http, $location;
            // this trick must be done so that we don't receive
            // `Uncaught Error: [$injector:cdep] Circular dependency found`
            $timeout(function () {
                Auth = $injector.get('Auth');
                $http = $injector.get('$http');
                $location = $injector.get('$location');

            });

            return {
                response: function (rejection) {
                    console.log(rejection);
                    var deferred = $q.defer();
                    if (rejection.status == 401) {
                        Auth.logindlg()
                            .then(function () {
                                deferred.resolve();
                                $state.go(toState.name, toParams);
                            }).catch(function () {
                            deferred.reject(rejection);
                            $location.path('/login');
                        });
                    } else {
                        deferred.resolve(rejection);
                    }
                    return deferred.promise;
                }
            };
        });
        $routeProvider.otherwise({redirectTo: '/login'});
    }]).config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/main', {
            templateUrl: '/app/user/main/main.html',
            controller: 'MainCtrl'
        }).when('/login', {
            templateUrl: '/app/user/login/login.html',
            controller: 'LoginCtrl'
        }).when('/register', {
            templateUrl: '/app/user/register/page.html',
            controller: 'RegisterCtrl'
        });
        $routeProvider.otherwise({redirectTo: '/login'});
    }]).config(function ($controllerProvider, $compileProvider, $filterProvider, $provide) {
        app.controller = $controllerProvider.register;
    }).factory('Auth', function ($http, $location, $q) {
        var logined;
        var name;
        var data = {name:null}
        return {
            setlogin: function (islogined) {
                logined = islogined;
            },
            logined:function(){
                return logined
            },
            getUserName: function(){
                return data.name;
            } ,
            userinfo:data,
            name :data.name,
            setUserName: function ( setname) {
                data.name = setname;
            }
        }
    }).controller('ModalInstanceCtrl', function ($scope, $modalInstance, data) {
        $scope.data = data;
        $scope.ok = function () {
            $modalInstance.close(data);
        };

        $scope.cancel = function () {
            $modalInstance.dismiss('cancel');
        };
    })