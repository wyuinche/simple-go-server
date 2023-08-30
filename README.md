# simple-go-server

This project is a simple API server implementation using _gin-gonic_.

The project provides basic jwt and CRUD methods for user accounts, product information, and order lists in general shopping malls.

|Index||
|-----|-|
|[Environment](#environment)|
|[Install and Run](#install-and-run)|
|[Test](#test)|
|[Architecture](#architecture)|

## Environment

This project has been developed on __Apple M2, Ventura 13.3.1 (a)__.

```sh
~$ go version
go version go1.20 darwin/arm64

~$ sqlite3 --version
3.39.5 2022-10-14 20:58:05 554764a6e721fab307c63a4f98cd958c8428a5d9d8edfde951858d6fd02daapl
```

## Install and Run

```sh
~$ go build -o ./run

~$ ./run
```

## Test

```sh
~$ go test -v ./...
```

## Architecture

### DB

__user table__
- uid: unique id (autoincrement, primary)
- userid: general user id (must be unique)
- role: manager, user
- password: hashed password

__product table__
- pid: unique product id (autoincrement, primary)
- name: product name
- price: product price

__order table__
- oid: unique order id (autoincrement, primary)
- uid: uid who orders
- date: last update date (unix int64)

__order product table__
- oid
- pid: product ordered with oid

### Basic Rules

1. A manager can be created only by another manager.
2. A user can only view orders created by him/her.
3. Only managers can look up the entire list of orders.
4. Only managers can register, update, and delete products.
5. A user can only delete his/her account.
6. A user can only delete and update his/her orders.

### Project Architecture

- [db](./db)
    - init database [connect.go](./db/connect.go), [database.go](./db/database.go)
    - implement user, product, order crud logic
- [handler](./handler)
    - init router and load api handlers [load.go](./handler/load.go)
    - declare api methods and urls
    - implement each api handlers
    - write test codes
- [model](./model)
    - declare user, product, order struct same with those in db tables
    - these structs are used to scan columns of db
- [router](./router)
    - implement router embedding gin.Engine
- [token](./token)
    - declare claims
    - create and verify access-tokens with jwt