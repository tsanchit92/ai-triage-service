AI Triage Service

AI Triage Service is a Go-based backend application designed to manage and triage incidents efficiently. It utilizes Docker for containerization and MySQL for data storage, ensuring a scalable and reliable solution for incident management.

Features

Incident Management: Create, update, and retrieve incidents.

MySQL Integration: Persistent storage using MySQL database.

Dockerized Environment: Easy setup and deployment with Docker Compose.

Migration Support: Automatic database migrations during startup.

API Documentation (Swagger): Interactive API docs available.

Architecture

The application follows a microservices architecture with the following components:

API Server: Handles HTTP requests and business logic.

Database: MySQL database for storing incident data.

Migrations: Ensures the database schema is up-to-date.

Prerequisites

Ensure you have the following installed:

Docker

Docker Compose

Setup and Installation

Clone the Repository

git clone https://github.com/tsanchit92/ai-triage-service.git
cd ai-triage-service


Configure Environment Variables

Copy the example environment configuration and modify as needed:

cp conf.env.example conf.env


Edit conf.env to set your database credentials and other configurations.

Build and Start the Application

Use Docker Compose to build and start the services:

docker compose up --build


This command will:

Build the Go application.

Create and start the MySQL container.

Apply database migrations.

Start the API server.

Access the Application

Once the services are up and running, you can access the API at:

http://localhost:8080/incidents/get


Use tools like Postman
 or curl
 to interact with the API.

Swagger / API Documentation

The application includes Swagger annotations for interactive API documentation. Once the API is running:

Open the Swagger UI in your browser at:

http://localhost:8080/swagger/index.html


You can explore all available endpoints, see request/response formats, and test API calls directly from the browser.

Swagger is automatically generated from Go annotations using swaggo/swag
.

Generating Swagger Docs

If you add or modify endpoints, regenerate the Swagger documentation:

swag init -g cmd/server/main.go


This will update the docs/ directory with the latest API specification.

Database Migrations

The application uses sqlx
 for database interactions and applies migrations located in the internal/migrations directory. Ensure that the migration files are present and correctly configured.

Troubleshooting

Database Connection Issues: Ensure that the MySQL container is running and accessible. Check the logs for any errors related to database connectivity.

Missing Migration Files: If you encounter errors related to missing migration files, verify that the internal/migrations directory contains the necessary SQL files.


curls - 
GET - 
curl --location 'http://localhost:8080/incidents/get' \
--header 'Content-Type: application/json'

POST -
curl --location 'http://localhost:8080/incidents/create' \
--header 'Content-Type: application/json' \
--data '{
    "title": "Server connection errors",
    "description": "Application failing to connect to Server intermittently",
    "affected_service": "Server"
  }'
