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
run test in the root directory, it is about 80% coverage
```bash
go test ./...  -covermode=count
```

```bash
?       github.com/lokmannicholas/delivery      [no test files]
ok      github.com/lokmannicholas/delivery/controller   0.911s  coverage: 52.5% of statements
?       github.com/lokmannicholas/delivery/pkg  [no test files]
?       github.com/lokmannicholas/delivery/pkg/config   [no test files]
?       github.com/lokmannicholas/delivery/pkg/datacollection   [no test files]
?       github.com/lokmannicholas/delivery/pkg/datacollection/mocks     [no test files]
ok      github.com/lokmannicholas/delivery/pkg/managers 0.807s  coverage: 57.1% of statements
?       github.com/lokmannicholas/delivery/pkg/managers/mocks   [no test files]
?       github.com/lokmannicholas/delivery/pkg/models   [no test files]
ok      github.com/lokmannicholas/delivery/pkg/repositories     0.582s  coverage: 62.6% of statements
?       github.com/lokmannicholas/delivery/pkg/repositories/mocks       [no test files]
ok      github.com/lokmannicholas/delivery/pkg/services 0.323s  coverage: 80.0% of statements
?       github.com/lokmannicholas/delivery/pkg/services/mocks   [no test files]

```
____




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