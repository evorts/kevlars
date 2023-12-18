# Kevlars

A golang backend development kits to accelerate building services. 

## Background

Having to write over and over again something like logger, db connection, caching, audit log, rest client, grpc client, telemetry, config/secret, healthcheck, feature flag, and so on. 
I find it waste a lot of time, while it should be re-usable everywhere in multiple service.
There-fore come to mind to build a development kits to help me reach that goal.

## Usage

Since on of the purpose of this development kit to accelerate building backend service, below we can see how easy and fast we can scaffold new RESTful API.
```go
scaffold.
    NewApp(scaffold.WithScope("restful_api")).
    WithDatabases().
    RunRestApiUseEcho(func(app *scaffold.Application, e *echo.Echo) {
        // do something such as routing and needed process
        // all resources that instantiate above (e.g. WithDatabase) are available under `app`
        // for example: 
        // if we want to get config value:
        // app.Config().GetString("key")
        // if we want to use database connection:
        // app.DefaultDB().Exec(ctx, q, args...)
})
```
Another sample when we need to run service as GRPC Server, it would simply be done by following sample:
```go
scaffold.
    NewApp(scaffold.WithScope("restful_api")).
    WithDatabases().
    RunGrpcServer(run func(app *Application, rpcServer *grpc.Server)) {
        // do something such as register proto and needed process
        // all resources that instantiate above (e.g. WithDatabase) are available under `app`
        // for example: 
        // if we want to get config value:
        // app.Config().GetString("key")
        // if we want to use database connection:
        // app.DefaultDB().Exec(ctx, q, args...)
})
```
For more samples, can refer to directory `examples` in this repository.

## Structure

