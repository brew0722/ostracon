package coregrpc

import (
	"context"

	abci "github.com/line/ostracon/abci/types"
	core "github.com/line/ostracon/rpc/core"
	rpctypes "github.com/line/ostracon/rpc/jsonrpc/types"
)

type broadcastAPI struct {
}

func (bapi *broadcastAPI) Ping(ctx context.Context, req *RequestPing) (*ResponsePing, error) {
	// kvstore so we can check if the server is up
	return &ResponsePing{}, nil
}

func (bapi *broadcastAPI) BroadcastTx(ctx context.Context, req *RequestBroadcastTx) (*ResponseBroadcastTx, error) {
	// NOTE: there's no way to get client's remote address
	// see https://stackoverflow.com/questions/33684570/session-and-remote-ip-address-in-grpc-go
	res, err := core.BroadcastTxCommit(&rpctypes.Context{}, req.Tx)
	if err != nil {
		return nil, err
	}

	return &ResponseBroadcastTx{
		CheckTx: &abci.ResponseCheckTx{
			Code: res.CheckTx.Code,
			Data: res.CheckTx.Data,
			Log:  res.CheckTx.Log,
		},
		DeliverTx: &abci.ResponseDeliverTx{
			Code: res.DeliverTx.Code,
			Data: res.DeliverTx.Data,
			Log:  res.DeliverTx.Log,
		},
	}, nil
}

type blockAPI struct {
}

func (bapi *blockAPI) Block(ctx context.Context, req *RequestBlock) (*ResponseBlock, error) {
	res, err := core.Block(&rpctypes.Context{}, &req.Height)
	if err != nil {
		return nil, err
	}

	protoBlock, err := res.Block.ToProto()
	if err != nil {
		return nil, err
	}

	return &ResponseBlock{
		Block: protoBlock,
	}, nil
}

func (bapi *blockAPI) BlockResults(ctx context.Context, req *RequestBlockResults) (*ResponseBlockResults, error) {
	res, err := core.BlockResults(&rpctypes.Context{}, &req.Height)
	if err != nil {
		return nil, err
	}

	return &ResponseBlockResults{
		Height:     res.Height,
		TxsResults: res.TxsResults,
		ResBeginBlock: &abci.ResponseBeginBlock{
			Events: res.BeginBlockEvents,
		},
		ResEndBlock: &abci.ResponseEndBlock{
			ValidatorUpdates:      res.ValidatorUpdates,
			ConsensusParamUpdates: res.ConsensusParamUpdates,
			Events:                res.EndBlockEvents,
		},
	}, nil
}
