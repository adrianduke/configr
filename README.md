# configr

[![Build Status](https://travis-ci.org/adrianduke/configr.svg?branch=master)](https://travis-ci.org/adrianduke/configr) [![Coverage Status](https://coveralls.io/repos/adrianduke/configr/badge.svg?branch=master&service=github)](https://coveralls.io/github/adrianduke/configr?branch=master) [![GoDoc](https://godoc.org/github.com/adrianduke/configr?status.svg)](https://godoc.org/github.com/adrianduke/configr)

Configr provides an abstraction above configuration sources, allowing you to use a single interface to expect and get all your configuration values.

**Features:**
- **Single interface for configuration values:** Simple API (Get(), String(), Bool()...)
- **Extendable config sources:** Load config from a file, database, environmental variables or any source you can get data from
- **Multiple source support:** Add as many sources as you can manage, FILO merge strategy employed (first source added has highest priority)
- **Nested Key Support:** `production.payment_gateway.public_key` `production.payment_gateway.private_key`
- **Value validation support:** Any matching key from every source is validated by your custom validators
- **Required keys support:** Ensure keys exist after parsing, otherwise error out
- **Blank config generator:** Register as many keys as you need and use the blank config generator
- **Custom blank config encoder support:** Implement an encoder for any data format and have a blank config generated in it
- **Type conversion support:** Your config has string "5" but you want an int 5? No problem
- **Comes pre-baked with JSON, TOML file support and Environmental Variables**
- **Satisfies github.com/yourheropaul/inj:Datasource:** Allows you to bypass the manual wiring of config values to struct properties (see below)

Built for a project at [HomeMade Digital](http://homemadedigital.com/), configrs primary goal was to eliminate user error when deploying projects with heavy configuration needs. The inclusion of required key support, value validators, descriptions and blank config generator allowed us to reduce pain for seperated client ops teams when deploying our apps. Our secondary goal was flexible configuration sources be it pulling from Mongo Document, DynamoDB Table, JSON or TOML files.

## Example

### Simple JSON File

Pre register some keys to expect from your configuration sources (typically via init() for small projects, or via your own object initiation system for larger projects):
```go
	configr.RequireKey("email.fromAddress", "Email from address")
	configr.RequireKey("email.subject", "Email subject")
	configr.RegisterKey("email.retryOnFail", "Retry sending email if it fails", false)
	configr.RegisterKey("email.maxRetries", "How many times to retry email resending", 3)
```

Create some configuration:
```json
{
	"email": {
		"fromAddress": "my@email.com",
		"subject": "A Subject",
		"retryOnFail": true,
		"maxRetries": 5
	}
}
```

Add a source:
```go
	configr.AddSource(configr.NewFile("/tmp/config.json"))
```

Parse your config:
```go
	if err := configr.Parse(); err != nil {
		...
	}
```

And use at your own leisure:
```go
	fromAddress, err := configr.String("email.fromAddress")
	if err != nil {
		...
	}
```

### Inj Datasource

Continuing from the simple JSON example above, you can use http://github.com/yourheropaul/inj to auto-wire in your configuration values, bypassing much of the typical config wiring boilerplate:

Pre register keys:
```go
	configr.RequireKey("email.fromAddress", "Email from address")
	configr.RequireKey("email.subject", "Email subject")
	configr.RegisterKey("email.retryOnFail", "Retry sending email if it fails", false)
	configr.RegisterKey("email.maxRetries", "How many times to retry email resending", 3)
```

Add the relevant inj struct tags with their corresponding key paths:
```go
type Email struct {
	FromAddress string `inj:"email.fromAddress"`
	Subject     string `inj:"email.subject"`
	MaxRetries  int    `inj:"email.maxRetries"`
	RetryOnFail bool   `inj:"email.retryOnFail"`
}
```

Add and setup your source (assume we're using the same config json as above):
```go
	configr.AddSource(configr.NewFile("/tmp/config.json"))
```

Parse your config:
```go
	if err := configr.Parse(); err != nil {
		...
	}
```

Setup inj with configr as its Datasource and commence the magic:
```go
	email := Email{}
	inj.Provide(&email) // Informs inj to perform DI on given instance
	inj.AddDatasource(configr.GetConfigr()) // Provides inj with a datasource to query

	if valid, errors := inj.Assert(); !valid { // Triggers the inj DI process
		...
	}
```

Marvel at the ease of auto-wiring:
```go
	fmt.Println("> Email Address:", email.FromAddress)
	fmt.Println("> Subject:", email.Subject)
	fmt.Println("> Max Retries:", email.MaxRetries)
	fmt.Println("> Retry on Fail:", email.RetryOnFail)
```

```
> Email Address: my@email.com
> Subject: A Subject
> Max Retries: 5
> Retry on Fail: true
```

More examples can be found in the `examples/` dir.

## Changes

**v0.4.0**

- Added new method `KeysToUnmarshal` to the `Source` interface, allows configr to tell your source what keys to expect, it also passes a key splitter func along so you can deconstruct nested keys to do as you please. See `./env_vars.go` for an example. Expected to be used in instances where a source doesn't have scan like functionality and needs to know the keys to search for in advance when it unmarshals.
- API Change: Added a new method `KeysToUnmarshal` to the `Source` interface, you'll need to add the method to any existing sources you have created, but it doesn't have to do anything. See `./file.go` for an example.

**v0.3.0**

- File source now supports registering encoders/decoders at a distance, check out the json and toml packages for examples
- API Change: `NewFileSource()` -> `NewFile()`

**v0.2.0**

- Add support for inj datasource

## TODO:
- Concurrent safety, particularly in multi `Parse()`'ing systems and when adding sources (will allow for hot reloads)
- ~~FileSource needs to be refactored to reduce dependency needs, something similar to sql package with a central register and blank importing the flavour you need~~
- More available sources, Env vars, Flags... etc
- Decide wether or not to ditch errors on the key getter methods (String, Get, Bool...). Alternative solution is to provide a 'Errored() bool' and 'Errors() []error or chan error' methods to Config interface.
	Arguments for:
		- Simpler interface when all you want is values
	Arguments against:
		- Error swallowing, decoupling of cause and effect (try to fetch key that cannot be converted to type ("aaa" -> Int()), user never checks configr for errors, system starts behaving weirdly)
		- Internal error managing will get funky in a concurrent environment, would have to use an error channel to pump the errors into, wouldn't be able to guarentee ordering or sacrafice performance for co-ordination
- Wrap validation errors
- Provide all primary types as getter methods
- ~~Add 'Keys' method to Source interface to accept keys and key name splitting func as parameters, provides keys for lookup for Sources that don't have 'scan' style interfaces, and potential performance improvements~~
