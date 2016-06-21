'use strict';

angular.module('videosystem.main', ['ngRoute', 'ui.bootstrap'])
    .controller('MainCtrl', function ($scope, Auth, $location, $q, $http, $sce) {
        $scope.userInfo = Auth.userinfo;
        $scope.username = $scope.userInfo.name;
        $scope.tabs = [];
        $scope.activetext = '';
        var init = $q.defer();
        $scope.init = function () {
            return init.promise;
        };
        //视频列表
        $scope.videolist = [];
        $scope.videomap = {};
        $scope.initVideo = initVideo;


        var lastActive;
        $scope.loadTab = loadTab;
        $scope.closeTab = closeTab;

        $scope.menu = [
            {
                "text": '视频审核',
                "js": "/app/verifyer/videomanager/videomanager.js",
                'html': '/app/verifyer/videomanager/videomanager.html'
            }];
        $scope.loadTab($scope.menu[0]);
        initPage();

        function initVideo(param){
            $http.get("/verifyer/video/list",{params:param}).then(function (ret) {
                $scope.videolist = ret.data;
                $scope.videomap = {};
                for (var i = 0; i < $scope.videolist.length; i++) {
                    var row = $scope.videolist[i];
                    $scope.videomap[row._id] = row;
                }
            });
        }

        function closeTab(idx) {
            console.log(idx);
            console.log($scope.tabs);
            $scope.tabs.splice(idx, 1);
            if ($scope.tabs[idx]) {
                $scope.tabs[idx].active = true;
            } else if ($scope.tabs[idx - 1]) {
                $scope.tabs[idx - 1].active = true;
            }
            console.log($scope.tabs);
            event.preventDefault();
        }

        function loadTab(menu) {
            //console.log(menu);
            if (menu.js) {
                $script(menu.js, function () {
                    var flag = false;
                    for (var i = 0; i < $scope.tabs.length; i++) {
                        if ($scope.tabs[i].text == menu.text) {
                            flag = true;

                        } else {
                            $scope.tabs[i].active = false;
                        }
                    }
                    if (!flag) {
                        menu.close = true;
                        $scope.tabs.push(menu)
                    }
                    menu.active = true;
                    $scope.$digest();
                });
            } else {
                var flag = false;
                for (var i = 0; i < $scope.tabs.length; i++) {
                    if ($scope.tabs[i].text == menu.text) {
                        flag = true;

                    } else {
                        $scope.tabs[i].active = false;
                    }
                }
                if (!flag) {
                    menu.close = true;
                    $scope.tabs.push(menu)
                }
                menu.active = true;
            }
        }

        function initPage(){
            initVideo({status: '等待审核'});
            $("#sidebar-collapse").on('click', function () {
                if (!$('#sidebar').is(':visible'))
                    $("#sidebar").toggleClass("hide");
                $("#sidebar").toggleClass("menu-compact");
                $(".sidebar-collapse").toggleClass("active");
                var b = $("#sidebar").hasClass("menu-compact");

                if ($(".sidebar-menu").closest("div").hasClass("slimScrollDiv")) {
                    $(".sidebar-menu").slimScroll({destroy: true});
                    $(".sidebar-menu").attr('style', '');
                }
                if (b) {
                    $(".open > .submenu")
                        .removeClass("open");
                } else {
                    if ($('.page-sidebar').hasClass('sidebar-fixed')) {
                        var position = (readCookie("rtl-support") || location.pathname == "/index-rtl-fa.html" || location.pathname == "/index-rtl-ar.html") ? 'right' : 'left';
                        $('.sidebar-menu').slimscroll({
                            height: 'auto',
                            position: position,
                            size: '3px',
                            color: themeprimary
                        });
                    }
                }
                //Slim Scroll Handle
            });
            //$("#qr-mini").qrcode({"size":30,text:$scope.username});

        }


    });