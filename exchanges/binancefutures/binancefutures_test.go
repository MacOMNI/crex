package binancefutures

import (
	"fmt"
	"testing"
	"time"

	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
)

func testExchange() Exchange {
	testConfig := configtest.LoadTestConfig("binancefutures")
	params := &Parameters{
		AccessKey: testConfig.AccessKey,
		SecretKey: testConfig.SecretKey,
		Testnet:   testConfig.Testnet,
		ProxyURL:  testConfig.ProxyURL,
	}
	ex := NewBinanceFutures(params)
	return ex
}

func TestBinanceFutures_GetTime(t *testing.T) {
	ex := testExchange()
	tm, err := ex.GetTime()
	if err != nil {
		t.Error(err)
		fmt.Println(err.Error())
		return
	}
	t.Logf("%v", tm)
}

func TestBinanceFutures_GetBalance(t *testing.T) {
	ex := testExchange()
	balance, err := ex.GetBalance("USDT")
	if err != nil {
		fmt.Println(err.Error())
		t.Error(err)
		return
	}
	fmt.Println("USDT = ", balance)
	t.Logf("%#v", balance)
}
func TestBinanceFutures_GetPosition(t *testing.T) {
	ex := testExchange()
	position, err := ex.GetPositions("ETHUSDT")
	if err != nil {
		fmt.Println(err.Error())
		t.Error(err)
		return
	}
	fmt.Println("position = ", position)
	t.Logf("%#v", position)
}
func TestBinanceFutures_GetOrderBook(t *testing.T) {
	ex := testExchange()
	ob, err := ex.GetOrderBook("ETHUSDT", 10)
	if err != nil {
		t.Error(err)
		fmt.Println(err.Error())

		return
	}
	fmt.Println("OrderBook = ", ob)

	t.Logf("%#v", ob)
}

func TestBinanceFutures_GetRecords(t *testing.T) {
	ex := testExchange()
	now := time.Now()
	start := now.Add(-300 * time.Minute)
	end := now
	records, err := ex.GetRecords("ETHUSDT",
		PERIOD_1MIN, start.Unix(), end.Unix(), 10)
	if err != nil {
		t.Error(err)
		fmt.Println(err.Error())

		return
	}

	for _, v := range records {
		fmt.Println("Records value = ", v)

		t.Logf("Timestamp: %v %#v", v.Timestamp, v)
	}
}

func TestBinanceFutures_GetOpenOrders(t *testing.T) {
	ex := testExchange()
	orders, err := ex.GetOpenOrders("BTCUSDT")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", orders)
}

func TestBinanceFutures_ChangeLeverage(t *testing.T) {
	ex := testExchange()
	binance := ex.(*BinanceFutures)
	err := binance.ChangeLeverage("BTCUSDT", 50)
	if err != nil {
		t.Error(err)
		return
	}
}
