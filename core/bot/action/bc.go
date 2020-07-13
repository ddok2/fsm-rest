package action

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"sync"
	"time"

	cfg "blockchain.automation/common/config"
	"blockchain.automation/core/bot"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const ginGroup = "/api/v2"

// BcBotAction represents Action for bc
type BcBotAction struct {
	jobs   *bot.Jobs
	client *resty.Client
	mutex  *sync.RWMutex
}

type wallet struct {
	WalletId     string `json:"wallet_id"`
	Balance      string `json:"balance"`
	UserId       string `json:"user_id"`
	UserName     string `json:"user_name"`
	WalletStatus string `json:"wallet_status"`
	Created      string `json:"created"`
	TokenId      string `json:"token_id"`
	TokenName    string `json:"token_name"`
	TokenType    string `json:"token_type"`
}

type issue struct {
	TxID          string `json:"txID"`
	WalletAddress string `json:"WalletAddress"`
	Amount        string `json:"amount"`
	Balance       string `json:"balance"`
	TxTime        string `json:"txTime"`
}

// @Param txFlag
// 1: TxRemittanceW2W,
// 2: TxSellToken,
// 3: TxSellElmo,
// 4: TxChargeHes,
// 6: TxRemittanceW2C,
// 7: TxRemittanceC2W,
// 8: TxRemittanceC2C,
// 9:CashOut
type trade struct {
	TxID                  string `json:"txID"`
	SenderWalletAddress   string `json:"senderWalletAddress"`
	SenderBalance         string `json:"senderBalance"`
	ReceiverWalletAddress string `json:"receiverWalletAddress"`
	ReceiverBalance       string `json:"receiverBalance"`
	Amount                string `json:"amount"`
	Fee                   string `json:"fee"`
	TxFlag                string `json:"txFlag"`
	TxTime                string `json:"txTime"`
	FeeToGo               string `json:"feeToGo"`
}

// NewBcBotAction return an initialized BcBotAction
func NewBcBotAction(j *bot.Jobs) *BcBotAction {
	// client := resty.New()
	// client.
	// 	SetRetryCount(3).
	// 	SetRetryWaitTime(10 * time.Second)
	return &BcBotAction{jobs: j, client: nil, mutex: new(sync.RWMutex)}
}

// NewWallet returns an rest call of /wallet
func (b *BcBotAction) NewWallet(bot *bot.Bot) (bool, error) {
	w := &wallet{
		WalletId:     bot.Id,
		Balance:      strconv.FormatFloat(bot.Balance, 'f', -1, 64),
		UserId:       bot.Id,
		UserName:     bot.Id,
		WalletStatus: "active",
		Created:      time.Now().Format("2006-01-02T15:04:05.000Z"),
		TokenId:      "general_token_id",
		TokenName:    "general_token",
		TokenType:    "GENERL",
	}

	isDone, err := b.requestHandler("/wallet", w)
	if err != nil {
		return isDone, err
	}
	return isDone, nil
}

// Issue returns an rest call of /issue
// issue coin admin wallet only.
func (b *BcBotAction) Issue(_ *bot.Bot) (bool, error) {

	adminBot := b.GetAdminBot()

	balance := strconv.FormatFloat(adminBot.Balance, 'f', -1, 64)
	chargeAmount, err := strconv.ParseFloat(b.jobs.Config.AdminChargeAmount, 64)
	if err != nil {
		return false, err
	}

	newIssue := &issue{
		TxID:          uuid.Must(uuid.NewUUID()).String(),
		WalletAddress: adminBot.Id,
		Amount:        b.jobs.Config.AdminChargeAmount,
		Balance:       balance,
		TxTime:        time.Now().Format("2006-01-02T15:04:05.000Z"),
	}

	isDone, err := b.requestHandler("/issue", newIssue)
	if err != nil {
		return false, err
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.jobs.AdminBot.Balance -= chargeAmount

	return isDone, nil
}

// Trade returns an rest call of /trade
func (b *BcBotAction) Trade(origin *bot.Bot, target *bot.Bot) (bool, error) {
	var isDone bool

	chargeAmount, err := strconv.ParseFloat(b.jobs.Config.ChargeAmount, 64)
	if err != nil {
		return false, err
	}
	amount, fee := getAmountAndFee(b.jobs.Config)
	if amount == -1 {
		return false, err
	}

	originBalanceStr := strconv.FormatFloat(origin.Balance, 'f', -1, 64)
	targetBalanceStr := strconv.FormatFloat(target.Balance, 'f', -1, 64)
	body := trade{
		TxID:                  uuid.Must(uuid.NewUUID()).String(),
		SenderWalletAddress:   origin.Id,
		SenderBalance:         originBalanceStr,
		ReceiverWalletAddress: target.Id,
		ReceiverBalance:       targetBalanceStr,
		Amount:                b.jobs.Config.RemittanceAmount,
		Fee:                   b.jobs.Config.RemittanceFee,
		TxFlag:                "1",
		TxTime:                time.Now().Format("2006-01-02T15:04:05.000Z"),
		FeeToGo:               b.jobs.Config.AdminId,
	}

	if origin.Id == body.FeeToGo {
		// issue to user from admin
		body.Amount = b.jobs.Config.ChargeAmount
		body.TxFlag = "2"
		isDone, err = b.requestHandler("/trade", body)
		if err != nil {
			return false, err
		}

		b.mutex.Lock()
		defer b.mutex.Unlock()
		b.jobs.AdminBot.Balance -= chargeAmount

		return true, nil
	}

	// trade w2w
	isDone, err = b.requestHandler("/trade", body)
	if err != nil {
		return false, err
	}

	// cal
	b.mutex.Lock()
	defer b.mutex.Unlock()

	origin.Balance -= amount - fee
	target.Balance += amount

	if fee > 0 {
		b.jobs.AdminBot.Balance += fee
	}

	return isDone, nil
}

// GetAdminBot return admin bot from Jobs
func (b *BcBotAction) GetAdminBot() *bot.Bot {
	return b.jobs.AdminBot
}

func (b *BcBotAction) requestHandler(target string, body interface{}) (bool, error) {
	client := resty.New()
	// client.
	// 	SetRetryCount(3).
	// 	SetRetryWaitTime(10 * time.Second)
	u := url.URL{
		Scheme: "http",
		Host:   b.jobs.Config.BoosterAddr + ":" + b.jobs.Config.BoosterPort,
		Path:   ginGroup + target,
	}

	res, err := client.R().
		EnableTrace().
		SetBody(body).
		Post(u.String())
	if err != nil {
		return false, err
	}

	if res.StatusCode() > 201 {
		return false, fmt.Errorf("something wrong: %v", res.Request)
	}
	return true, nil
}

// Handler handled bot 's actions
func (b *BcBotAction) Handler() {
	if b.jobs.Queue.Size() == 0 {
		for _, bt := range b.jobs.BotList {
			_, err := b.NewWallet(bt)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = bt.SetState(bot.BotStateNew)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		// issue to admin wallet
		if _, err := b.Issue(b.jobs.AdminBot); err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	for {
		get, err := b.jobs.Queue.Get()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		target := get.(*bot.Bot)
		event := target.GetState()

		switch event {
		case bot.BotStateNew:
			if err = target.SetState(bot.BotStateTrade); err != nil {
				fmt.Println(err.Error())
				return
			}

		case bot.BotStateIssue:
			chargeAmount, err := strconv.ParseFloat(b.jobs.Config.ChargeAmount, 64)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			go func() {
				if b.jobs.AdminBot.Balance < chargeAmount {
					if _, err = b.Issue(b.jobs.AdminBot); err != nil {
						fmt.Println(err.Error())
						return
					}
				}
				if _, err = b.Trade(b.jobs.AdminBot, target); err != nil {
					fmt.Println(err.Error())
					return
				}
			}()

			if err = target.SetState(bot.BotStateTrade); err != nil {
				fmt.Println(err.Error())
				return
			}

		case bot.BotStateTrade:
			amount, fee := getAmountAndFee(b.jobs.Config)
			if amount == -1 {
				return
			}

			if target.Balance > amount+fee {
				go func() {
					receiver := b.jobs.BotList[rand.Intn(b.jobs.Config.MemberCount)]
					if _, err = b.Trade(target, receiver); err != nil {
						fmt.Println(err.Error())
						return
					}
				}()
				if err = target.SetState(bot.BotStateTrade); err != nil {
					fmt.Println(err.Error())
					return
				}
			} else {
				// trade admin -> wallet
				err = target.SetState(bot.BotStateIssue)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}

		}
		b.jobs.Queue.Put(target)
		fmt.Println(target.Id, ": ", target.GetState())

		d, err := strconv.Atoi(b.jobs.Config.TxInterval)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		time.Sleep(time.Duration(d) * time.Millisecond)
	}
}

func getAmountAndFee(config *cfg.Config) (float64, float64) {
	amount, err := strconv.ParseFloat(config.RemittanceAmount, 64)
	if err != nil {
		fmt.Println(err.Error())
		return -1, -1
	}
	fee, err := strconv.ParseFloat(config.RemittanceFee, 64)
	if err != nil {
		fmt.Println(err.Error())
		return -1, -1
	}
	return amount, fee
}
