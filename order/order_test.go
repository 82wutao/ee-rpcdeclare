package order

import (
	"context"
	"testing"
)

func TestOrderSubmit(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *OrderSubmitReq
		want    *OrderSubmitResp
		wantErr bool
	}{
		{
			name: "submit",
			req: &OrderSubmitReq{
				UserID:     2,
				OrderParam: nil,
			},
			want: &OrderSubmitResp{
				UserID:  2,
				OrderID: 201,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OrderSubmit(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderSubmit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.OrderID != tt.want.OrderID {
				t.Errorf("OrderSubmit() = %v, want %v", got.OrderID, tt.want.OrderID)
			}
		})
	}
}
