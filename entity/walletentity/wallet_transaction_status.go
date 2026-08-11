package walletentity

type WalletTransactionStatus string

const (
	WalletTransactionStatusPending  WalletTransactionStatus = "PENDING"
	WalletTransactionStatusComplete WalletTransactionStatus = "COMPLETE"
	WalletTransactionStatusReversed WalletTransactionStatus = "REVERSED"
	WalletTransactionStatusFailed   WalletTransactionStatus = "FAILED"
)
