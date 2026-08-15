package justanotherpanel

type ActionType string

const (
	ActionTypeServices     ActionType = "services"
	ActionTypeAdd          ActionType = "add"
	ActionTypeStatus       ActionType = "status"
	ActionTypeBalance      ActionType = "balance"
	ActionTypeRefill       ActionType = "refill"
	ActionTypeRefillStatus ActionType = "refill_status"
	ActionTypeCancel       ActionType = "cancel"
)
