package mm

type wallet struct {
	ID         int32
	userID     string
	walletType string
}

type account struct {
	ID          int32
	cents       int64
	accountType string
	walletID    int32
}

type transaction struct {
	ID                       int32
	pid                      string
	srcUserID                string
	dstUserID                string
	srcAccountWalletID       string
	dstAccountWalletID       string
	srcAccountID             string
	dstAccountID             string
	srcAccountType           string
	dstAccountType           string
	finalDstMerchantWalletID int32
	amount                   int64
}
