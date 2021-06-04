package main

import (
	"log"
	"time"

	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
	"github.com/coinrust/crex/exchanges"
)

type BasicStrategy struct {
	StrategyBase
}

func (s *BasicStrategy) OnInit() error {
	return nil
}

func (s *BasicStrategy) OnTick() error {
	// currency := "BTC"
	// symbol := "BTC-PERPETUAL"

	balance, err := s.Exchange.GetBalance("USDT")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("balance usdt =", balance)
	log.Printf("balance: %#v", balance)
	positions, err := s.Exchange.GetPositions("ETHUSDT")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range positions {
		log.Printf("position: %#v", v)

	}
	//log.Printf("positions: %#v", positions)

	return nil
}

func (s *BasicStrategy) Run() error {
	// run loop
	for {
		s.OnTick()
		time.Sleep(120 * time.Second)
	}
	return nil
}

func (s *BasicStrategy) OnExit() error {
	return nil
}

func main() {
	testConfig := configtest.LoadTestConfig("binancefutures")
	// exchange := exchanges.NewBinanceFutures(params)
	
	exchange := exchanges.NewExchange(exchanges.BinanceFutures,
		ApiProxyURLOption(testConfig.ProxyURL), // 使用代理
		ApiAccessKeyOption(testConfig.AccessKey),
		ApiSecretKeyOption(testConfig.SecretKey),
		ApiTestnetOption(false))
	s := &BasicStrategy{}

	s.Setup(TradeModeLiveTrading, exchange)

	s.OnInit()
	s.Run()
	s.OnExit()
}
