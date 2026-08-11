package seeder

import (
	"scaling-up-rest-vs-grpc/internal/data/model"
)

// ---------- Depth 0 ----------

func (f fakeOrderDepthZero) toProto() *model.OrderDepthZero {
	return &model.OrderDepthZero{
		OrderId:               f.OrderID,
		OrderNumber:           f.OrderNumber,
		OrderDate:             f.OrderDate,
		OrderStatus:           f.OrderStatus,
		PaymentMethod:         f.PaymentMethod,
		PaymentStatus:         f.PaymentStatus,
		SubtotalAmount:        float64(f.SubtotalAmount),
		ShippingFee:           float64(f.ShippingFee),
		TotalAmount:           float64(f.TotalAmount),
		PromoCode:             f.PromoCode,
		ShippingMethod:        f.ShippingMethod,
		TrackingNumber:        f.TrackingNumber,
		EstimatedDeliveryDate: f.EstimatedDeliveryDate,
		ItemCount:             f.ItemCount,
		OrderNotes:            f.OrderNotes,
	}
}

// ToOrderDepthZeroResponse generates n fake flat order records and wraps
// them into a model.OrderDepthZeroResponse. Shared by the DepthZero, One,
// Hundred, and Thousand methods, which only differ in n.
func ToOrderDepthZeroResponse(n int) (*model.OrderDepthZeroResponse, error) {
	orders, err := GenerateOrdersDepthZero(n)
	if err != nil {
		return nil, err
	}
	data := make([]*model.OrderDepthZero, len(orders))
	for i, o := range orders {
		data[i] = o.toProto()
	}
	return &model.OrderDepthZeroResponse{Orders: data}, nil
}

// ---------- Depth 2 ----------

func (f fakeOrderDepthTwo) toProto() *model.OrderDepthTwo {
	return &model.OrderDepthTwo{
		OrderId:     f.OrderID,
		OrderNumber: f.OrderNumber,
		OrderDate:   f.OrderDate,
		OrderStatus: f.OrderStatus,
		TotalAmount: float64(f.TotalAmount),
		Customer: &model.CustomerDepthTwo{
			CustomerId:  f.Customer.CustomerID,
			FullName:    f.Customer.FullName,
			Email:       f.Customer.Email,
			Phone:       f.Customer.Phone,
			LoyaltyTier: f.Customer.LoyaltyTier,
			Address: &model.AddressDepthTwo{
				AddressId:     f.Customer.Address.AddressID,
				RecipientName: f.Customer.Address.RecipientName,
				AddressLine1:  f.Customer.Address.AddressLine1,
				PostalCode:    f.Customer.Address.PostalCode,
				AddressType:   f.Customer.Address.AddressType,
			},
		},
	}
}

// ToOrderDepthTwoResponse generates n fake depth-2 order records and
// wraps them into a model.OrderDepthTwoResponse.
func ToOrderDepthTwoResponse(n int) (*model.OrderDepthTwoResponse, error) {
	orders, err := GenerateOrdersDepthTwo(n)
	if err != nil {
		return nil, err
	}
	data := make([]*model.OrderDepthTwo, len(orders))
	for i, o := range orders {
		data[i] = o.toProto()
	}
	return &model.OrderDepthTwoResponse{Orders: data}, nil
}

// ---------- Depth 4 ----------

func (f fakeOrderDepthFour) toProto() *model.OrderDepthFour {
	return &model.OrderDepthFour{
		OrderId:     f.OrderID,
		OrderDate:   f.OrderDate,
		TotalAmount: float64(f.TotalAmount),
		Customer: &model.CustomerDepthFour{
			CustomerId: f.Customer.CustomerID,
			FullName:   f.Customer.FullName,
			Email:      f.Customer.Email,
			Address: &model.AddressDepthFour{
				AddressId:    f.Customer.Address.AddressID,
				AddressLine1: f.Customer.Address.AddressLine1,
				PostalCode:   f.Customer.Address.PostalCode,
				Region: &model.RegionDepthFour{
					RegionId:   f.Customer.Address.Region.RegionID,
					RegionName: f.Customer.Address.Region.RegionName,
					Timezone:   f.Customer.Address.Region.Timezone,
					Country: &model.CountryDepthFour{
						CountryId:   f.Customer.Address.Region.Country.CountryID,
						CountryName: f.Customer.Address.Region.Country.CountryName,
						CountryCode: f.Customer.Address.Region.Country.CountryCode,
					},
				},
			},
		},
	}
}

// ToOrderDepthFourResponse generates n fake depth-4 order records and
// wraps them into a model.OrderDepthFourResponse.
func ToOrderDepthFourResponse(n int) (*model.OrderDepthFourResponse, error) {
	orders, err := GenerateOrdersDepthFour(n)
	if err != nil {
		return nil, err
	}
	data := make([]*model.OrderDepthFour, len(orders))
	for i, o := range orders {
		data[i] = o.toProto()
	}
	return &model.OrderDepthFourResponse{Orders: data}, nil
}
