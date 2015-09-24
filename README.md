# configr

Configr provides an abstraction above configuration sources, allowing you to use a single interface to expect and get all your configuration values.

**Features:**
- Single interface for configuration values: Simple API (Get(), String(), Bool()...)
- Extendable config sources: Load config from a file, database or any source you can get data from
- Multiple source support: Add as many sources as you can manage, FILO merge strategy employed (first source added has highest priority)
- Value validation support: Any matching key from every source is validated by your custom validators
- Required keys support: Ensure keys exist after parsing, otherwise error out
- Blank config generator: Register as many keys as you need and use the blank config generator
- Custom blank config encoder support: Implement an encoder for any data format and have a blank config generated in it
- Comes pre-baked with JSON and TOML file support

Built for a project at [HomeMade Digital](http://homemadedigital.com/), configrs primary goals were to eliminate user error when deploying projects with heavy configuration needs. The inclusion of required key support, validators, value descriptions and blank config generator allowed us to reduce pain for seperated client ops teams when deploying our apps. Our secondary goal was configurable configuration sources be it from pulling from Mongo, DynamoDB, JSON or TOML files.

