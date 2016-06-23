package main

import (
	"fmt"
	"net/mail"
	"os"

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
	MaxRetries  int  // Note env vars are string
	RetryOnFail bool // Note env vars are string
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
	configr.AddSource(configr.NewEnvVars("configr"))

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

	fmt.Printf("Values:\n\t%#v\n\n", email)

	// Un-polute environment
	os.Unsetenv("CONFIGR_EMAIL_FROMADDRESS")
	os.Unsetenv("CONFIGR_EMAIL_SUBJECT")
	os.Unsetenv("CONFIGR_EMAIL_RETRYONFAIL")
	os.Unsetenv("CONFIGR_EMAIL_MAXRETRIES")
}
