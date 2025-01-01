# Web API - Pencatatan Daya Listrik Rumah Tangga

This web API provides services for recording electrical power in the house, what devices are there and how much power. This API was built using Go and PostgreSQL as the database.

## Key Features
- Store device data and electrical power
- Displays device data
- Provides an endpoint to search for device data by ID
- Add, update and delete device data

## Technologies Used
- **Go (Golang)** - The programming language used for the backend API.
- **PostgreSQL** - Relational database used to store energy data.
- **github/lib/pq** - PostgreSQL Driver for Go. 

## Prerequisites
Before running the API, make sure you have installed the following:
- [Go](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/) or use [Docker](https://docs.docker.com/get-started/get-docker/) to running PosgreSQL service.

## Installation

**Follow these steps to set up this project locally:**

   ```bash
   git clone https://github.com/kalislami/daya-listrik-api.git
   cd daya-listrik-api
   go mod tidy
   ```

## Directory Structure
#### /daya-listrik-api
│
├── /cmd/server/main.go => Contains the main code to run the HTTP server.      
│    
├── /internal   
│   ├── /handlers   
│   │   └── energy_records.go => Handles the logic for HTTP requests.   
│   ├── /models   
│   │   └── energy_record.go => Contains data structures and types used in the application.   
│   ├── /repository   
│   │   └── energy_record_repository.go => Contains code to interact with the database.   
│   └── /db/postgres.go => Contains the configuration for the PostgreSQL database connection.      
│       
├── /tests   
│   ├── api_test.go => Contains unit test code for the API.   
│   ├── api_benchmark_test.go => Contains code to benchmark the API.   
│   ├── mock.go => Mocks the repository for unit tests.   
├── /migrations => Scripts for database migration if needed.   
├── go.mod   
└── go.sum   

## Running Unit Test and Benchmark
**Run the command below, and it will display unit test and benchmark result:**

   ```bash
   #run all unit-test and benchmark
   go test -v -bench . ./tests

   #run all unit-test only
   go test -v ./tests

   #run all benchmark only
   go test -v -bench . ./tests -run=^$

   #run spesific unit-test
   go test -v ./tests -run=unit_test_func_name

   #run spesific benchmark
   go test -v -bench=benchmark_func_name ./tests -run=^$
   ```

## API Endpoints
#### The list of API endpoints can be checked in the Postman collection in this repository.

## Running the Application

#### 1. Manually

- **Make sure there is a service connection to PosgreSQL.**
- **Create .env file, an example is in the env.example file in this respository**

- **Run the command below:**

```bash
go run cmd/server/main.go
   ```
- **The API can be accessed at localhost:8080.**

#### 2. Using Docker

- **Ensure Docker is installed.**
- **Run the following commands:**

```bash
docker compose up -d
   ```
- **The API can be accessed at localhost:8080.**

## Licence
#### This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing
#### Please fork the repo and submit pull requests for contributions.

## Contact
#### For questions, contact me at [email](mailto:kamalgoritm@gmail.com).
