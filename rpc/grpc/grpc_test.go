package coregrpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/ostracon/abci/example/kvstore"
	core_grpc "github.com/line/ostracon/rpc/grpc"
	rpctest "github.com/line/ostracon/rpc/test"
)

func TestMain(m *testing.M) {
	// start a tendermint node in the background to test against
	app := kvstore.NewApplication()
	node := rpctest.StartTendermint(app)

	code := m.Run()

	// and shut down proper at the end
	rpctest.StopTendermint(node)
	os.Exit(code)
}

func TestBroadcastTx(t *testing.T) {
	broadcastClient, _ := rpctest.GetGRPCClient()
	res, err := broadcastClient.BroadcastTx(
		context.Background(),
		&core_grpc.RequestBroadcastTx{Tx: []byte("this is a tx")},
	)
	require.NoError(t, err)
	require.EqualValues(t, 0, res.CheckTx.Code)
	require.EqualValues(t, 0, res.DeliverTx.Code)
}

func TestBlockResults(t *testing.T) {
	_, blockClient := rpctest.GetGRPCClient()
	res, err := blockClient.BlockResults(
		context.Background(),
		&core_grpc.RequestBlockResults{Height: 1},
	)
	require.NoError(t, err)
	require.EqualValues(t, 1, res.GetHeight())
}
