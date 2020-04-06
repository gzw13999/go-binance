package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/binance-exchange/go-binance"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(os.Getenv("BINANCE_SECRET")),
	}
	ctx, cancelCtx := context.WithCancel(context.Background())
	// use second return value for cancelling request
	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		os.Getenv("BINANCE_APIKEY"),
		hmacSigner,
		logger,
		ctx,
	)
	b := binance.NewBinance(binanceService)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	kech, done, err := b.TradeWebsocket(binance.TradeWebsocketRequest{
		Symbol: "ETHBTC",
	})
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case ke := <-kech:
				fmt.Printf("%#v\n", ke)
			case <-done:
				break
			}
		}
	}()

	fmt.Println("waiting for interrupt")
	<-interrupt
	fmt.Println("canceling context")
	cancelCtx()
	fmt.Println("waiting for signal")
	<-done
	fmt.Println("exit")
	return

	kl, err := b.Klines(binance.KlinesRequest{
		Symbol:   "BNBETH",
		Interval: binance.Hour,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", kl)

	newOrder, err := b.NewOrder(binance.NewOrderRequest{
		Symbol:      "BNBETH",
		Quantity:    1,
		Price:       999,
		Side:        binance.SideSell,
		TimeInForce: binance.GTC,
		Type:        binance.TypeLimit,
		Timestamp:   time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(newOrder)

	res2, err := b.QueryOrder(binance.QueryOrderRequest{
		Symbol:     "BNBETH",
		OrderID:    newOrder.OrderID,
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", res2)

	res4, err := b.OpenOrders(binance.OpenOrdersRequest{
		Symbol:     "BNBETH",
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", res4)

	res3, err := b.CancelOrder(binance.CancelOrderRequest{
		Symbol:    "BNBETH",
		OrderID:   newOrder.OrderID,
		Timestamp: time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", res3)

	res5, err := b.AllOrders(binance.AllOrdersRequest{
		Symbol:     "BNBETH",
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", res5[0])

	res6, err := b.Account(binance.AccountRequest{
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", res6)

	res7, err := b.MyTrades(binance.MyTradesRequest{
		Symbol:     "BNBETH",
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", res7)

	res9, err := b.DepositHistory(binance.HistoryRequest{
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", res9)

	res8, err := b.WithdrawHistory(binance.HistoryRequest{
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", res8)

	ds, err := b.StartUserDataStream()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", ds)

	err = b.KeepAliveUserDataStream(ds)
	if err != nil {
		panic(err)
	}

	err = b.CloseUserDataStream(ds)
	if err != nil {
		panic(err)
	}
}
