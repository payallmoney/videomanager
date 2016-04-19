'use strict';
app.controller('WorkerPageCtrl', function ($scope, i18nService, $modal, $log, cfpLoadingBar, $http, $window) {
    $scope.worker = {};
    $scope.msg = "";
    $scope.worker_save = worker_save;
    $scope.setWorker = setWorker;
    $scope.setEdit = setEdit;

    function worker_save() {
        if(valid()){
            return $http.post("/admin/worker/add",$scope.worker)
        }
    }
    function valid(){
        if(!$scope.worker.account){
            $scope.msg = "帐号不能为空";
            return false;
        }
        if(!$scope.worker.password){
            $scope.msg = "密码不能为空";
            return false;
        }
        if($scope.worker.password !== $scope.worker.password_rp){
            $scope.msg = "两次密码不匹配不能为空";
            return false;
        }
        $scope.msg = "";
        return true;
    }
    function setWorker(worker){
        $scope.worker = worker;
    }
    function setEdit(edit){
        $scope.edit = edit;
    }


});

app.controller('WorkerCtrl', function ($scope, i18nService, $modal, $log, cfpLoadingBar, $http, $window) {
    //初始化查询参数
    //以下处理是为了防止树内容一次加载导致的性能问题,改为了点击加载 , 使用了非deepcopy
    $scope.clients = [];

    init();
    $("#worker_dialog").dialog({
        resizable: false,
        maxHeight: 200,
        width: 300,
        autoOpen: false,
        modal: true,
        height: 300,
        buttons: {
            "保存": function () {
                var dialog =  $(this);
                $("#worker_dialog").scope().worker_save().then(function(resp){
                    console.log(resp);
                    if(resp.data.Success){
                        dialog.dialog("close");
                        init();
                    }
                })
            },
            "关闭": function () {
                $(this).dialog("close");
            }
        },
        close: function (event, ui) {
        }
    });
    $scope.worker_del = worker_del;
    $scope.worker_edit = worker_edit;
    $scope.worker_add = worker_add;

    function worker_add() {
        $("#worker_dialog").scope().setWorker({});
        $("#worker_dialog").scope().setEdit(false);
        $("#worker_dialog").dialog("option", "title", "新增工作人员");
        $("#worker_dialog").dialog("open");
    }

    function worker_edit(idx){
        var worker = $scope.workers[idx];
        $("#worker_dialog").scope().setWorker(worker);
        $("#worker_dialog").scope().setEdit(true);
        $("#worker_dialog").dialog("option", "title", "修改工作人员");
        $("#worker_dialog").dialog("open");
    }
    function worker_del(idx){
        var worker = $scope.workers[idx];
        if ($window.confirm("确定删除工作人员\"" + worker._id + "\"吗?") == 1) {
            $http.post("/admin/worker/del",worker).then(function (ret) {
                $scope.workers.splice(idx, 1);
            });
        }
    }

    function init(){
        $http.get("/admin/workers").then(function (ret) {
            $scope.workers = ret.data;
        });
    }

});

