## Organono Microservice
A simple lightweight service for managing companies

#### Requirements

- Golang v1.17
- Docker
- Postman
- Install godo: `go get -u gopkg.in/godo.v2/cmd/godo`

#### Run Application Locally

- Clone the repo `git clone https://github.com/vonmutinda/organono.git && cd organono`
- Create `.env` and `.env.test` using `.env.sample` and `.env.test.sample` files respectively
- Start PostgresSQL server `make up`
- Install goose `go install github.com/pressly/goose/cmd/goose@latest`
- Run Migrations `make migrate`

```shell
➜  organono git:(main) ✗ make migrate
goose -dir 'app/db/migrations' postgres "postgres://admin:password@localhost:5433/organono_backend?sslmode=disable" up
2022/08/02 18:39:12 OK    20220801141410_initial.sql
2022/08/02 18:39:12 goose: no migrations to run. current version: 20220801141410
```

- Install deps `go mod tidy`
- Populate countries `make populate_countries`

```shell
➜  organono git:(main) ✗ make populate_countries
go run scripts/populate_countries.go -e .env
Inserting Afghanistan
Inserting Aland Islands
Inserting Albania
Inserting Algeria
Inserting American Samoa
Inserting Andorra
Inserting Angola
Inserting Anguilla
Inserting Antarctica
Inserting Antigua and Barbuda
...
```

- Start Server `make server`
- Run tests `make test` or `make test-lite`

```shell
➜  organono git:(main) ✗ make test-lite
godo test-lite -- -e .env.test
?       github.com/vonmutinda/organono/Gododir  [no test files]
?       github.com/vonmutinda/organono/app/db   [no test files]
?       github.com/vonmutinda/organono/app/entities     [no test files]
?       github.com/vonmutinda/organono/app/forms        [no test files]
?       github.com/vonmutinda/organono/app/logger       [no test files]
?       github.com/vonmutinda/organono/app/providers    [no test files]
ok      github.com/vonmutinda/organono/app/repos        (cached)
ok      github.com/vonmutinda/organono/app/services     (cached)
?       github.com/vonmutinda/organono/app/utils        [no test files]
ok      github.com/vonmutinda/organono/app/web/api/companies    5.264s
?       github.com/vonmutinda/organono/app/web/api/sessions     [no test files]
?       github.com/vonmutinda/organono/app/web/auth     [no test files]
?       github.com/vonmutinda/organono/app/web/ctxhelper        [no test files]
?       github.com/vonmutinda/organono/app/web/middleware       [no test files]
?       github.com/vonmutinda/organono/app/web/router   [no test files]
?       github.com/vonmutinda/organono/app/web/webutils [no test files]
?       github.com/vonmutinda/organono/cmd      [no test files]
?       github.com/vonmutinda/organono/scripts  [no test files]
test-lite 6972ms
```

#### Endpoints

Download Postman Collection for a better testing experience [Click to Download](https://github.com/vonmutinda/organono/blob/main/Organono.postman_collection.json)

##### 1. Auth

No authentication is required when interacting with CREATE and DELETE APIs from Cyprus.

Provide correct `username` and `password` to acquire an auth token that you'll need to access resources from protected routes.

HTTP POST `localhost:3000/v1/auth`

Request Body

```shell
{
  "username": "xm",
  "password": "password"
}
```

Response Body

```shell
{
  "success": true
}
```

Authentication Token `X-ORGANONO-Token`

<img src="assets/login-to-organono.png"  alt="login page">

##### 2. Create Company

Set `X-ORGANONO-Token` in the Request Header

<img src="assets/set-x-organono-token.png"  alt="login page">

HTTP POST `localhost:3000/v1/companies`

Request Body

```shell
{
    "name": "Trading Point",
    "code": "TP",
    "phone": "35790034567",
    "website": "https://xm.com",
    "country": "cyprus"
}
```

Response Body

```shell
{
    "id": 1,
    "name": "Trading Point",
    "code": "TP",
    "country": "Cyprus",
    "website": "https://xm.com",
    "phone": "+35790034567",
    "phone_number": {
        "country_code": "+357",
        "number": "90034568"
    },
    "operation_status": "",
    "created_at": "2022-08-02T21:33:33.7311+03:00",
    "updated_at": "2022-08-02T21:33:33.7311+03:00"
}
```

##### 3. Update Company

HTTP PUT `localhost:3000/v1/companies/{id}`

Reqest Body

```shell
{
    "name": "Trading CFDs",
    "code": "TCF"
}
```

Response Body

```shell
{
    "id": 1,
    "name": "Trading CFDs",
    "code": "TCF",
    "country": "Cyprus",
    "website": "https://xm.com",
    "phone": "",
    "phone_number": {
        "country_code": "+357",
        "number": "90034567"
    },
    "operation_status": "active",
    "created_at": "2022-08-02T19:23:23.443329Z",
    "updated_at": "2022-08-02T21:35:32.440811+03:00"
}
```

##### 4. Delete Company

HTTP DELETE `localhost:3000/v1/companies/{id}`

Response Body

```shell
{
    "id": 1,
    "name": "Trading Point",
    "code": "BP",
    "country": "Cyprus",
    "website": "https://xm.com",
    "phone": "",
    "phone_number": {
        "country_code": "+357",
        "number": "90034567"
    },
    "operation_status": "active",
    "created_at": "2022-08-02T21:33:33.7311Z",
    "updated_at": "2022-08-02T21:33:33.7311Z"
}
```

##### 5. List Companies
HTTP GET `localhost:3000/v1/companies`

Response Body
```shell
{
    "companies": [
        {
            "id": 2,
            "name": "Google",
            "code": "GGL",
            "country": "United States",
            "website": "https://google.com",
              ...
        },
        {
            "id": 3,
            "name": "Safaricom",
            "code": "SAF",
            "country": "Kenya",
            "website": "https://safaricom.com",
              ...
        }
    ],
    "pagination": {
        "count": 3,
        "next_page": null,
        "num_pages": 1,
        "page": 1,
        "per": 20,
        "prev_page": null
    }
}
```

##### 6. Logout

Deactivate active session

HTTP DELETE `localhost:3000/v1/auth`

Response Body

```shell
{
  "success": true
}
```
