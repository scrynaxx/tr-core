package inner

import (
	pbinner "github.com/scrynaxx/tr-core/contracts/generated/services/inner/order"
)

type Controller struct {
}

func NewController() pbinner.OrderServiceServer {
	return &Controller{}
}
