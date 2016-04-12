'use strict';


app.controller('VideoManagerCtrl', function ($scope, i18nService, $modal, $log,cfpLoadingBar,$http, $sce,$window) {
    //初始化查询参数
    //以下处理是为了防止树内容一次加载导致的性能问题,改为了点击加载 , 使用了非deepcopy
    $scope.clients = [];

    $scope.appendList = function(data){
        console.log(data);
        for(var i = 0 ;i < data.length;i++){
            var item = data[i];
            var row = {_id:item,name:item,src:$sce.trustAsUrl('/uploadvideo/'+item)};
            $scope.$parent.videolist.push(row);
            $scope.$parent.videomap[row._id]= row;
            $scope.$parent.$digest();
        }
    }

    $scope.video_del = function(c,i){
        if($window.confirm("确定删除视频\""+ c.name+"\"吗?")==1){
            $http.get("/video/del/"+ c._id).then(function(ret){
                $scope.$parent.videolist.splice(i,1);
                $scope.$parent.$digest();
            });
        }
    }
    $scope.change_name = function(c){
        $http.get("/video/changename/"+ c._id+'/'+ c.newname).then(function(ret){
            c.name = c.newname;
            c.show=false;
            $scope.$parent.videomap[c._id]= c;
            $scope.$parent.$digest();
        });
    }
    $scope.change_name_cancel = function(c){
        c.newname = c.name;
        c.show=false;
    }

});
