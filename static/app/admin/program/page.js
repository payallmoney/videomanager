'use strict';


app.controller('VideoClientCtrl', function ($scope, i18nService, $modal, $log,cfpLoadingBar,$http, $sce,$window) {
    //初始化查询参数
    //以下处理是为了防止树内容一次加载导致的性能问题,改为了点击加载 , 使用了非deepcopy
    $scope.programs = [];
    $http.get("/video/program/list").then(function(ret){
        console.log(ret.data);
        $scope.programs = ret.data;
    });

    $scope.program_del = function(c,i){
        if($window.confirm("确定删除版本\""+ c._id+"\"吗?")==1){
            $http.get("/video/program/delete/"+ c._id).then(function(ret){
                $scope.programs.splice(i,1);
            });
        }
    };

    $scope.program_append = function(data){
        for(var i = 0 ;i < data.length;i++){
            if(data[i]){
                var item = eval("("+data[i]+")");
                $scope.programs.push(item);
            }
        }
    }
    $scope.url = function(url){
        return $sce.trustAsUrl(""+url)
    }


});
