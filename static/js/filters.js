'use strict';

angular.module('videosystem.filters', [])

    .filter('test',  function() {
        return function(text) {
            return "test";
        };
    }).filter('substr',  function() {
        return function(datetext,start,end) {
            if(end<0){
                end = datetext.length +end;
            }
            return datetext.substring(start,end);
        };
    }).filter('status',  function() {
        return function(status) {
            if (status == 2) {
                return '已完成';
            } else if (status == 1) {
                return '未完成';
            } else if (status == 0) {
                return '正在创建';
            } else {
                return '其他';
            }
        };
    });