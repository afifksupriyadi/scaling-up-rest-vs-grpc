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

func (f fakeOrderDepthTwoDocument) toProto() *model.OrderDepthTwoDocument {
	return &model.OrderDepthTwoDocument{
		Order: &model.OrderDepthTwo{
			OrderId:     f.Order.OrderID,
			OrderNumber: f.Order.OrderNumber,
			OrderDate:   f.Order.OrderDate,
			OrderStatus: f.Order.OrderStatus,
			TotalAmount: float64(f.Order.TotalAmount),
		},
		Customer: &model.CustomerDepthTwo{
			CustomerId:  f.Customer.CustomerID,
			FullName:    f.Customer.FullName,
			Email:       f.Customer.Email,
			Phone:       f.Customer.Phone,
			LoyaltyTier: f.Customer.LoyaltyTier,
		},
		Address: &model.AddressDepthTwo{
			AddressId:     f.Address.AddressID,
			RecipientName: f.Address.RecipientName,
			AddressLine1:  f.Address.AddressLine1,
			PostalCode:    f.Address.PostalCode,
			AddressType:   f.Address.AddressType,
		},
	}
}

// ToOrderDepthTwoResponse generates n fake depth-2 order documents and
// wraps them into a model.OrderDepthTwoResponse.
func ToOrderDepthTwoResponse(n int) (*model.OrderDepthTwoResponse, error) {
	docs, err := GenerateOrdersDepthTwo(n)
	if err != nil {
		return nil, err
	}
	data := make([]*model.OrderDepthTwoDocument, len(docs))
	for i, d := range docs {
		data[i] = d.toProto()
	}
	return &model.OrderDepthTwoResponse{Orders: data}, nil
}

// ---------- Depth 4 ----------

func (f fakeOrderDepthFourDocument) toProto() *model.OrderDepthFourDocument {
	return &model.OrderDepthFourDocument{
		Order: &model.OrderDepthFour{
			OrderId:     f.Order.OrderID,
			OrderDate:   f.Order.OrderDate,
			TotalAmount: float64(f.Order.TotalAmount),
		},
		Customer: &model.CustomerDepthFour{
			CustomerId: f.Customer.CustomerID,
			FullName:   f.Customer.FullName,
			Email:      f.Customer.Email,
		},
		Address: &model.AddressDepthFour{
			AddressId:    f.Address.AddressID,
			AddressLine1: f.Address.AddressLine1,
			PostalCode:   f.Address.PostalCode,
		},
		Region: &model.RegionDepthFour{
			RegionId:   f.Region.RegionID,
			RegionName: f.Region.RegionName,
			Timezone:   f.Region.Timezone,
		},
		Country: &model.CountryDepthFour{
			CountryId:   f.Country.CountryID,
			CountryName: f.Country.CountryName,
			CountryCode: f.Country.CountryCode,
		},
	}
}

// ToOrderDepthFourResponse generates n fake depth-4 order documents and
// wraps them into a model.OrderDepthFourResponse.
func ToOrderDepthFourResponse(n int) (*model.OrderDepthFourResponse, error) {
	docs, err := GenerateOrdersDepthFour(n)
	if err != nil {
		return nil, err
	}
	data := make([]*model.OrderDepthFourDocument, len(docs))
	for i, d := range docs {
		data[i] = d.toProto()
	}
	return &model.OrderDepthFourResponse{Orders: data}, nil
}
