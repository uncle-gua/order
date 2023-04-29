package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/uncle-gua/gobinance/futures"
	"github.com/uncle-gua/log"
)

var config struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

var (
	help             bool
	cfgFile          string
	symbol           string
	side             string
	positionSide     string
	_type            string
	reduceOnly       bool
	quantity         string
	price            string
	newClientOrderId string
	stopPrice        string
	closePosition    bool
	activationPrice  string
	timeInForce      string
	callbackRate     string
	workingType      string
	priceProtect     bool
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&cfgFile, "c", "", "config file")
	flag.StringVar(&symbol, "symbol", "", "symbol")
	flag.StringVar(&side, "s", "", "side")
	flag.StringVar(&positionSide, "ps", "", "position side")
	flag.StringVar(&_type, "t", "MARKET", "type")
	flag.BoolVar(&reduceOnly, "ro", false, "reduce only")
	flag.StringVar(&quantity, "q", "", "quantity")
	flag.StringVar(&price, "p", "", "price")
	flag.StringVar(&newClientOrderId, "co", "", "new client order id")
	flag.StringVar(&stopPrice, "sp", "", "stop price")
	flag.BoolVar(&closePosition, "cp", false, "close position")
	flag.StringVar(&activationPrice, "ap", "", "activation price")
	flag.StringVar(&timeInForce, "tf", "", "time in force")
	flag.StringVar(&callbackRate, "cr", "", "callback rate")
	flag.StringVar(&workingType, "wt", "", "working type")
	flag.BoolVar(&priceProtect, "pp", false, "price protect")
}

func main() {
	log.SetLevel(log.TraceLevel)

	flag.Parse()

	if help || cfgFile == "" || symbol == "" || side == "" || positionSide == "" || _type == "" {
		flag.Usage()
		log.Exit(0)
	}

	if strings.LastIndex(cfgFile, ".") == -1 {
		cfgFile = cfgFile + ".json"
	}

	body, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(body, &config)
	if err != nil {
		log.Panic(err)
	}

	serv := futures.NewClient(config.ApiKey, config.ApiSecret).
		NewCreateOrderService().
		Symbol(symbol).
		Side(futures.SideType(side)).
		PositionSide(futures.PositionSideType(positionSide)).
		Type(futures.OrderType(_type))

	if reduceOnly {
		serv.ReduceOnly(reduceOnly)
	}

	if quantity != "" {
		serv.Quantity(quantity)
	}

	if price != "" {
		serv.Price(price)
	}

	if newClientOrderId != "" {
		serv.NewClientOrderID(newClientOrderId)
	}

	if stopPrice != "" {
		serv.StopPrice(stopPrice)
	}

	if closePosition {
		serv.ClosePosition(closePosition)
	}

	if activationPrice != "" {
		serv.ActivationPrice(activationPrice)
	}

	if timeInForce != "" {
		serv.TimeInForce(futures.TimeInForceType(timeInForce))
	}

	if callbackRate != "" {
		serv.CallbackRate(callbackRate)
	}

	if workingType != "" {
		serv.WorkingType(futures.WorkingType(workingType))
	}

	resp, err := serv.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Info(resp)
}
