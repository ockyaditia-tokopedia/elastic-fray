package marketplace

import (
	"time"

	"github.com/elastic-fray/entity/user"
)

type (
	Promo struct {
		OrderID          int64         `json:"order_id"`
		PaymentID        int64         `json:"payment_id"`
		ShopID           int64         `json:"shop_id"`
		InvoiceRefNum    string        `json:"invoice_ref_num"`
		PaymentGatewayID int           `json:"payment_gateway_id"`
		Amount           float64       `json:"amount"`
		ShippingID       int           `json:"shipping_id"`
		IsGoldShop       bool          `json:"is_gold_shop"`
		SellerData       user.UserData `json:"seller_data"`
		BuyerData        user.UserData `json:"buyer_data"`
		PromoDetail      PromoData     `json:"promo_detail"`
		DeviceID         string        `json:"device_id"`
		CreateTime       time.Time     `json:"create_time"`
		Source           string        `json:"source"`
		Platform         string        `json:"platform"`
		GroupID          int           `json:"group_id"`
	}

	PromoData struct {
		PromoID              int64      `json:"promo_id,omitempty" validate:"required"`
		PromoName            string     `json:"promo_name,omitempty"`
		VoucherCode          string     `json:"voucher_code,omitempty"`
		AdsID                string     `json:"ads_id,omitempty" validate:"required"`
		Benefit              string     `json:"benefit,omitempty"`
		FingerPrint          string     `json:"finger_print,omitempty"`
		GatewayCode          string     `json:"gateway_code,omitempty"`
		IsBackdoor           bool       `json:"is_backdoor,omitempty"`
		IsUnlimited          bool       `json:"is_unlimited,omitempty"`
		PromoCodeUsageID     int        `json:"promo_code_usage_id,omitempty"`
		PromoCodeID          int        `json:"promo_code_id,omitempty"`
		BinaryPromoType      int64      `json:"binary_promo_type,omitempty"`
		ProductCode          string     `json:"product_code,omitempty"`
		Status               int        `json:"status,omitempty"`
		Code                 string     `json:"code,omitempty"`
		Similarity           float64    `json:"similarity,omitempty"`
		AdsIDChecking        bool       `json:"ads_id_checking,omitempty" validate:"required"`
		IsPostCheck          bool       `json:"is_post_check,omitempty" validate:"required"`
		IsCoupon             bool       `json:"is_coupon,omitempty"`
		CashbackEarned       float64    `json:"cashback_earned,omitempty"`
		GroupID              int        `json:"group_id,omitempty"`
		DiscountAmount       float64    `json:"discount_amount,omitempty"`
		IsExclusive          bool       `json:"is_exclusive,omitempty"`
		PromoRule            *PromoRule `json:"promo_rule,omitempty" validate:"required"`
		IsFraud              bool       `json:"is_fraud"`
		IsMaxBenefit         bool       `json:"is_max_benefit,omitempty"`
		BenefitAmount        float64    `json:"benefit_amount,omitempty"`
		BenefitPercentage    float64    `json:"benefit_percentage,omitempty"`
		MaxBenefitPercentage float64    `json:"max_benefit_percentage,omitempty"`
		TokoPointsEarned     float64    `json:"tokopoints_earned,omitempty"`
		TokoPointsRate       float64    `json:"tokopoints_rate,omitempty"`
	}

	PromoRule struct {
		Counter  int       `json:"counter,omitempty" validate:"required"`
		Expired  int64     `json:"expired,omitempty"`
		Status   int       `json:"status,omitempty"`
		Days     int       `json:"days,omitempty"`
		Coverage *Coverage `json:"coverage,omitempty"`
	}

	Coverage struct {
		Categories []int    `json:"categories,omitempty"`
		Products   []string `json:"products,omitempty"`
		ServiceIds []int64  `json:"service_ids,omitempty"`
		PromoID    int      `json:"promo_id,omitempty"`
	}
)
