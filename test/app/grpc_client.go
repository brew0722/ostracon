package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"context"

	tmjson "github.com/line/ostracon/libs/json"
	coregrpc "github.com/line/ostracon/rpc/grpc"
)

var grpcAddr = "tcp://localhost:36656"

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Must enter a transaction to send (hex)")
		os.Exit(1)
	}
	tx := args[1]
	txBytes, err := hex.DecodeString(tx)
	if err != nil {
		fmt.Println("Invalid hex", err)
		os.Exit(1)
	}

	broadcastClientGRPC, blockClientGRPC := coregrpc.StartGRPCClient(grpcAddr)
	{
		res, err := broadcastClientGRPC.BroadcastTx(context.Background(), &coregrpc.RequestBroadcastTx{Tx: txBytes})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bz, err := tmjson.Marshal(res)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(bz))
	}

	{
		res, err := blockClientGRPC.BlockResults(context.Background(), &coregrpc.RequestBlockResults{Height: 1})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bz, err := tmjson.Marshal(res)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(bz))
	}
}
