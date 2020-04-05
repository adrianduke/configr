package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrianduke/configr"
	_ "github.com/adrianduke/configr/sources/file/json"
	"github.com/davecgh/go-spew/spew"
)

func DefaultConfig() Config {
	return Config{
		Email: DefaultEmail(),
		User:  User{},
	}
}

type Config struct {
	Email Email
	User  User
}

type User struct {
	Name    string `configr:"username"`
	Address Address
}

type Address struct {
	Line1, Line2, Line3 string
	Country             string
	Postcode            string `configr:"zipcode,required"`
}

func DefaultEmail() Email {
	return Email{
		MaxRetries:  5,
		RetryOnFail: true,
	}
}

type Email struct {
	FromAddress string `configr:",required"`
	Subject     string
	MaxRetries  int
	RetryOnFail bool
}

func main() {
	// Register our expected fields and default values
	configDefaults := DefaultConfig()
	configr.RegisterFromStruct(&configDefaults, configr.ToLowerCamelCase)

	// Wont work if GOPATH contains multiple DIRs
	path := filepath.Join(os.Getenv("GOPATH"), "src/github.com/adrianduke/configr/examples/from-struct/config.json")
	f := configr.NewFile(path)
	configr.AddSource(f)

	if err := configr.Parse(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := Config{}
	configr.Unmarshal(&config)
	spew.Dump(config)
}
