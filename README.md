# Software Engineer Challenge (Backend)

## installation
download the go packages cache for the project

run 
```bash
go mod download
```
download the go packages to vendor folder

run 
```bash
go mod vendor
```

which is also included in `start.sh`

## Set Environment variable
In  `docker-compose.ymal`, 
set up the default mysql variable
```bash
      MYSQL_DATABASE: '<MYSQL_DATABASE>'
      MYSQL_USER: '<MYSQL_USERNAME>'
      MYSQL_PASSWORD: '<MYSQL_PASSWORD>'
      MYSQL_ROOT_PASSWORD: '<MYSQL_ROOT_PASSWORD>'
```

then set up the variable for the app
 ```bash
      MAP_API_KEY: '<GOOGLE_MAP_API_KEY>'
      MYSQL_ADDR: 'delievery-mysql:3306'
      MYSQL_USER: '<MYSQL_USERNAME>'
      MYSQL_PASSWORD: '<MYSQL_PASSWORD>'
      MYSQL_DB: '<MYSQL_DATABASE>'
```
the `MYSQL_ADDR` is the container name and the port of mysql



## Run app
Run the `start.sh` in the project root directory 
```bash
. start.sh
```