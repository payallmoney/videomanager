'use strict';
app.controller('VerifyerPageCtrl', function ($scope, i18nService, $modal, $log, cfpLoadingBar, $http, $q) {
    $scope.verifyer = {};
    $scope.msg = "";
    $scope.verifyer_save = verifyer_save;
    $scope.setverifyer = setverifyer;
    $scope.setEdit = setEdit;

    function verifyer_save() {
        if(valid()){
            return $http.post("/admin/verifyer/add",$scope.verifyer)
        }
    }
    function valid(){
        if(!$scope.verifyer.account){
            $scope.msg = "帐号不能为空!";
            return false;
        }
        if(!$scope.verifyer.password){
            $scope.msg = "密码不能为空!";
            return false;
        }
        if($scope.verifyer.password !== $scope.verifyer.password_rp){
            $scope.msg = "两次密码不匹配!";
            return false;
        }
        $scope.msg = "";
        return true;
    }
    function setverifyer(verifyer){
        $scope.verifyer = verifyer;
    }
    function setEdit(edit){
        $scope.edit = edit;
    }


});

app.controller('VerifyerCtrl', function ($scope, i18nService, $modal, $log, cfpLoadingBar, $http, $window) {
    //初始化查询参数
    //以下处理是为了防止树内容一次加载导致的性能问题,改为了点击加载 , 使用了非deepcopy
    $scope.clients = [];

    init();
    $("#verifyer_dialog").dialog({
        resizable: false,
        maxHeight: 200,
        width: 300,
        autoOpen: false,
        modal: true,
        height: 300,
        buttons: {
            "保存": function () {
                var dialog =  $(this);
                var save = $("#verifyer_dialog").scope().verifyer_save();
                if(save){
                    save.then(function(resp){
                        console.log(resp);
                        if(resp.data.Success){
                            dialog.dialog("close");
                            init();
                        }
                    })
                }else{
                    $("#verifyer_dialog").scope().$digest();
                }
            },
            "关闭": function () {
                $(this).dialog("close");
            }
        },
        close: function (event, ui) {
        }
    });
    $scope.verifyer_del = verifyer_del;
    $scope.verifyer_edit = verifyer_edit;
    $scope.verifyer_add = verifyer_add;

    function verifyer_add() {
        $("#verifyer_dialog").scope().setverifyer({});
        $("#verifyer_dialog").scope().setEdit(false);
        $("#verifyer_dialog").dialog("option", "title", "新增审核人员");
        $("#verifyer_dialog").dialog("open");
    }

    function verifyer_edit(idx){
        var verifyer = $scope.verifyers[idx];
        $("#verifyer_dialog").scope().setverifyer(verifyer);
        $("#verifyer_dialog").scope().setEdit(true);
        $("#verifyer_dialog").dialog("option", "title", "修改审核人员");
        $("#verifyer_dialog").dialog("open");
    }
    function verifyer_del(idx){
        var verifyer = $scope.verifyers[idx];
        if ($window.confirm("确定删除审核人员\"" + verifyer._id + "\"吗?") == 1) {
            $http.post("/admin/verifyer/del",verifyer).then(function (ret) {
                $scope.verifyers.splice(idx, 1);
            });
        }
    }

    function init(){
        $http.get("/admin/verifyers").then(function (ret) {
            $scope.verifyers = ret.data;
        });
    }

});

