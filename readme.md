# How to run the application
- First install docker and docker-compose ( https://docs.docker.com/compose/install/ )
- Install GO (https://go.dev/doc/install)
- Clone this repository `git clone https://github.com/swavan/api-gateway.git`
- Go inside the cloned repo folder `cd api-gateway`
- Start postgres database using docker-compose `docker compose up -d`
- Create a environment file `.local.env`, Put below the values on it
```
CONFIG_FILE_PATH=../../data
CONFIG_FILE_NAME=config
CONFIG_FILE_EXTENSION=yaml
DATABASE_URL="user=postgres host=localhost password=postgres dbname=postgres sslmode=disable"
AUTH_CONFIG_FILE_NAME=auth
AUTH_CONFIG_FILE_LOCATION=../../data
AUTH_CONFIG_FILE_EXTENSION=yaml
SECRET_SALT=N1PCdw3M2B1TfJhoaY2mL736p2vCUc47
APP_NAME=swavan-api-gateway
```
- Load the environment variable in the terminal `export $(xargs < .local.env)`
- Run alert service `go run cmd/alert/alert.go`
- Open another terminal and load environment variable `export $(xargs < .local.env)`
- Run API gateway `go run cmd/gateway/main.go`
- Now visit `http://localhost:8000/health`

~~~~You manage to run application~~~~~


