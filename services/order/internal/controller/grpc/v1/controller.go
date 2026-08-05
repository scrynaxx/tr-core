package v1

import (
	"context"

	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/order"
)

type Controller struct {
}

func NewController() pbv1.OrderServiceServer {
	return &Controller{}
}

func (c Controller) GetOrder(ctx context.Context, orderRequest *pbv1.GetOrderRequest) (*pbv1.GetOrderResponse, error) {
	//TODO implement me
	panic("implement me")
}
