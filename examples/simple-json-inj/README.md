This example simplifies the 'simple-json' one by using http://github.com/yourheropaul/inj to auto-wire our config values into struct properties.

```
$ go run $GOPATH/src/github.com/adrianduke/configr/examples/simple-json-inj/example.go
> Email Address: my@email.com
> Subject: A Subject
> Max Retries: 5
> Retry on Fail: true
```
