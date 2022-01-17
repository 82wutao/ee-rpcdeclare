package order

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

// func (os *OrderService) Query(ctx context.Context, req *OrderQueryReq, resp *OrderQueryResp) error {
// 	return nil
// }
// func (os *OrderService) Submit(ctx context.Context, req *Req, resp *Resp) error {
// 	return nil
// }
// func (os *OrderService) Cancel(ctx context.Context, req *Req, resp *Resp) error {
// 	return nil
// }
