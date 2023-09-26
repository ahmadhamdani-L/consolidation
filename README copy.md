# MCASH Financial Consolidation - Core Service

This is the core backend service of MCASH Financial Consolidation project.

## Installation

``` bash
# clone the repo
$ git clone https://gitlab.codeoffice.net/all/mcash/konsolidasi/core.git

# go into app's directory
$ cd core
```

## Configuration

| **ENV**                    | **YAML**                  | **env only** | **Req** | **Type** | **Options**                                             | **Default** | **Description**                                                                                               |
| :---                       | :---                      | :---:        | :---:   | :---     | :---                                                    | :---        | :---                                                                                                          |
| `ENV`                      |                           | √            | √       | string   | local, development, staging, production                 |             | Environment flag.                                                                                             |
| `LOG_LEVEL`                |                           | √            |         | string   | disabled, debug, info, warn, warning, err, error, fatal | disabled    | Log level.                                                                                                    |
| `HOST`                     |                           | √            |         | string   |                                                         |             | Service hostname.                                                                                             |
| `SCHEMES`                  |                           | √            |         | csv      | One or both of http and https                           | http        | List of available schemes.                                                                                    |
| `CONFIG_DIRECTORY_PATH`    |                           | √            |         | string   |                                                         |             | A directory path of the rest configuartion files. When `ENV != local`, it becomes a mandatory.                |
| `CONFIG_ENCRYPTED`         |                           | √            |         | boolean  |                                                         | false       | When `true` the config files will be decrypted on preload.                                                    |
| `STORAGE_DIRECTORY_PATH`   |                           | √            | √       | string   |                                                         |             | A directory path to store imported/exported files.                                                            |
| `JWT_KEY`                  | `jwt.key`                 |              | √       | string   |                                                         |             | Backend JWT secret key.                                                                                       |
| `DB_HOST`                  | `db.host`                 |              | √       | string   |                                                         |             | Database server hostname/ip.                                                                                  |
| `DB_PORT`                  | `db.port`                 |              |         | string   |                                                         | 5432        | Database server port.                                                                                         |
| `DB_NAME`                  | `db.name`                 |              | √       | string   |                                                         |             | Database name.                                                                                                |
| `DB_USER`                  | `db.user`                 |              | √       | string   |                                                         |             | Username to access DB server.                                                                                 |
| `DB_PASS`                  | `db.pass`                 |              | √       | string   |                                                         |             | User's password to access DB server.                                                                          |
| `DB_SSLMODE`               | `db.sslmode`              |              |         | string   | disable, allow, prefer, require, verify-ca, verify-full | disable     | Visit [PGSSLMODE](https://www.postgresql.org/docs/15/libpq-connect.html#LIBPQ-CONNECT-SSLMODE) for more info. |
| `DB_TZ`                    | `db.tz`                   |              |         | string   |                                                         | UTC         | Visit [Wiki TZ](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) for more info.                  |
| `DB_CONNECTION_LIFETIME`   | `db.connection_lifetime`  |              |         | string   | Available suffix: ns, us, ms, s, m, h                   | 0s          | Max amount of time a connection maybe reused. When zero, the connection will never get closed.                |
| `DB_MAX_OPEN_CONNECTIONS`  | `db.max_open_connections` |              |         | integer  |                                                         | 0           | Max amount of open connection to the database.                                                                |
| `DB_MAX_IDLE_CONNECTIONS`  | `db.max_idle_connections` |              |         | integer  |                                                         | 2           | Max amount of idle open connection to the database.                                                           |
| `REDIS_HOST`               | `redis.host`              |              | √       | string   |                                                         |             | Redis server hostname/ip.                                                                                     |
| `REDIS_PORT`               | `redis.port`              |              |         | string   |                                                         | 6379        | Redis server port.                                                                                            |
| `REDIS_DB`                 | `redis.db`                |              |         | integer  |                                                         | 0           | Redis database index number.                                                                                  |
| `REDIS_PASSWORD`           | `redis.password`          |              | √       | string   |                                                         |             | Redis passkey.                                                                                                |
| `REDIS_POOL_SIZE`          | `redis.pool_size`         |              |         | integer  |                                                         | 10          | Max amount of redis open connection.                                                                          |

## Build & Run

For local environment write down the necessary [environment variables](#configuration) in a file, name the file as `.env.local`, and put it in the poject's root directory. Then execute the following command:

``` bash
# Run application
$ ENV=local go run main.go
```

## Documentation

Install environment:

``` bash
# get swagger package 
$ go get github.com/swaggo/swag

# move to swagger dir
$ cd $GOPATH/src/github.com/swaggo/swag

# install swagger cmd 
$ go install cmd/swag
```

Generate documentation:

``` bash
# generate swagger doc
$ swag init --propertyStrategy snakecase
```

to see the results, run app and access {{url}}/swagger/index.html

## Description

This project built in clean architecture that contains :

1. Http
2. Factory
3. Middleware
4. Handler
5. Binder
6. Validation
7. Service
8. Repository
9. Model
10. Database

This project have some default function :

- Context
- Validator
- Transaction
- Pagination & Sort
- Filter
- Env
- Response
- Redis
- Elasticsearch
- Log

This project have some default endpoint :

- Auth
  - Login
  - Register
- Sample
  - Get (with pagination, sort, & filter)
  - GetByID
  - Create (with transaction)
  - Update (with transaction)
  - Delete

## Author

CodeID Backend Team
