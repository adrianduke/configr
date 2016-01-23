package main

import (
	"fmt"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/adrianduke/configr"
	_ "github.com/adrianduke/configr/sources/file/json"
	"github.com/yourheropaul/inj"
)

const (
	EmailFromAddressKey = "email.fromAddress"
	EmailSubjectKey     = "email.subject"
	EmailRetryOnFailKey = "email.retryOnFail"
	EmailMaxRetriesKey  = "email.maxRetries"
)

type Email struct {
	FromAddress string `inj:"email.fromAddress"`
	Subject     string `inj:"email.subject"`
	MaxRetries  int    `inj:"email.maxRetries"`
	RetryOnFail bool   `inj:"email.retryOnFail"`
}

func init() {
	configr.RequireKey(EmailFromAddressKey, "Email from address", func(v interface{}) error {
		_, err := mail.ParseAddress(v.(string))
		return err
	})
	configr.RequireKey(EmailSubjectKey, "Email subject")
	configr.RegisterKey(EmailRetryOnFailKey, "Retry sending email if it fails", false)
	configr.RegisterKey(EmailMaxRetriesKey, "How many times to retry email resending", 3)
}

func main() {
	// Wont work if GOPATH contains multiple DIRs
	path := filepath.Join(os.Getenv("GOPATH"), "src/github.com/adrianduke/configr/examples/simple-json-inj/config.json")
	configr.AddSource(configr.NewFile(path))

	if err := configr.Parse(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	email := Email{}
	inj.Provide(&email)                     // Informs inj to perform DI on given instance
	inj.AddDatasource(configr.GetConfigr()) // Provides inj with a datasource to query

	if valid, errors := inj.Assert(); !valid { // Triggers the inj DI process
		fmt.Println(errors)
		os.Exit(1)
	}

	fmt.Println("> Email Address:", email.FromAddress)
	fmt.Println("> Subject:", email.Subject)
	fmt.Println("> Max Retries:", email.MaxRetries)
	fmt.Println("> Retry on Fail:", email.RetryOnFail)
}
