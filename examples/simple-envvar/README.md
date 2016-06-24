```
$ source $GOPATH/src/github.com/adrianduke/configr/examples/simple-envvar/envvars.sh
$ go run $GOPATH/src/github.com/adrianduke/configr/examples/simple-envvar/example.go
	Values:
		main.Email{FromAddress:"test@testing.com", Subject:"My email subject", MaxRetries:5, RetryOnFail:false}
```

Alternatively:
```
$ CONFIGR_EMAIL_FROMADDRESS=test@testing.com CONFIGR_EMAIL_SUBJECT='Test subject' go run $GOPATH/src/github.com/adrianduke/configr/examples/simple-envvar/example.go
	Values:
		main.Email{FromAddress:"test@testing.com", Subject:"Test subject", MaxRetries:5, RetryOnFail:false}

```
