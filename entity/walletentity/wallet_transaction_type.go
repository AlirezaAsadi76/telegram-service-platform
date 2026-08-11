package walletentity

type WalletTransactionType string

// DEPOSIT / WITHDRAW / REFUND

const (
	WalletTransactionDeposit  WalletTransactionType = "DEPOSIT"
	WalletTransactionWithdraw WalletTransactionType = "WITHDRAW"
	WalletTransactionRefund   WalletTransactionType = "REFUND"
)
