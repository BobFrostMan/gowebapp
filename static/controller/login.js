app.controller("LoginCtrl", loginController);
console.log("Loaded controller: 'LoginCtrl'");

function loginController($scope, $http, userManagementService){
    $scope.errVisible = "hidden";
    $scope.loginHandler = login;
    $scope.getUsers = getAllUsers;

    function login(){
        userManagementService.auth($scope, $http);
    };

    function getAllUsers(){
        userManagementService.getUsers($scope, $http);
    };
}