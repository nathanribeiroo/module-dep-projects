package main

import (
	"fmt"

	"github.com/nathanribeiroo/module-dep-projects/httpclient"
)

func main() {
	res, code, err := httpclient.NewHttpClient(httpclient.OptionsHttpclient{RetryCount: 5, Timeout: 2}).
		SetUrl("https://httpbin.org/get").
		SetHeader("test", "val1").
		SetHeader("test2", "val2").
		SendGet()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
	fmt.Println(code)

}
