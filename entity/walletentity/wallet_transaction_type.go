package walletentity

type WalletTransactionType string

// DEPOSIT / WITHDRAW / REFUND

const (
	WalletTransactionTypeDeposit  WalletTransactionType = "DEPOSIT"
	WalletTransactionTypeWithdraw WalletTransactionType = "WITHDRAW"
	WalletTransactionTypeRefund   WalletTransactionType = "REFUND"
)
