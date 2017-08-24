var app = angular.module('loginModule', []);
app.controller("LoginCtrl", function($scope){
    $scope.errVisible = "hidden";

    $scope.loginHandler = function(){
        $scope.name = $scope.name_value;
        $scope.password = $scope.pass_value;
        $scope.errVisible = "visible";
        $scope.error = "You've entered " + $scope.name + " and " + $scope.password;
    };
})