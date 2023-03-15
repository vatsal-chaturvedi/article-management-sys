[![Build & Test](https://github.com/vatsal-chaturvedi/article-management-sys/actions/workflows/build_test.yml/badge.svg)](https://github.com/vatsal-chaturvedi/article-management-sys/actions/workflows/build_test.yml) [![codecov](https://codecov.io/gh/vatsal-chaturvedi/article-management-sys/branch/main/graph/badge.svg?token=S4Q2G7L25O)](https://codecov.io/gh/vatsal-chaturvedi/article-management-sys)
# Backend API for a Article Management System
This is a RESTful HTTP API for a simple blog system built in Go. The API implements three endpoints for creating, retrieving, and listing articles. 
The solution uses clean architecture and design patterns and is tested using unit and integration tests. The application can be run in Docker, and the repository contains a docker-compose.yml file and a start.sh bash script for setting up the relevant services and applications. 
The solution uses a MySQL database, and the installation and initialization of the DB are done in start.sh.

## Running the Application
* Run the following command to start the application:
```
./start.sh
```
This will build and start the Docker containers for the application and the database.
Once the containers are up and running, the API can be accessed at `http://localhost:8080.`