# rodgoapp

Example of using Rod for testing a Go-App application.

This actually is made to demonstrate a problem with rod not being able to hijack a request and change it to have an empty body.

## Some info

As this is backend and frontend where the frontend runs as WASM, there are two components to compile. I separate part of the code by using build tags.

## Run the test

Running the test is done by:

```shell
GOARCH=wasm GOOS=js go build -o web/app.wasm && \
go build -o tester && \
./tester
```

or simpler like this:

```go
go generate && ./tester
```
