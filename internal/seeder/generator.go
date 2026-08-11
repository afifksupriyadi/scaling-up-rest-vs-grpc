package seeder

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// ---------- Depth 0 (shared by DepthZero, One, Hundred, Thousand) ----------

// fakeOrderDepthZero is the gofakeit template for a flat, 15-field order
// record. OrderID, OrderNumber, OrderDate, TrackingNumber, and
// EstimatedDeliveryDate are overwritten manually in
// GenerateOrdersDepthZero, since date/ID/tracking-number formatting needs
// to be deterministic rather than random.
type fakeOrderDepthZero struct {
	OrderID               string
	OrderNumber           string
	OrderDate             string
	OrderStatus           string  `fake:"{randomstring:[pending,processing,shipped,delivered,cancelled]}"`
	PaymentMethod         string  `fake:"{randomstring:[credit_card,debit_card,bank_transfer,e_wallet,cash_on_delivery]}"`
	PaymentStatus         string  `fake:"{randomstring:[paid,unpaid,refunded,failed]}"`
	SubtotalAmount        float32 `fake:"{float32range:10,5000}"`
	ShippingFee           float32 `fake:"{float32range:1,50}"`
	TotalAmount           float32 `fake:"{float32range:11,5050}"`
	PromoCode             string  `fake:"{randomstring:[NEWYEAR10,FREESHIP,SAVE20,WELCOME5,NONE]}"`
	ShippingMethod        string  `fake:"{randomstring:[standard,express,same_day,pickup]}"`
	TrackingNumber        string
	EstimatedDeliveryDate string
	ItemCount             int32  `fake:"{number:1,10}"`
	OrderNotes            string `fake:"{sentence:6}"`
}

// GenerateOrdersDepthZero produces n fake flat order records, used by the
// DepthZero, One, Hundred, and Thousand methods, which share this shape
// and differ only in how many records are requested.
func GenerateOrdersDepthZero(n int) ([]fakeOrderDepthZero, error) {
	orders := make([]fakeOrderDepthZero, n)
	now := time.Now()
	for i := range orders {
		if err := gofakeit.Struct(&orders[i]); err != nil {
			return nil, fmt.Errorf("generate order depth-zero %d: %w", i, err)
		}
		orders[i].OrderID = fmt.Sprintf("ORD-%08d", i+1)
		orders[i].OrderNumber = fmt.Sprintf("ORDNUM-%08d", i+1)
		orders[i].TrackingNumber = fmt.Sprintf("TRK%08d", i+1)
		orders[i].OrderDate = now.AddDate(0, 0, -i).Format("2006-01-02")
		orders[i].EstimatedDeliveryDate = now.AddDate(0, 0, 7-i%7).Format("2006-01-02")
	}
	return orders, nil
}

// ---------- Depth 2 ----------

// fakeOrderDepthTwo is the gofakeit template for the Order entity at depth 2.
type fakeOrderDepthTwo struct {
	OrderID     string
	OrderNumber string
	OrderDate   string
	OrderStatus string  `fake:"{randomstring:[pending,processing,shipped,delivered,cancelled]}"`
	TotalAmount float32 `fake:"{float32range:11,5050}"`
}

// fakeCustomerDepthTwo is the gofakeit template for the Customer entity at depth 2.
type fakeCustomerDepthTwo struct {
	CustomerID  string
	FullName    string `fake:"{name}"`
	Email       string `fake:"{email}"`
	Phone       string `fake:"{phone}"`
	LoyaltyTier string `fake:"{randomstring:[bronze,silver,gold,platinum]}"`
}

// fakeAddressDepthTwo is the gofakeit template for the Address entity at depth 2.
type fakeAddressDepthTwo struct {
	AddressID     string
	RecipientName string `fake:"{name}"`
	AddressLine1  string `fake:"{street}"`
	PostalCode    string `fake:"{zip}"`
	AddressType   string `fake:"{randomstring:[home,office,warehouse]}"`
}

// fakeOrderDepthTwoDocument bundles the three depth-2 entities into one
// single-object chain, matching the OrderDepthTwoDocument proto message.
type fakeOrderDepthTwoDocument struct {
	Order    fakeOrderDepthTwo
	Customer fakeCustomerDepthTwo
	Address  fakeAddressDepthTwo
}

// GenerateOrdersDepthTwo produces n fake depth-2 order documents.
func GenerateOrdersDepthTwo(n int) ([]fakeOrderDepthTwoDocument, error) {
	docs := make([]fakeOrderDepthTwoDocument, n)
	now := time.Now()
	for i := range docs {
		if err := gofakeit.Struct(&docs[i]); err != nil {
			return nil, fmt.Errorf("generate order depth-two %d: %w", i, err)
		}
		docs[i].Order.OrderID = fmt.Sprintf("ORD-%08d", i+1)
		docs[i].Order.OrderNumber = fmt.Sprintf("ORDNUM-%08d", i+1)
		docs[i].Order.OrderDate = now.AddDate(0, 0, -i).Format("2006-01-02")
		docs[i].Customer.CustomerID = fmt.Sprintf("CUST-%08d", i+1)
		docs[i].Address.AddressID = fmt.Sprintf("ADDR-%08d", i+1)
	}
	return docs, nil
}

// ---------- Depth 4 ----------

// fakeOrderDepthFour is the gofakeit template for the Order entity at depth 4.
type fakeOrderDepthFour struct {
	OrderID     string
	OrderDate   string
	TotalAmount float32 `fake:"{float32range:11,5050}"`
}

// fakeCustomerDepthFour is the gofakeit template for the Customer entity at depth 4.
type fakeCustomerDepthFour struct {
	CustomerID string
	FullName   string `fake:"{name}"`
	Email      string `fake:"{email}"`
}

// fakeAddressDepthFour is the gofakeit template for the Address entity at depth 4.
type fakeAddressDepthFour struct {
	AddressID    string
	AddressLine1 string `fake:"{street}"`
	PostalCode   string `fake:"{zip}"`
}

// fakeRegionDepthFour is the gofakeit template for the Region entity at depth 4.
type fakeRegionDepthFour struct {
	RegionID   string
	RegionName string `fake:"{state}"`
	Timezone   string `fake:"{timezone}"`
}

// fakeCountryDepthFour is the gofakeit template for the Country entity at depth 4.
type fakeCountryDepthFour struct {
	CountryID   string
	CountryName string `fake:"{country}"`
	CountryCode string `fake:"{countryabr}"`
}

// fakeOrderDepthFourDocument bundles the five depth-4 entities into one
// single-object chain, matching the OrderDepthFourDocument proto message.
type fakeOrderDepthFourDocument struct {
	Order    fakeOrderDepthFour
	Customer fakeCustomerDepthFour
	Address  fakeAddressDepthFour
	Region   fakeRegionDepthFour
	Country  fakeCountryDepthFour
}

// GenerateOrdersDepthFour produces n fake depth-4 order documents.
func GenerateOrdersDepthFour(n int) ([]fakeOrderDepthFourDocument, error) {
	docs := make([]fakeOrderDepthFourDocument, n)
	now := time.Now()
	for i := range docs {
		if err := gofakeit.Struct(&docs[i]); err != nil {
			return nil, fmt.Errorf("generate order depth-four %d: %w", i, err)
		}
		docs[i].Order.OrderID = fmt.Sprintf("ORD-%08d", i+1)
		docs[i].Order.OrderDate = now.AddDate(0, 0, -i).Format("2006-01-02")
		docs[i].Customer.CustomerID = fmt.Sprintf("CUST-%08d", i+1)
		docs[i].Address.AddressID = fmt.Sprintf("ADDR-%08d", i+1)
		docs[i].Region.RegionID = fmt.Sprintf("REG-%08d", i+1)
		docs[i].Country.CountryID = fmt.Sprintf("CTY-%08d", i+1)
	}
	return docs, nil
}
