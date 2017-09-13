# gowebapp
Repository for Golang pet-project

# Project structure
- config/       - application settings
- static/       - location of statically served files like CSS and JS

- pet/app/executor/     - finite-state-machine executor (fsm stored as json in database)
- pet/app/model/        - database queries
- pet/app/route/        - route information and middleware
- pet/app/shared/       - packages for templates, sessions, and json

- pet/app/main.go   - application entry point

# Build project
To build project you should have GB package manager installed.
1. Install [GB package manager](https://getgb.io/), from site or using command:
```
go get github.com/constabulary/gb/...
```
Also make sure that your $GOPATH/bin directory added to $PATH variable.

2. Download all project dependencies:
```
gb vendor update -all
```

# Run project
1. Configure settings that are suitable for you. Settings located in config/ directory
2. Launch pet.* file from bin directory