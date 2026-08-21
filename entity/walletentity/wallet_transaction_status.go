package walletentity

type WalletTransactionStatus string

const (
	WalletTransactionStatusPending  WalletTransactionStatus = "PENDING"
	WalletTransactionStatusComplete WalletTransactionStatus = "COMPLETED"
	WalletTransactionStatusReversed WalletTransactionStatus = "REVERSED"
	WalletTransactionStatusFailed   WalletTransactionStatus = "FAILED"
)
