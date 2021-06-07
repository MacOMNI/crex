package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/configtest"
	"github.com/coinrust/crex/exchanges"
	"github.com/gin-gonic/gin"
)

var (
	coinArray = []string{"EOS", "BTC", "DOGE", "LTC", "ETH", "FIL", "DOT", "BNB", "LTC", "LINK", "QTUM"}
)

type RobotInMessage struct {
	MessageType string `json:"msgtype"`
	Text        struct {
		Content string `json:"content,omitempty"`
	} `json:"text,omitempty"`
	MessageID         string `json:"msgId"`
	CreatedAt         int64  `json:"createAt"`
	ConversationID    string `json:"conversationId"`
	ConversationType  string `json:"conversationType"`
	ConversationTitle string `json:"conversationTitle"`
	SenderID          string `json:"senderId"`
	SenderNick        string `json:"senderNick"`
	SenderCorpID      string `json:"senderCorpId"`
	SenderStaffID     string `json:"senderStaffId"`
}

func main() {
	testConfig := configtest.LoadTestConfig("binancefutures")
	// exchange := exchanges.NewBinanceFutures(params)

	exchange := exchanges.NewExchange(exchanges.BinanceFutures,
		ApiProxyURLOption(testConfig.ProxyURL), // 使用代理
		ApiAccessKeyOption(testConfig.AccessKey),
		ApiSecretKeyOption(testConfig.SecretKey),
		ApiTestnetOption(false))

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// gin.SetMode(gin.ReleaseMode)
	router.POST("/", func(c *gin.Context) {
		//log.Printf("text: %#v", c.Copy())
		//body := c.Request.Body
		var message RobotInMessage
		if err := c.ShouldBind(&message); err != nil {
			c.Writer.Write(BuildDingMsg(err.Error()))
			return
		}
		content := strings.ToUpper(strings.TrimSpace(message.Text.Content))
		if _, err := Contain(content, coinArray); err != nil {
			c.Writer.Write(BuildDingMsg(message.Text.Content))
			return
		}
		balance, err := exchange.GetBalance("USDT")
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println("balance usdt =", balance)
		positions, err := exchange.GetPositions(content + "USDT")
		if err != nil {
			log.Fatal(err)
			c.Writer.Write(BuildDingMsg(err.Error()))
			return
		}
		if positions != nil {
			for _, v := range positions {
				// type Position struct {
				// 	Symbol    string  `json:"symbol"`     // 标
				// 	OpenPrice float64 `json:"open_price"` // 开仓价
				// 	Size      float64 `json:"size"`       // 仓位大小
				// 	AvgPrice  float64 `json:"avg_price"`  // 平均价
				// 	Profit    float64 `json:"profit"`     //浮动盈亏
				// }
				var res string
				res = "合约币种: " + v.Symbol
				res = res + "\n资产净值: " + fmt.Sprintf("%v", balance.Equity)
				res = res + "\n可用资产: " + fmt.Sprintf("%v", balance.Available)
				res = res + "\n当前价格: " + fmt.Sprintf("%v", v.MarkPrice)
				res = res + "\n持仓价格: " + fmt.Sprintf("%v", v.AvgPrice)
				if v.Size < 0 { //做空
					res = res + "\n空单数量: " + fmt.Sprintf("%v", v.Size*-1)
				} else { // 做多
					res = res + "\n多单数量: " + fmt.Sprintf("%v", v.Size*-1)
				}
				res = res + "\n合约倍数: " + fmt.Sprintf("%v", v.Leverage)
				res = res + "\n浮动盈亏: " + fmt.Sprintf("%v", v.Profit)
				c.Writer.Write(BuildDingMsg(res))
			}
		} else {
			c.Writer.Write(BuildDingMsg(message.Text.Content))
		}
		//c.String(http.StatusOK, "Hello World")
	})
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	router.Run("0.0.0.0:8080")
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}
func BuildDingMsg(msg string) []byte {
	reply := struct {
		Msgtype string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}{}
	reply.Msgtype = "text"
	reply.Text.Content = msg
	data, err := json.Marshal(reply)
	if err != nil {
		return []byte(err.Error())
	}
	return data
}
