app.service("userManagementService", userService);

function userService(){

    this.getUsers = function ($scope, $http){
        console.log("Getting list of existing users")
        $http({
            method : "GET",
            url : "/users"
        }).then(
            function onSuccess(response) {
                $scope.data = response.data;
                console.log("Server returned response:");
                console.dir(response.data);
            }, function onError(error, status) {
                $scope.status = { message: error, status: status};
                console.error("Server responded with error: " +  $scope.status);
            }
        );
    }

    this.loginStub = function ($scope){
        console.log("Invoke login stub")
        $scope.name = $scope.login;
        $scope.password = $scope.pass;
        $scope.errVisible = "visible";
        $scope.error = "You've entered " + $scope.name + " and " + $scope.password;
    }

    this.auth = function($scope, $http){
        console.log("Attempt to authorize as '" + $scope.login + "'...")
        //Use jquery for serializing json data
        var data =
            $.param({
                login: $scope.login,
                pass: $scope.pass
            });

            var config = {
                headers : {
                    'Content-Type': 'application/x-www-form-urlencoded;charset=utf-8;'
                }
            }

        $http.post("/auth", data, config).then(
            function(response){
                $scope.data = response.data;
                logSuccess(response);
                window.location.href = "dashboard.html";
            },
            function(response){
                logError(response);
                showError($scope, response);
            }
        );
    }

    function logSuccess(response) {
        console.log("Server returned response:");
        console.dir(response.data);
    }

    function logError(response) {
        var err = { message: response.data.message, status: response.status};
        console.error("Server responded with error:")
        console.dir(err);
    }

    function showError($scope, response){
        console.log(status);
        if (response.status === 403){
            $scope.error = "Invalid login or password";
        }

        if (response.status >= 500){
            $scope.error = "Service is currently unavailable. Please, retry attempt in a few minutes";
        }
        console.log($scope.error);
        $scope.showError = true;
    }

}




