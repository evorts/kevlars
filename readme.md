# Kevlars

A golang backend development kits to accelerate building services. 

## Background

Having to write over and over again something like logger, db connection, caching, audit log, rest client, grpc client, telemetry, config/secret, healthcheck, feature flag, and so on. 
I find it waste a lot of time, while it should be re-usable everywhere in multiple service.
There-fore come to mind to build a development kits to help me reach that goal.
During those development process in private, I find it that this might help anyone whom has the same goal as mine. Thus, I decide to make this project open source.
To anyone using this library and have feedback or would like to contribute, please don't hesitate to do so.

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

### Captcha

This package is used to generate captcha image and validate user input.
> Note: currently only support generating base64 encoded image.

To use this package, simply import it and use it like below:
```go
cb64 := captcha.NewB64(captcha.B64WithType(captcha.B64TypeDigit)).MustInit()
// generate base 64 image
// capId => is the captcha id use to validate
// capImg => base 64 encoded image
capId, _, capImg := cb64.Generate()

// to verify user input:
// userInput => should be the variable contains the value of user input
// capId => use the generated id above
isValid := cb64.Verify(capId, userInput, true)
fmt.Println(isValid)
```
