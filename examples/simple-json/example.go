package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrianduke/configr"
)

const (
	EmailFromAddressKey = "email.fromAddress"
	EmailSubjectKey     = "email.subject"
	EmailRetryOnFailKey = "email.retryOnFail"
	EmailMaxRetriesKey  = "email.maxRetries"
)

type Email struct {
	FromAddress string
	Subject     string
	MaxRetries  int  // Note config.json contains string
	RetryOnFail bool // Note config.json contains string
}

func init() {
	configr.RequireKey(EmailFromAddressKey, "Email from address")
	configr.RequireKey(EmailSubjectKey, "Email subject")
	configr.RegisterKey(EmailRetryOnFailKey, "Retry sending email if it fails", false)
	configr.RegisterKey(EmailMaxRetriesKey, "How many times to retry email resending", 3)
}

func main() {
	// Wont work if GOPATH contains multiple DIRs
	path := filepath.Join(os.Getenv("GOPATH"), "src/github.com/adrianduke/configr/examples/simple-json/config.json")
	configr.AddSource(configr.NewFileSource(path))

	if err := configr.Parse(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fromAddress, err := configr.String(EmailFromAddressKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	emailSubject, err := configr.String(EmailSubjectKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	retryOnFail, err := configr.Bool(EmailRetryOnFailKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	maxRetries, err := configr.Int(EmailMaxRetriesKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	email := Email{
		FromAddress: fromAddress,
		Subject:     emailSubject,
		MaxRetries:  maxRetries,
		RetryOnFail: retryOnFail,
	}

	fmt.Printf("Gotten values:\n\t%#v\n\n", email)

	//
	// Generate blank config
	//
	f := configr.NewFileSource("")
	f.SetEncoding(configr.JSON)

	// f implements the encoder interface
	configBytes, err := configr.GenerateBlank(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Blank config:\n%s", string(configBytes))
}
