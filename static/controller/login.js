var app = angular.module('loginModule', []);
app.controller("LoginCtrl", loginController);

function loginController($scope, $http){
    $scope.errVisible = "hidden";
    $scope.loginHandler = login;

    function login(){
            $scope.name = $scope.name_value;
            $scope.password = $scope.pass_value;

            $http({
                method : "GET",
                url : "/users"
            }).then(
                function mySuccess(response) {
                    $scope.data = response.data;
                }, function myError(response) {
                    $scope.status = response.statusText;
                }
            );

            $scope.errVisible = "visible";
            $scope.error = "You've entered " + $scope.name + " and " + $scope.password;
    };
}