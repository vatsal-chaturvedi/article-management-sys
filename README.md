[![Build & Test](https://github.com/vatsal-chaturvedi/article-management-sys/actions/workflows/build_test.yml/badge.svg)](https://github.com/vatsal-chaturvedi/article-management-sys/actions/workflows/build_test.yml) [![codecov](https://codecov.io/gh/vatsal-chaturvedi/article-management-sys/branch/main/graph/badge.svg?token=S4Q2G7L25O)](https://codecov.io/gh/vatsal-chaturvedi/article-management-sys)
## Backend API for a Article Management System
* This is a REST HTTP API for a article management system built using Go. 
* It implements three endpoints for creating, retrieving, and listing articles. 
### Features: 
* Uses clean architecture and design patterns and is tested using unit and integration tests. The application can be run in Docker, and the repository contains a docker-compose.yml file and a start.sh bash script for setting up the relevant services and applications. 
* Uses a MySQL database, and the installation and initialization of the DB are done when `start.sh` is executed.
* Get all article endpoint uses pagination and default limit is set to 20 so that the response time is fast and you can provide header query params for key `limit` and `page` as integers to change them.
* Get endpoints uses caching middleware for caching the response for 10 seconds.
## Running the Application
* Run the following command to start the application:
```
sh start.sh
```
* This will build and start the Docker containers for the application and the database.
Once the containers are up and running, the API can be accessed at `http://localhost:8080`
* You can test the api using postman, just import the [Postman Collection](./article-management-system.postman_collection.json) into your postman app.

## Testing the Application
* GitHub actions have been used with codecov to generate the project code coverage as a badge with over 80% lines of code covered by unit tests.
* Additionally the test cases can be executed by
```bash
go test ./... --cover
```