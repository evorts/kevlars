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

## Package Structure and Usage

### Audit
This package is used for audit log.
> Note: currently only support postgres and mysql

Usage as follows:
```go
dbm := db.New(
	db.DriverPostgreSQL, 
	"host=localhost port=5432 user=db_user password=b4Zd3x6aLRM6mKs2S3 dbname=db_name sslmode=disable", 
	db.WithMaxOpenConnection(30), 
)
al := audit.New(dbm).MustInit()
err := al.Add(context.Background(), audit.Record{
	Action: "do_something",
	CreatedByName: "system",
	BeforeChanged: map[string]interface{}{
		"name": "kevlars",
    },
	AfterChanged: map[string]interface{}{
		"name": "shield",
	}
})
fmt.Println(err)
```

### In Memory

This package is used for In Memory data.
> Currently, supported valkey and redis as provider. Will consider add another provider when deemed necessary.

To utilize this package, simply:
```go
ctx := context.Background()
// instantiate with must connect to ensure connection is established else panic
c := inmemory.NewRedis(address).MustConnect(ctx)

// get value user hash
uh := c.GetString(ctx, "user_hash")
fmt.Println(uh)

// delete user hash
err := c.Del(ctx, "user_hash")
fmt.Println(err)
```

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

### DB

This package is used to connect to database provider
> Database supported: `PostgreSQL`, `MySQL` and `SQL Server`
> 
> Meanwhile, for SQL Server, it's not fully tested yet -- use at your own risk

Usage:
```go
dbm := db.New(
	db.DriverPostgreSQL, 
	"host=localhost port=5432 user=db_user password=b4Zd3x6aLRM6mKs2S3 dbname=db_name sslmode=disable", 
	db.WithMaxOpenConnection(30), 
)
rows, err := dbm.Query(context.Background(), "SELECT * FROM users")
if err != nil {
	return
}
// ... continue your implementation here
```

### FFlag

This package is used to manage feature flag. Use this on any block of code that you want to be able to turn on/off.
Sample usage:
```go
dbm := db.New(
    db.DriverPostgreSQL,
    "host=localhost port=5432 user=db_user password=b4Zd3x6aLRM6mKs2S3 dbname=db_name sslmode=disable",
    db.WithMaxOpenConnection(30),
)
ff := fflag.New(fflag.WithDB(dbm)).MustInit()
ctx := context.Background()
ff.Add(ctx, fflag.Record{
	Feature:"feature_name", 
	Enabled: true, 
	LastChangedBy:"who_is_me",
})
ff.ExecWhenEnabled(ctx, "feature_name", func() {
    // this will be executed when feature_name are enabled	
})

```

### Messaging

This package is used to send messaging to selected provider. 
> Currently only support sending message to `telegram`. Will add other provider as things progress.

To use this package, simply import it and use it like below:
```go
tgSender := messaging.NewTelegramSender(
        messaging.TelegramWithTarget(app.Config().GetString("messaging.telegram.default_target")),
        messaging.TelegramWithToken(app.Config().GetString("messaging.telegram.token")),
        messaging.TelegramWithSanitizer(messaging.StandardSanitizer()), 
    ).MustInit()
err := tgSender.SendMessage(`hello!`)
fmt.Println(err)
```