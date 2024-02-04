package outsideErrors

import "errors"

var (
	NotEnoughMoney          = errors.New("not enough money")
	NoSuchWallet            = errors.New("no such wallet")
	NoSuchOutgoingWallet    = errors.New("no such outgoing wallet")
	NoSuchDestinationWallet = errors.New("no such destination wallet")
	InvalidMoneyFormat      = errors.New("invalid money format")
	InvalidWalletFormat     = errors.New("invalid wallet format")
)
