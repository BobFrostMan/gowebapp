var app = angular.module('app', []);

app.directive("header", header);
app.directive("footer", footer);

function header(){
    return {
        templateUrl: 'header.html'
    };
}

function footer(){
    return {
        templateUrl: 'footer.html'
    };
}