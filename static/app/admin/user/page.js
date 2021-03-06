'use strict';
app.controller('UserPageCtrl', function ($scope, i18nService, $modal, $log, cfpLoadingBar, $http, $window) {
    $scope.user = {};
    $scope.msg = "";
    $scope.user_save = user_save;
    $scope.setUser  = setUser ;
    $scope.setEdit = setEdit;

    function user_save() {
        if(valid()){
            return $http.post("/admin/user/add",$scope.user)
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

app.controller('UserCtrl', function ($scope, i18nService, $modal, $log, cfpLoadingBar, $http, $window) {
    //初始化查询参数
    //以下处理是为了防止树内容一次加载导致的性能问题,改为了点击加载 , 使用了非deepcopy
    $scope.users = [];

    init();
    $("#user_dialog").dialog({
        resizable: false,
        maxHeight: 200,
        width: 300,
        autoOpen: false,
        modal: true,
        height: 300,
        buttons: {
            "关闭": function () {
                $(this).dialog("close");
            }
        },
        close: function (event, ui) {
        }
    });
    $scope.user_del = user_del;
    $scope.user_edit = user_edit;
    $scope.user_add = user_add;

    function user_add() {
        $("#user_dialog").scope().setUser({});
        $("#user_dialog").scope().setEdit(false);
        $("#user_dialog").dialog("option", "title", "新增工作人员");
        $("#user_dialog").dialog("open");
    }

    function user_edit(idx){
        var user = $scope.users[idx];
        console.log(user);
        $("#user_dialog").scope().setUser(user);
        $("#user_dialog").scope().setEdit(true);
        $("#user_dialog").dialog("option", "title", "修改工作人员");
        $("#user_dialog").dialog("open");
    }
    function user_del(idx){
        var worker = $scope.users[idx];
        if ($window.confirm("确定删除工作人员\"" + worker._id + "\"吗?") == 1) {
            $http.post("/admin/user/del",worker).then(function (ret) {
                if(ret.data.Success){
                    $scope.users.splice(idx, 1);
                }else{
                    alert('删除失败!'+ret.data.Msg);
                }
            });
        }
    }

    function init(){
        $http.get("/admin/users").then(function (ret) {
            $scope.users = ret.data;
        });
    }

});

