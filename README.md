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

