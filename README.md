# gowebapp
Repository for Golang pet-project

# Project structure
- config/       - application settings and initial files
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
To run project you should have [MongoDB installed](https://www.mongodb.com/download-center#community).

1. Fill your local MongoDB database with initial data. To do this you need to execute config/init.js script on your db server.
Use commands below:
For MongoDB v2
```
mongo pet-operational --eval path_to_this_file/init_entities.js
```
For MongoDB v2.6+
```
mongo localhost:27017/pet-operational path_to_this_file/init_entities.js
```

2. Configure settings that are suitable for you and matching your DB server. Settings located in config/ directory
3. Execute initial Database
4. Launch pet.* file from bin directory