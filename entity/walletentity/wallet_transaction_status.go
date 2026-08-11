package walletentity

type WalletTransactionStatus string

const (
	WalletTransactionPending  WalletTransactionStatus = "PENDING"
	WalletTransactionComplete WalletTransactionStatus = "COMPLETE"
	WalletTransactionReversed WalletTransactionStatus = "REVERSED"
	WalletTransactionFailed   WalletTransactionStatus = "FAILED"
)
