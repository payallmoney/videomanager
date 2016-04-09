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
        'videosystem.main',
        'videosystem.filters'
    ]).config(['$routeProvider','$sceProvider', function ($routeProvider,$sceProvider) {
        $sceProvider.enabled(false);
        $routeProvider.otherwise({redirectTo: '/main'});
    }]).config(function ($controllerProvider, $compileProvider, $filterProvider, $provide) {
        app.controller = $controllerProvider.register;
    }).factory('Auth', function ($http, $location, $q) {
        var logined;
        var user;
        return {
            setlogin: function (islogined) {
                logined = islogined;
            },
            getUser: function () {
                return user;
            },
            checklogin: function () {
                var deferred = $q.defer();
                if (!logined) {
                    TaskService.getCurrentUser({
                        callback: function (data) {
                            console.log("data==", data)
                            if (data == null) {
                                $location.path('/videosystem/main');
                                deferred.reject(data);
                            } else {
                                user = data;
                                deferred.resolve(data);
                            }
                        }, errorHandler: function (message) {
                            $location.path('/videosystem/main');
                            deferred.reject(message);
                        }
                    });
                } else {
                    TaskService.getCurrentUser({
                        callback: function (data) {
                            user = data;
                            deferred.resolve(data);
                        }, errorHandler: function (message) {
                            $location.path('/videosystem/main');
                            deferred.reject(message);
                        }
                    });
                }
                return deferred.promise;
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
    });