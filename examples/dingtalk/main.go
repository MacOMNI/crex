package main

import (
	"fmt"

	"github.com/coinrust/crex/configtest"
	"github.com/hugozhu/godingtalk"
)

func main() {
	testConfig := configtest.LoadTestConfig("binancefutures")
	// exchange := exchanges.NewBinanceFutures(params)
	c := godingtalk.NewDingTalkClient("", "")
	c.AccessToken = testConfig.DingTalk
	res, err := c.SendRobotMarkdownMessage(testConfig.DingTalk, "binance", "btc price is xxx")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println("res = ", res)

}
