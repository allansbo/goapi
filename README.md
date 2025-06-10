# Location API

This project aimed to train concepts of the Golang language, project structure, documentation using swagger and the creation of the API itself.

The project manages fleet vehicles overtime, through of a simple CRUD operation.

#### Receive data about:
- The location (latitude and longitude)
- If the vehicle is moving, stopped or offline
- The speed recorded at the time of collection

## Technologies

- Golang 1.24.4
- Fiber framework
- MongoDB
- Docker

## How to run

1. Rename `.env.example` into `.env`
2. Set up the environment variables at `.env` file
   1. If change the default values for `DB_PORT` and `APP_PORT`, change the value from `ports` section at `docker-compose.yml`
3. Run the `docker-compose.yml` or `makefile`
   1. ```shell
      docker compose up -d --build
      ```
   2. ```shell
      make up
      ```

## The swagger

When the project is running, the default route to swagger will be:

- http://localhost:8080/docs/index.html

## The Health Check

The API implements routes to check the health from the application, both provided by the Fiber framework. [View about](https://docs.gofiber.io/api/middleware/healthcheck/).

- http://localhost:8080/livez
- http://localhost:8080/readyz

Both returns "OK" if the application is healthy