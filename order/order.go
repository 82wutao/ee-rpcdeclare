package order

import (
	"context"

	"github.com/82wutao/ee-rpcdeclare/rpcx"
)

type OrderService int

func (os *OrderService) HandleName() string {
	return "order-service"
}

type OrderQueryReq struct {
	UserID int
}
type OrderQueryResp struct {
	UserID int
	Orders []int
}

type OrderSubmitReq struct {
	UserID     int
	OrderParam interface{}
}
type OrderSubmitResp struct {
	UserID  int
	OrderID int
}
type OrderCancelReq struct {
	UserID  int
	OrderID int
}
type OrderCancelResp struct {
	UserID int
	Suc    bool
}

func OrderQuery(ctx context.Context, req *OrderQueryReq) (*OrderQueryResp, error) {
	var os OrderService
	cli, err := rpcx.NewClientByP2P(rpcx.HostPort{Proto: "tcp", Host: "localhost", Port: 9000}, os.HandleName())
	if err != nil {
		return nil, err
	}

	var resp OrderQueryResp
	err = cli.Call(context.Background(), "OrderQuery", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
func OrderSubmit(ctx context.Context, req *OrderSubmitReq) (*OrderSubmitResp, error) {
	var os OrderService
	cli, err := rpcx.NewClientByP2P(rpcx.HostPort{Proto: "tcp", Host: "localhost", Port: 9000}, os.HandleName())
	if err != nil {
		return nil, err
	}

	var resp OrderSubmitResp
	err = cli.Call(context.Background(), "OrderSubmit", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
func OrderCancel(ctx context.Context, req *OrderCancelReq) (*OrderCancelResp, error) {
	var os OrderService
	cli, err := rpcx.NewClientByP2P(rpcx.HostPort{Proto: "tcp", Host: "localhost", Port: 9000}, os.HandleName())
	if err != nil {
		return nil, err
	}

	var resp OrderCancelResp
	err = cli.Call(context.Background(), "OrderCancel", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
