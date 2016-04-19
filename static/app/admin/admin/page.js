'use strict';
app.controller('AdminPageCtrl', function ($scope, i18nService, $modal, $log, cfpLoadingBar, $http, $window) {
    $scope.user = {};
    $scope.msg = "";
    $scope.admin_save = admin_save;
    $scope.setUser = setUser;
    $scope.setEdit = setEdit;

    function admin_save() {
        if(valid()){
            return $http.post("/admin/admin/add",$scope.user)
        }
    }
    function valid(){
        if(!$scope.user._id){
            $scope.msg = "帐号不能为空!";
            return false;
        }
        if(!$scope.user.password){
            $scope.msg = "密码不能为空!";
            return false;
        }
        if($scope.user.password !== $scope.user.password_rp){
            $scope.msg = "两次密码不匹配!";
            return false;
        }
        $scope.msg = "";
        return true;
    }
    function setUser(user){
        $scope.user = user;
    }
    function setEdit(edit){
        $scope.edit = edit;
    }


});

app.controller('AdminCtrl', function ($scope, i18nService, $modal, $log, cfpLoadingBar, $http, $window) {
    //初始化查询参数
    //以下处理是为了防止树内容一次加载导致的性能问题,改为了点击加载 , 使用了非deepcopy
    $scope.users = [];

    init();
    $("#admin_dialog").dialog({
        resizable: false,
        maxHeight: 200,
        width: 300,
        autoOpen: false,
        modal: true,
        height: 300,
        buttons: {
            "保存": function () {
                var dialog =  $(this);
                var save = $("#admin_dialog").scope().admin_save();
                if(save){
                    save.then(function(resp){
                        console.log(resp);
                        if(resp.data.Success){
                            dialog.dialog("close");
                            init();
                        }
                    })
                }else{
                    $("#admin_dialog").scope().$digest();
                }
            },
            "关闭": function () {
                $(this).dialog("close");
            }
        },
        close: function (event, ui) {
        }
    });
    $scope.admin_del = admin_del;
    $scope.admin_edit = admin_edit;
    $scope.admin_add = admin_add;

    function admin_add() {
        $("#admin_dialog").scope().setUser({});
        $("#admin_dialog").scope().setEdit(false);
        $("#admin_dialog").dialog("option", "title", "新增管理员");
        $("#admin_dialog").dialog("open");
    }

    function admin_edit(idx){
        var worker = $scope.users[idx];
        $("#admin_dialog").scope().setUser(worker);
        $("#admin_dialog").scope().setEdit(true);
        $("#admin_dialog").dialog("option", "title", "修改管理员");
        $("#admin_dialog").dialog("open");
    }
    function admin_del(idx){
        var user = $scope.users[idx];
        if ($window.confirm("确定删除管理员\"" + user._id + "\"吗?") == 1) {
            $http.post("/admin/admin/del",user).then(function (ret) {
                if(ret.data.Success){
                    $scope.users.splice(idx, 1);
                }else{
                    alert('删除失败!'+ret.data.Msg);
                }
            });
        }
    }

    function init(){
        $http.get("/admin/admins").then(function (ret) {
            $scope.users = ret.data;
        });
    }

});

