package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"blockchain.automation/cmd/common"
	"blockchain.automation/utils"
	uuid "github.com/satori/go.uuid"
)

type Members struct {
	Members     []*Member
	BoosterAddr string
	Count       int
	Prefix      string
}

var (
	config   *common.Config
	members  *Members
	exchange Member
	logger   *utils.Logger
	ch       chan bool

	pause bool
)

const (
	RemittanceAmount = 100
	RemittanceFee    = 3
	ChargeAmount     = 10000
)

func SetExit(flag bool) {
	pause = flag
}

func Resume() {
	ch <- true
}

func NewMembers() *Members {
	if members == nil {
		members = new(Members)
	}

	return members
}

func (m *Members) Initialize() {

	logger = utils.NewLogger()
	config = common.NewConfig()

	m.BoosterAddr = config.BoosterAddr
	m.Count = config.MemberCount
	m.Prefix = config.MemberPrefix

	m.Members = make([]*Member, m.Count)

	messageQueue = NewMessageQueue()
	messageQueue.Initialize()

	ch = make(chan bool)

	pause = false

}

func (m *Members) Schedule() {

	logger.Info("Scheduler start...")
	for {
		// <- ch
		if pause == true {
			<-ch
		}

		if member, err := messageQueue.GetMessage(); err != nil {

		} else {

			event := member.State()

			switch event {
			case Closed:
				err := member.SetState(Register)
				if err != nil {
					logger.Info(err)
				}

			case Register:
				err := member.SetState(Transfer)
				if err != nil {
					logger.Info(err)
				}

			case Transfer:
				if member.Balance > RemittanceAmount+RemittanceFee {
					go func() {
						m.Remittance(member)
					}()
					err := member.SetState(Transfer)
					if err != nil {
						logger.Info(err)
					}
				} else if member.Balance <= 0 {
					err := member.SetState(Charge)
					if err != nil {
						logger.Info(err)
					}
				}
			case Charge:
				go func() {
					m.SellElmo(member)
				}()
				err := member.SetState(Transfer)
				if err != nil {
					logger.Info(err)
				}

			case End:

			}
			messageQueue.SendMessage(member)
			logger.Info(member.MemberId, ": ", member.State())
			// go func() {
			//	ch <- true
			// }()
			time.Sleep(1)
		}

	}
}

func (m *Members) CreateMembers() {

	exchange.MemberId = "elmo"
	exchange.VsCode = "EXC"
	exchange.CountryCode = "GHA"
	exchange.CurrencyCode = "GHC"
	exchange.MemberRole = "EX"
	exchange.WalletAddress = "elmo"
	exchange.CreateDate = "2019-09-10 10:10:10" // time.Now().Format("2006-01-12 15:04:05")

	if config.OperationMode == "exchange" {
		AdminLogin()
		exchange.Balance = GetElmoBalance()
		logger.Info("exchange balance :", exchange.Balance)
	}

	for i := 0; i < m.Count; i++ {
		m.Members[i] = NewMember()
		m.Members[i].MemberId = m.Prefix + strconv.Itoa(i) + "@gmail.com"
		m.Members[i].VsCode = "CVS"
		m.Members[i].CountryCode = "GHA"
		m.Members[i].CurrencyCode = "GHC"
		m.Members[i].MemberRole = "UP0"
		m.Members[i].WalletAddress = m.Prefix + strconv.Itoa(i) + "@gmail.com"
		m.Members[i].CreateDate = "2019-09-10 10:10:10" // time.Now().Format("2006-01-12 15:04:05")
	}
}

func (m *Members) Signup() {

	count := 0

	for i := 0; i < m.Count; i++ {
		count++

		go func(i int) {
			res := Signup(m.Members[i])
			if res == true {
				m.Members[i].Init()
				UserLogin(m.Members[i])
				messageQueue.SendMessage(m.Members[i])
			}

		}(i)
	}
}

func (m *Members) AdminLogin() {
	AdminLogin()
}

func (m *Members) UserLogin(user *Member) {
	UserLogin(user)
}

func (m *Members) SellElmo(user *Member) {
	var res bool
	if config.OperationMode == "blockchain" {
		res = m.Charge(user, "10000")
	} else {
		res = SellElmo(exchange, user, "10000")
	}
	if res {
		exchange.Balance -= ChargeAmount
		user.Balance += ChargeAmount
	}

}

func (m *Members) Remittance(user *Member) {
	receiver := m.Members[rand.Intn(m.Count)]
	if user.MemberId == receiver.MemberId {
		return
	}

	var res bool
	if config.OperationMode == "blockchain" {
		res = m.TransferCoin(user, receiver)
	} else {
		res = Remittance(user, receiver, "100")
	}
	if res == true {
		user.Balance -= RemittanceAmount + RemittanceFee
		receiver.Balance += RemittanceAmount
		exchange.Balance += RemittanceFee
	}
}

func (m *Members) RegisterMembers() {

	if config.OperationMode == "blockchain" {
		uuidExchange := uuid.NewV4()
		uuidStrExchange := uuidExchange.String()
		resp, err := http.PostForm("http://"+m.BoosterAddr+":8080/transaction/registeruser",
			url.Values{"txID": {uuidStrExchange}, "memberId": {exchange.MemberId}, "vsCode": {exchange.VsCode},
				"countryCode": {exchange.CountryCode}, "currencyCode": {exchange.CurrencyCode},
				"memberRole": {exchange.MemberRole}, "walletAddress": {exchange.WalletAddress},
				"txTime": {exchange.CreateDate}})
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		bytes, _ := ioutil.ReadAll(resp.Body)
		str := string(bytes)

		logger.Info("exchange creation : ", str)
		logger.Info("Creating", m.Count, "members...")

		resp, err = http.PostForm("http://"+m.BoosterAddr+":8080/transaction/issuecoin",
			url.Values{"txID": {uuidStrExchange}, "walletAddress": {exchange.WalletAddress},
				"amount": {"10000000000"}, "txTime": {"2019-09-10 10:10:10"}})
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		exchange.Balance += 10000000000
	}

	count := 0

	for i := 0; i < m.Count; i++ {
		count++

		go func(i int) {
			uuid := uuid.NewV4()
			uuidStr := uuid.String()
			resp, err := http.PostForm("http://"+m.BoosterAddr+":8080/transaction/registeruser",
				url.Values{"txID": {uuidStr}, "memberId": {m.Members[i].MemberId}, "vsCode": {m.Members[i].VsCode},
					"countryCode": {m.Members[i].CountryCode}, "currencyCode": {m.Members[i].CurrencyCode},
					"memberRole": {m.Members[i].MemberRole}, "walletAddress": {m.Members[i].WalletAddress},
					"txTime": {m.Members[i].CreateDate}})
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			bytes, _ := ioutil.ReadAll(resp.Body)
			str := string(bytes)
			logger.Info(m.Members[i].MemberId, ": member creation", str)

			if count >= m.Count {

			}

			m.Members[i].Init()
			messageQueue.SendMessage(m.Members[i])
			// ch <-true
		}(i)

	}

}

func (m *Members) TransferCoin(sender *Member, receiver *Member) bool {

	uuid := uuid.NewV4()
	uuidStr := uuid.String()
	resp, err := http.PostForm("http://"+m.BoosterAddr+":8080/transaction/transfercoin",
		url.Values{"txID": {uuidStr}, "senderWalletAddress": {sender.WalletAddress},
			"receiverWalletAddress": {receiver.WalletAddress},
			"amount":                {"100"}, "fee": {"3"}, "txFlag": {"1"}, "txTime": {"2019-09-10 10:10:10"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes)

	res := strings.Contains(str, "{\"message\":\"OK\",\"status\":200}")
	if res {
		logger.Info("Transfer Coin Success: ", sender.MemberId, "(", sender.Balance, ") to ", receiver.MemberId,
			"(", receiver.Balance, ")")

	} else {
		logger.Info("Transfer Coin Failure: ", sender.MemberId, " to ", receiver.MemberId, " ", str)
	}

	return res
}

func (m *Members) Charge(member *Member, amount string) bool {

	uuid := uuid.NewV4()
	uuidStr := uuid.String()

	resp, err := http.PostForm("http://"+m.BoosterAddr+":8080/transaction/transfercoin",
		url.Values{"txID": {uuidStr}, "senderWalletAddress": {"elmo"},
			"receiverWalletAddress": {member.WalletAddress},
			"amount":                {amount}, "fee": {"0"}, "txFlag": {"2"}, "txTime": {"2019-09-10 10:10:10"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes)

	res := strings.Contains(str, "{\"message\":\"OK\",\"status\":200}")
	if res {
		logger.Info("Charge Success: ", member.MemberId, "(", member.Balance, ")")

	} else {
		logger.Info("Charge Failure: ", member.MemberId, " ", str)
	}

	return res
}

func (m *Members) GetUserBalance(walletAddr string) string {

	url := "http://" + m.BoosterAddr + ":8080/transaction/getbalance?walletAddr=" + walletAddr
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	bodyMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(body), &bodyMap)

	if err != nil {
		panic(err)
	}

	balance := bodyMap["message"]

	// v := reflect.ValueOf(balance)
	// var floatType = reflect.TypeOf(float64(0))
	// fv := v.Convert(floatType)
	// return fv.Float()
	str := fmt.Sprintf("%v", balance)
	return str
}

func (m *Members) Report() string {
	success := 0
	failure := 0

	logger.Info("exchange balance: ", exchange.Balance, GetElmoBalance())

	var elmores string
	if exchange.Balance == GetElmoBalance() {
		elmores = "Elmo Balance correct!!!"
	} else {
		elmores = "Elmo Balance invalid!!!"
	}

	logger.Info(elmores)
	for i := 0; i < m.Count; i++ {
		userBalance := m.GetUserBalance(m.Members[i].WalletAddress)
		logger.Info(m.Members[i].WalletAddress, ": ", m.Members[i].Balance, userBalance)
		balance := strconv.FormatFloat(m.Members[i].Balance, 'f', -1, 64)
		if balance == userBalance {
			success++
		} else {
			failure++
		}
	}

	logger.Info("User Count :", m.Count)
	logger.Info("success : ", success, "/", m.Count)
	logger.Info("failure : ", failure, "/", m.Count)

	return elmores
}

type FileItem struct {
	Key      string // image_content
	FileName string // test.jpg
	Content  []byte // []byte
}

func Signup(user *Member) bool {
	signupBody := map[string]string{"birth": "19810430", "city": "accra", "cp_no": "01033150014",
		"dev_token": "asdfasdfasdfasdfasdfasdfasdfasdf", "dev_type": "AND",
		"dev_ver": "1.7.1", "dtl_addr": "no32", "eml": user.MemberId,
		"fnm": "song", "id": user.MemberId,
		"lnm": "jinun", "mtr_id": "P24180503070", "nnm": "jay",
		"pw": "111111", "sex": "M",
		"strt_addr": "silverrd", "tel_no": "01033150014"}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for k, v := range signupBody {
		bodyWriter.WriteField(k, v)
	}

	var paramFile FileItem

	fileWriter, err := bodyWriter.CreateFormFile(paramFile.Key, paramFile.FileName)
	if err != nil {
		logger.Info(err)
		return false
	}

	fileWriter.Write(paramFile.Content)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	// logger.Info(bodyBuf.String())
	resp, err := http.Post("http://"+config.ExchangeAddr+":80/api/user/signup", contentType, bodyBuf)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	str := fmt.Sprintf("%s", resp)
	res := strings.Contains(str, "200")
	if res == true {
		logger.Info(user.MemberId, "Signup Success: ", str)
	} else {
		logger.Info(user.MemberId, "Signup Fail: ", str)
	}

	return res
}

func AdminLogin() {

	form := url.Values{}
	form.Add("id", "admin01")
	form.Add("pw", "admin01")
	payload := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", "http://"+config.ExchangeAddr+":80/api/user/login", payload)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		// logger.Info("Found a cookie named:", cookie.Name, cookie.Value)
		exchange.Cookies = append(exchange.Cookies, cookie.Name+"="+cookie.Value)
	}

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) // 바이트를 문자열로
	res := strings.Contains(str, "\"status\":true,\"message\":\"Success! \"")
	if res == true {
		logger.Info("Admin Login Success: ", str)
	} else {
		logger.Info("Admin Login Fail: ", str)
	}
}

func UserLogin(user *Member) {

	form := url.Values{}
	form.Add("id", user.MemberId)
	form.Add("pw", "111111")
	payload := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", "http://"+config.ExchangeAddr+":80/api/user/login", payload)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		// logger.Info("Found a cookie named:", cookie.Name, cookie.Value)
		user.Cookies = append(user.Cookies, cookie.Name+"="+cookie.Value)
	}

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) // 바이트를 문자열로
	res := strings.Contains(str, "\"status\":true,\"message\":\"Success! \"")
	if res == true {
		logger.Info(user.MemberId, "User Login Success: ", str)
	} else {
		logger.Info(user.MemberId, "User Login Fail: ", str)
	}
}

func SellElmo(seller Member, user *Member, amount string) bool {
	form := url.Values{}
	form.Add("cash_amount", amount)
	form.Add("wallet_address", user.MemberId)
	form.Add("wallet_key", "1111")
	form.Add("type", "11")
	payload := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", "http://"+config.ExchangeAddr+":80/api/station/common/etoken/user/sell", payload)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for _, v := range seller.Cookies {
		req.Header.Add("Cookie", v)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) // 바이트를 문자열로
	res := strings.Contains(str, "\"status\":true,\"message\":\"Success! \"")
	if res == true {
		logger.Info("Sell to", user.MemberId, "Elmo Success: ", str)
	} else {
		logger.Info("Sell to", user.MemberId, "Elmo Fail: ", str)
	}

	return res
}

func Remittance(sender *Member, receiver *Member, amount string) bool {

	form := url.Values{}
	form.Add("walletAddress", receiver.MemberId)
	form.Add("amount", amount)
	form.Add("token", "TEST-TOKEN")
	payload := strings.NewReader(form.Encode())

	req, _ := http.NewRequest("POST", "http://"+config.ExchangeAddr+":80/api/wallet/remittance", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for _, v := range sender.Cookies {
		req.Header.Add("Cookie", v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) // 바이트를 문자열로
	res := strings.Contains(str, "\"status\":true,\"message\":\"Success! \"")
	if res == true {
		logger.Info(sender.MemberId, "to", receiver.MemberId, "Remittance Success: ", str)
	} else {
		logger.Info(sender.MemberId, "to", receiver.MemberId, "Remittance Fail: ", str)
	}

	return res
}

func Publish(exchange Member, amount string) {
	form := url.Values{}
	form.Add("amount", amount)
	form.Add("token", "TEST-TOKEN")
	payload := strings.NewReader(form.Encode())

	req, _ := http.NewRequest("POST", "http://"+config.ExchangeAddr+":80/api/admin/publish/elmo", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for _, v := range exchange.Cookies {
		req.Header.Add("Cookie", v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) // 바이트를 문자열로
	res := strings.Contains(str, "\"status\":true,\"message\":\"Success! \"")
	if res == true {
		logger.Info("Publish Success: ", str)
	} else {
		logger.Info("Publish Fail: ", str)
	}
}

func GetElmoBalance() float64 {

	url := "http://" + config.ExchangeAddr + ":80/api/station/common/etoken/balance?wallet-key=20"

	req, _ := http.NewRequest("GET", url, nil)
	for _, v := range exchange.Cookies {
		req.Header.Add("Cookie", v)
	}
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	bodyMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(body), &bodyMap)
	if err != nil {
		panic(err)
	}

	subBody := bodyMap["object"].(map[string]interface{})
	balance := subBody["elmo_blc"]

	v := reflect.ValueOf(balance)
	var floatType = reflect.TypeOf(float64(0))
	fv := v.Convert(floatType)

	return fv.Float()
}
