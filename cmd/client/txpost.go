package client

/*

func Signup(user *model.Member) {
	signupBody := map[string]string{"birth": "19810430", "city": "accra", "cp_no": "01033150014",
		"dev_token": "asdfasdfasdfasdfasdfasdfasdfasdf", "dev_type": "AND",
		"dev_ver": "1.7.1", "dtl_addr": "no32","eml": user.MemberId,
		"fnm": "song","id": user.MemberId,
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
		return
	}

	fileWriter.Write(paramFile.Content)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	//logger.Info(bodyBuf.String())
	resp, err := http.Post("https://www.elmo.africa/api/user/signup", contentType, bodyBuf)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	res := strings.Contains(fmt.Sprintf(resp), "200 OK")
	if res == true {
		logger.Info(user.MemberId, "Signup Success: ", str)
	} else {
		logger.Info(user.MemberId, "Signup Fail: ", str)
	}
}

func AdminLogin() {

	form := url.Values{}
	form.Add("id", "admin01")
	form.Add("pw", "admin01")
	payload := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", "https://www.elmo.africa/api/user/login", payload)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)

	for _, cookie := range resp.Cookies() {
		fmt.Println("Found a cookie named:", cookie.Name, cookie.Value)
		cookies = append(cookies, cookie.Name+"="+cookie.Value)
	}

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	res := strings.Contains(str, "\"status\":true,\"message\":\"Success! \"")
	if res == true {
		logger.Info("Admin Login Success: ", str)
	} else {
		logger.Info("Admin Login Fail: ", str)
	}
}


func UserLogin(user *model.Member) {

	form := url.Values{}
	form.Add("id", user.MemberId)
	form.Add("pw", "111111")
	payload := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", "https://www.elmo.africa/api/user/login", payload)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)

	for _, cookie := range resp.Cookies() {
		fmt.Println("Found a cookie named:", cookie.Name, cookie.Value)
		cookies = append(user.Cookies, cookie.Name+"="+cookie.Value)
	}

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로
	res := strings.Contains(str, "\"status\":true,\"message\":\"Success! \"")
	if res == true {
		logger.Info(user.MemberId, "User Login Success: ", str)
	} else {
		logger.Info(user.MemberId, "User Login Fail: ", str)
	}
}


func SellElmo(seller model.Member, user *model.Member, amount string) bool {
	form := url.Values{}
	form.Add("cash_amount", amount)
	form.Add("wallet_address", user.MemberId)
	form.Add("wallet_key", "1111")
	form.Add("type", "11")
	payload := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", "https://www.elmo.africa/api/station/common/etoken/user/sell", payload)
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
	str := string(bytes) //바이트를 문자열로
	res := strings.Contains(str, "\"status\":true,\"message\":\"Success! \"")
	if res == true {
		logger.Info("Sell to", user.MemberId, "Elmo Success: ", str)
	} else {
		logger.Info("Sell to", user.MemberId, "Elmo Fail: ", str)
	}

	return res
}

func Remittance(sender *model.Member, receiver *model.Member, amount string) bool {

	form := url.Values{}
	form.Add("walletAddress", receiver.MemberId)
	form.Add("amount", amount)
	form.Add("token", "TEST-TOKEN")
	payload := strings.NewReader(form.Encode())

	req, _ := http.NewRequest("POST", "https://www.elmo.africa/api/wallet/remittance", payload)
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
	str := string(bytes) //바이트를 문자열로
	res := strings.Contains(str, "\"status\":true,\"message\":\"Success! \"")
	if res == true {
		logger.Info(sender.MemberId, "to", receiver.MemberId, "Remittance Success: ", str)
	} else {
		logger.Info(sender.MemberId, "to", receiver.MemberId, "Remittance Fail: ", str)
	}

	return res
}


func Publish(exchange model.Member, amount string) {
	form := url.Values{}
	form.Add("amount", amount)
	form.Add("token", "TEST-TOKEN")
	payload := strings.NewReader(form.Encode())

	req, _ := http.NewRequest("POST", "https://www.elmo.africa/api/admin/publish/elmo", payload)
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
	str := string(bytes) //바이트를 문자열로
	res := strings.Contains(str, "\"status\":true,\"message\":\"Success! \"")
	if res == true {
		logger.Info("Publish Success: ", str)
	} else {
		logger.Info("Publish Fail: ", str)
	}
} */
