'use strict';


app.controller('ClientManagerCtrl', function ($scope, i18nService, $modal, $log, cfpLoadingBar, $http, $window) {
    //初始化查询参数
    //以下处理是为了防止树内容一次加载导致的性能问题,改为了点击加载 , 使用了非deepcopy
    $scope.clients = [];
    $http.get("/video/clients").then(function (ret) {
        $scope.clients = ret.data;
        console.log(ret.data);
        $scope.clients[0].show = true;
        //for (var i = 0; i < $scope.clients.length; i++) {
        //    var row = $scope.clients[i];
        //    if (row.videolist) {
        //        row.videolist = [];
        //        for (var j = 0; j < row.videolist.length; j++) {
        //            var item = row.videolist[j];
        //            row.videolist.push({_id: item, new_id: item, show: false});
        //        }
        //    }
        //}
    });
    $scope.client_add = function (name) {
        $http.get("/video/client/add/" + name).then(function (ret) {
            $scope.clients.push({_id: name, name: name});
        });
    };

    $scope.client_del = function (c, idx) {
        if ($window.confirm("确定删除设备\"" + c._id + "\"吗?") == 1) {
            $http.get("/video/client/del/" + c._id).then(function (ret) {
                $scope.clients.splice(idx, 1);
            });
        }
    };

    $scope.video_add = function (c) {
        c.show = true;
        if (!c.videolist) {
            c.videolist = [];
        }
        c.videolist.push({show: true});
        console.log(c)
    };
    $scope.video_save = function (c, v,idx) {
        console.log(c, v);
        if (v.new_id && v.new_id != v._id) {
            if (!v._id) {
                $http.get("/video/client/videoadd/" + c._id + '/' + v.new_id).then(function (ret) {
                    v._id = v.new_id;
                    v.show = false;
                });
            } else {
                $http.get("/video/client/videochange/" + c._id + '/' + idx + '/' + v.new_id).then(function (ret) {
                    v._id = v.new_id;
                    v.show = false;
                });
            }
        } else {
            if (v._id) {
                v.show = false;
            }
        }
    };
    $scope.video_save_cancel = function (c, v) {
        v.new_id = v._id;
        v.show = false;
    };

    $scope.video_del = function (c, v, idx) {
        console.log(c, v, idx);
        var v = c.videolist[idx];

        if (v._id) {
            console.log($scope.$parent.videomap[v._id]);
            if ($window.confirm("确定将视频\"" + $scope.$parent.videomap[v._id].name + "\"从列表中移除吗?") == 1) {
                $http.get("/video/client/videodel/" + c._id + '/' + idx).then(function (ret) {
                    c.videolist.splice(idx, 1);
                });
            }
        } else {
            c.videolist.splice(idx, 1);
        }
    }

});
