package entity

import (
	"github.com/google/uuid"
)

type PaymentStatus uint8

const (
	PaymentStatusPending  PaymentStatus = 1
	PaymentStatusPaid     PaymentStatus = 2
	PaymentStatusFailed   PaymentStatus = 3
	PaymentStatusRefunded PaymentStatus = 4
)

func (p PaymentStatus) String() string {
	switch p {
	case PaymentStatusPending:
		return "PENDING"
	case PaymentStatusPaid:
		return "PAID"
	case PaymentStatusFailed:
		return "FAILED"
	case PaymentStatusRefunded:
		return "REFUNDED"
	default:
		return "UNKNOWN"
	}
}

type OrderType uint8

const (
	OrderTypeBasic   OrderType = 1
	OrderTypePremium OrderType = 2
	OrderTypeCustom  OrderType = 3
)

func (o OrderType) String() string {
	switch o {
	case OrderTypeBasic:
		return "BASIC"
	case OrderTypePremium:
		return "PREMIUM"
	case OrderTypeCustom:
		return "CUSTOM"
	default:
		return "UNKNOWN"
	}
}

type ProductionStatus uint8

const (
	ProductionStatusPending    ProductionStatus = 1
	ProductionStatusProcessing ProductionStatus = 2
	ProductionStatusCompleted  ProductionStatus = 3
	ProductionStatusFailed     ProductionStatus = 4
)

func (p ProductionStatus) String() string {
	switch p {
	case ProductionStatusPending:
		return "PENDING"
	case ProductionStatusProcessing:
		return "PROCESSING"
	case ProductionStatusCompleted:
		return "COMPLETED"
	case ProductionStatusFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

type DeliveryMethod uint8

const (
	DeliveryMethodDownload DeliveryMethod = 1
	DeliveryMethodEmail    DeliveryMethod = 2
	DeliveryMethodBoth     DeliveryMethod = 3
)

func (d DeliveryMethod) String() string {
	switch d {
	case DeliveryMethodDownload:
		return "DOWNLOAD"
	case DeliveryMethodEmail:
		return "EMAIL"
	case DeliveryMethodBoth:
		return "BOTH"
	default:
		return "UNKNOWN"
	}
}

type Order struct {
	ID                    uuid.UUID        `json:"id" bson:"_id"`
	JobID                 uuid.UUID        `json:"JobId" bson:"JobId"`
	UserEmail             string           `json:"userEmail" bson:"userEmail"`
	UserName              string           `json:"userName" bson:"userName"`
	StripePaymentIntentID string           `json:"stripePaymentIntentId,omitempty" bson:"stripePaymentIntentId,omitempty"`
	StripeCustomerID      string           `json:"stripeCustomerId,omitempty" bson:"stripeCustomerId,omitempty"`
	Amount                float64          `json:"amount" bson:"amount"`
	Currency              string           `json:"currency" bson:"currency"`
	PaymentStatus         PaymentStatus    `json:"paymentStatus" bson:"paymentStatus"`
	PaidAt                int64            `json:"paidAt,omitempty" bson:"paidAt,omitempty"`
	OrderType             OrderType        `json:"orderType" bson:"orderType"`
	Requirements          string           `json:"requirements,omitempty" bson:"requirements,omitempty"`
	ProductionJobID       *uuid.UUID       `json:"productionJobId,omitempty" bson:"productionJobId,omitempty"`
	ProductionStatus      ProductionStatus `json:"productionStatus" bson:"productionStatus"`
	DeliveryMethod        DeliveryMethod   `json:"deliveryMethod" bson:"deliveryMethod"`
	DeliveredAt           int64            `json:"deliveredAt,omitempty" bson:"deliveredAt,omitempty"`
	CustomerNotes         string           `json:"customerNotes,omitempty" bson:"customerNotes,omitempty"`
	SupportTicketID       string           `json:"supportTicketId,omitempty" bson:"supportTicketId,omitempty"`
	IPAddress             string           `json:"ipAddress,omitempty" bson:"ipAddress,omitempty"`
	UserAgent             string           `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
	ExpiresAt             int64            `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
	CreatedAt             int64            `json:"createdAt" bson:"createdAt"`
	UpdatedAt             int64            `json:"updatedAt" bson:"updatedAt"`
}
