'use strict';


app.controller('VideoManagerCtrl', function ($scope, i18nService, $modal, $log,cfpLoadingBar,$http, $sce,$window) {
    //初始化查询参数
    //以下处理是为了防止树内容一次加载导致的性能问题,改为了点击加载 , 使用了非deepcopy
    $scope.clients = [];
    $scope.query = {status: '等待审核'};
    $scope.queryvideo = refresh_video_list;
    $scope.video_verify = video_verify;
    $scope.appendList = function(data){
        console.log(data);
        for(var i = 0 ;i < data.length;i++){
            var row = eval("("+data[i]+")");
            //var row = {_id:item._id,name:item,src:'/uploadvideo/'+item};
            $scope.$parent.videolist.push(row);
            $scope.$parent.videomap[row._id]= row;
            //$scope.$parent.$digest();
        }
    };

    function refresh_video_list() {
        $scope.$parent.initVideo($scope.query);
    }


    function video_verify(c,i){
        if($window.confirm("确定视频\""+ c.name+"\"通过审核吗?")==1){
            $http.get("/verifyer/video/verify/"+ c._id).then(function(ret){
                $scope.$parent.videolist[i].status='正常';
                $window.alert('审核成功!')
                //$scope.$parent.$digest();
            });
        }
    }

    $scope.change_name_cancel = function(c){
        c.newname = c.name;
        c.show=false;
    }

});
