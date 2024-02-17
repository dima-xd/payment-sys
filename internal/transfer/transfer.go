package transfer

import (
	"github.com/dima-xd/payment-sys/internal/iban"
)

type Transfer struct {
	SenderIBAN   iban.IBAN `json:"sender_IBAN"`
	ReceiverIBAN iban.IBAN `json:"receiver_IBAN"`
	Amount       float64   `json:"amount"`
}
