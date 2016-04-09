'use strict';


app.controller('TaskManagerCtrl', function ($scope, i18nService, $modal, $log,cfpLoadingBar) {
    //初始化查询参数
    $scope.query = {cat: null, rule: null, status: ''};
    //设置树参数
    //以下处理是为了防止树内容一次加载导致的性能问题,改为了点击加载 , 使用了非deepcopy
    $scope.$parent.init().then(function () {
        var base = angular.extend({}, $scope.$parent.district[0]);
        base.children = null;
        $scope.query.district = base.id;
        $scope.currentdistname = base.text;
        $scope.dist = {children: [base]};
        var rootid = $scope.$parent.district[0].id;
        var rootlength = $scope.$parent.district[0].id.length;
        $scope.extcount = 0;
        if (rootid.substring(rootid.length - 2) == '00') {
            $scope.extcount = 1;
        }
    });

    $scope.treeclick = function (item) {
        item.active = !item.active;
        $scope.query.district = item.id;
        $scope.currentdistname = item.text;
        if (!item.children) {
            var root = $scope.$parent.district[0];
            var count = (item.id.length - 6 ) / 3 + $scope.extcount;
            for (var i = 1; i <= count; i++) {
                var rootlength = root.id.length;
                var extcount = 0;
                if (root.id.substring(root.id.length - 2) == '00') {
                    extcount = 1;
                    rootlength = 4;
                }
                var childid = item.id.substr(0, rootlength + (rootlength == 4 ? i * 2 : i * 3));
                if (root.children) {
                    for (var j = 0; j < root.children.length; j++) {
                        if (root.children[j].id == childid) {
                            root = root.children[j];
                            break;
                        }
                    }
                }
            }
            if (root.children) {
                var children = [];
                for (i = 0; i < root.children.length; i++) {
                    var childitem = angular.extend({}, root.children[i]);
                    childitem.children = null;
                    children.push(childitem);
                }
                item.children = children;
            }
        }
    };
    //表格参数
    i18nService.setCurrentLang('zh-cn');
    $scope.gridOptions = {
        data: [],
        enableSorting: false,
        enableColumnMenus: false,
        enableColumnResizing: true,
        columnDefs: [
            {displayName: '任务日期', field: 'smsdate', width: 80, cellFilter: "substr:0:10"},
            {displayName: '任务类型', field: 'examname', minWidth:80 ,width:120 , enableColumnResizing: true},
            {
                displayName: '状态',
                field: 'status',
                cellTemplate: '<div class="ui-grid-cell-contents">' +
                ' <span ng-show="COL_FIELD==0">正在创建</span>' +
                ' <span ng-show="COL_FIELD==2">已完成</span> ' +
                ' <button  ng-show="COL_FIELD==1" ' +
                ' ng-click=" grid.appScope.onCellClick(row.entity)"> ' +
                ' 处理任务</button></div>',
                width: 80
            },
            {displayName: '档案号', field: 'fileno', width: 140, pinnedLeft: true},
            {displayName: '姓名', field: 'personname', width: 70, pinnedLeft: true},
            {displayName: '身份证号', field: 'idnumber', width: 150},
            {displayName: '生日', field: 'birthday', cellFilter: "date:'yyyy-MM-dd'", width: 80},
            {displayName: '电话', field: 'tel', width: 100},
            //{displayName: '内容', field: 'msg', width: 200}
            //{displayName:'完成情况',field: 'status'}
        ],
        paginationPageSizes: [20],
        paginationPageSize: 20,
        paginationCurrentPage: 1,
        useExternalPagination: true,
        //totalItems: $scope.totalItems,
        rowHeight: 30,
        enableRowSelection: true,
        enableRowHeaderSelection: false,
        multiSelect: false,
        modifierKeysToMultiSelect: false,
        noUnselect: true,
        onRegisterApi: function (gridApi) {
            $scope.gridApi = gridApi;
            $scope.gridApi.pagination.on.paginationChanged($scope, function () {
                $scope.querydata();
            });
        }
    };

    //TODO 执行点击函数
    $scope.onCellClick = function (row) {
        console.log(row);
        //打开模态窗口
        var taskurl = row['inputpage'] + '?fileNo=' + row['fileno'] + '&isNext=1&loadtaskdefault=true&taskid=' + row['id'];
        var modaldata = {url: taskurl, title: row['examname'] + ':' + row['personname'] + " " + row['msg']};
        var modalInstance = $modal.open({
            templateUrl: 'OldWindowModalContent.html',
            controller: 'ModalInstanceCtrl',
            windowClass: 'oldwindow',
            size: 'lg',
            backdrop: false,
            resolve: {
                data: function () {
                    return modaldata;
                }
            }
        });
        modalInstance.opened.then(function () {
            $('#oldwindowiframe').on('load', function () {

            });
        });
        modalInstance.result.then(function (selectedItem) {
            $scope.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };
    //初始化查询参数
    //初始化分类下拉
    TaskService.getTaskCatOption(function (data) {
        var cats = [];
        $.each(data, function (key, value) {
            if (value[0] != '') {
                cats.push({
                    id: value[0],
                    name: value[1],
                    ord: value[2]
                })
            }
        });
        $scope.cats = cats;
        $scope.$digest();
    });
    //初始化类型下拉
    TaskService.getTaskRuleOption(function (data) {
        var rules = [];
        $.each(data, function (key, value) {
            if (value[0] != '') {
                rules.push({
                    id: value[0],
                    name: value[1],
                    ord: value[3],
                    parent: value[2]
                })
            }
        });
        $scope.rules = rules;
        $scope.$digest();
    });
    //初始化类型下拉
    TaskService.getQueryParams(function (data) {
        $scope.dropdowns = data;
        $scope.dropdown = data[0]['key'];
        $scope.$digest();
    });
    //日期校验20110101-20120101的类型
    $scope.dateValidate = function (code) {
        var dateregstr = "(?:(?!0000)[0-9]{4}(?:(?:0[1-9]|1[0-2])(?:0[1-9]|1[0-9]|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[0-9]{2}(?:0[48]|[2468][048]|[13579][26])|(?:0[48]|[2468][048]|[13579][26])00)-02-29)";
        //var regexp1 = new RegExp("^(" + dateregstr + '|' + dateregstr + "-" + dateregstr + ")$");
        return "^(" + dateregstr + '|' + dateregstr + "-" + dateregstr + ")$";
    };
    $scope.$watch("query.cat", function (newval, oldval) {
        if (!inarray($scope.query.rule, $scope.rules, "id")) {
            $scope.query.rule = "";
        }
    });
    function inarray(str, array, name) {
        if (array && array.length) {
            $.each(array, function (key, value) {
                if (name && value[name] == str) {
                    return true;
                } else if (!name && value == str) {
                    return true;
                }
            });
        }
        return false;
    }

    $scope.filterRule = function (item) {
        if (!$scope.query.cat  || item.parent == $scope.query.cat ) {
            return true;
        }
        return false;
    };
    $scope.$watch("query.district", function (newval, oldval) {
        if (newval != oldval) {
            $scope.querydata();
        }
    });
    $scope.querydata = function () {
        var datestr = $scope.query.querydate;
        if (datestr) {
            var begindate = datestr;
            var enddate = datestr;
            if (datestr.indexOf("-") > 0) {
                var strs = datestr.split("-");
                begindate = strs[0];
                enddate = strs[1];
            }
            $scope.query.begindate =begindate;
            $scope.query.enddate =begindate;
        }
        for(var i = 0 ;i <$scope.dropdowns.length;i++){
            $scope.query[$scope.dropdowns[i].key]= null;
        }
        if($scope.dropdownvalue){
            $scope.query[$scope.dropdown]= $scope.dropdownvalue;
        }
        cfpLoadingBar.start();
        TaskService.queryLogsnew($scope.query, {
            start: ($scope.gridOptions.paginationCurrentPage - 1) * 20,
            limit: 20
        }, function (data) {
            cfpLoadingBar.complete()
            $scope.gridOptions.data = data.data;
            $scope.gridOptions.totalItems = data.totalSize;
            $scope.$digest();
        });
    }
});

function setFrameLoaded(obj) {
    $('#oldwindowiframe').contents().find('body').append($('<link href="/videosystem/css/oldwindow.css" rel="stylesheet"/>'));
    $('#oldwindowiframe').contents().find('.quit.img').on("click",function(){
        angular.element(obj).scope().cancel();
    });

    //$('#oldwindowiframe').contents().find('body').append($('<link href="/tasksystem/lib/beyond/css/bootstrap.min.css" rel="stylesheet"/>'));
}