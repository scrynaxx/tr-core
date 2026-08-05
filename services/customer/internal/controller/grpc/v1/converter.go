package v1

import (
	pbv1 "github.com/scrynaxx/tr-core/contracts/generated/services/v1/customer"
	"github.com/scrynaxx/tr-core/services/customer/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toCustomer(customer model.Customer) *pbv1.Customer {
	result := &pbv1.Customer{
		Id:        customer.ID.String(),
		Name:      customer.Name,
		Phone:     customer.Phone,
		Email:     customer.Email,
		CreatedAt: timestamppb.New(customer.CreatedAt),
		UpdatedAt: timestamppb.New(customer.UpdatedAt),
	}
	if customer.ArchivedAt != nil {
		result.ArchivedAt = timestamppb.New(*customer.ArchivedAt)
	}

	return result
}

func toListCustomersResponse(customers []model.Customer) *pbv1.ListCustomersResponse {
	result := &pbv1.ListCustomersResponse{
		Customers: make([]*pbv1.Customer, len(customers)),
	}

	for i := range customers {
		result.Customers[i] = toCustomer(customers[i])
	}

	return result
}
