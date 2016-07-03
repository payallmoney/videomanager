'use strict';


app.controller('ClientStatusCtrl', function ($scope, i18nService, $modal, $log, cfpLoadingBar, $http, $window) {
    //初始化查询参数
    $scope.clients = [];



    $scope.queryClients = queryClients;
    queryClients();
    function queryClients() {

        $http.get("/clients").then(function (ret) {
            $scope.clients = ret.data;
            $scope.clients[0].show = true;
        });
    }


    function video_select_change(video){
        console.log(video);
        var currentitem = null;
        for(var i = 0;i<$scope.$parent.videolist.length;i++){
            var item = $scope.$parent.videolist[i];
            if(item._id==video.new_id){
                currentitem = item;
                break;
            }
        }
        if(currentitem != null && currentitem.status != "正常"){
            video.new_id = null;
            alert("视频正在审核, 暂时无法选择 , 请等待审核成功后点击刷新再试!");
        }
    }
    function client_add (name) {
        $http.get("/client/add/" + name).then(function (resp) {
            if(resp.data.success){
                $scope.clients.push({_id: name, name: name});
            }else{
                alert(resp.data.msg);
            }
        });
    }

    function  client_del(c, idx) {
        if ($window.confirm("确定删除设备\"" + c._id + "\"吗?") == 1) {
            $http.get("/client/del/" + c._id).then(function (ret) {
                $scope.clients.splice(idx, 1);
            });
        }
    }
    function client_unbind (c, idx) {
        if ($window.confirm("确定解除设备\"" + c._id + "\"的绑定吗?") == 1) {
            $http.get("/client/unbind/" + c._id).then(function (ret) {
                $scope.clients.splice(idx, 1);
            });
        }
    }

    function video_add (c) {
        c.show = true;
        if (!c.videolist) {
            c.videolist = [];
        }
        c.videolist.push({show: true});
        console.log(c)
    }
    function video_save (c, v,idx) {
        console.log(c, v);
        if (v.new_id && v.new_id != v._id) {
            if (!v._id) {
                $http.get("/client/videoadd/" + c._id + '/' + v.new_id).then(function (ret) {
                    v._id = v.new_id;
                    v.show = false;
                });
            } else {
                $http.get("/client/videochange/" + c._id + '/' + idx + '/' + v.new_id).then(function (ret) {
                    v._id = v.new_id;
                    v.show = false;
                });
            }
        } else {
            if (v._id) {
                v.show = false;
            }
        }
    }
    function  video_save_cancel(c, v) {
        v.new_id = v._id;
        v.show = false;
    }

    function video_del(c, v, idx) {
        console.log(c, v, idx);
        var v = c.videolist[idx];

        if (v._id && $scope.$parent.videomap[v._id]) {
            console.log($scope.$parent.videomap[v._id]);
            if ($window.confirm("确定将视频\"" + $scope.$parent.videomap[v._id].name + "\"从列表中移除吗?") == 1) {
                $http.get("/client/videodel/" + c._id + '/' + idx).then(function (ret) {
                    c.videolist.splice(idx, 1);
                });
            }
        } else {
            $http.get("/client/videodel/" + c._id + '/' + idx).then(function (ret) {
                c.videolist.splice(idx, 1);
            });
        }
    }

});
