# gowebapp
Repository for Golang pet-project

# Project structure
config/       - application settings
static/       - location of statically served files like CSS and JS
template/     - HTML templates

pet/app/controller/   - page logic organized by HTTP methods (GET, POST)
pet/app/model/        - database queries
pet/app/route/        - route information and middleware
pet/app/shared/       - packages for templates, sessions, and json

main.go   - application entry point
