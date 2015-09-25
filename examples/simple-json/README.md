```
$ go run $GOPATH/src/github.com/adrianduke/configr/examples/simple-json/example.go
Gotten values:
	main.Email{FromAddress:"my@email.com", Subject:"A Subject", MaxRetries:5, RetryOnFail:true}

Blank config:
{
	"email": {
		"fromAddress": "*** Email from address ***",
		"maxRetries": 3,
		"retryOnFail": false,
		"subject": "*** Email subject ***"
	}
}
```