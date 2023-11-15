package blockchain

import (
	"context"
	"errors"
	"log"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/global"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/wallet"
)

func Transfer(receiverAddress string, amount uint64, comment string) error {
	client, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		log.Printf("Unable to create lite client: %v\n", err)
		return err
	}

	w, err := wallet.DefaultWalletFromSeed(global.Conf.AppConf.TonSeed, client)
	if err != nil {
		log.Printf("Unable to create wallet: %v\n", err)
		return err
	}
	balance, err := w.GetBalance(context.TODO())
	if err != nil {
		log.Printf("Unable to get balance: %v\n", err)
		return err
	}
	needTonAmount := amount + 100_000_000
	if balance < needTonAmount {
		log.Printf("balance is not enough: now %d but need %d\n", balance, needTonAmount)
		return errors.New("balance is not enough")
	}

	simpleTransfer := wallet.SimpleTransfer{
		Amount:  tlb.Grams(amount),
		Address: tongo.MustParseAddress(receiverAddress).ID,
		Comment: comment,
	}

	err = w.Send(context.TODO(), simpleTransfer)
	if err != nil {
		log.Printf("Unable to generate transfer message: %v", err)
		return err
	}

	return nil
}
