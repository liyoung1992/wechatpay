# wechatpay
微信支付SDK for Go！


## 安装

`go get -u github.com/liyoung1992/wechatpay`

## 帮助
如果在集成过程中遇到问题，请联系：liyoung_1992@163.com

## 目前实现的接口

- 扫码支付（NATIVE ）

- H5支付 （MWEB）

- 公众号支付 （JSAPI ）

- APP支付 （APP）

- 退款

- 退款查询

## TODO 
- 小程序支付
- 单元测试


## 集成方式
强烈建议开发前仔细阅读[微信支付官方文档](https://pay.weixin.qq.com/wiki/doc/api/index.html)


### 创建支付

```go

	wechat_cert, err := ioutil.ReadFile("config/wechat/apiclient_cert.pem")
	if err != nil {
		panic(err)
	}
	wechat_key, err := ioutil.ReadFile("config/wechat/apiclient_key.pem")

	wechat_client = wechatpay.New(os.Getenv("WECHAT_APPID"),os.Getenv("WECHAT_MCHID"),
	os.Getenv("WECHAT_APIKEY"),wechat_key,wechat_cert)

	if err != nil {
		panic(err)
	}

```

### 发起扫码支付(其他支付改对应的tradetype即可)

```go

	var pay_data wechatpay.UnitOrder
	pay_data.NotifyUrl = os.Getenv("WECHAT_NOTIFY_URL")
	pay_data.TradeType = "NATIVE"
	pay_data.Body = payweb.Subject
	pay_data.SpbillCreateIp =  c.ClientIP()

	pay_data.TotalFee = 1
	pay_data.OutTradeNo = payweb.OrderId
	result ,err:= wechat_client.Pay(pay_data)
    ......

```
APP支付和公众号支付都是先返回：预支付交易单，然后用预支付交易码在进行支付操作

## 异步通知

```go

//微信扫码回调地址(gin框架)
func (this *WechatPay) PayNotifyUrl(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err, "read notify body error")
	}

	var wx_notify_req PayNotifyResult
	err = xml.Unmarshal(body, &wx_notify_req)
	if err != nil {
		log.Error(err, "read http body xml failed! err :"+err.Error())
	}
	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)

	reqMap["return_code"] = wx_notify_req.ReturnCode
	reqMap["return_msg"] = wx_notify_req.ReturnMsg
	reqMap["appid"] = wx_notify_req.AppId
	reqMap["mch_id"] = wx_notify_req.MchId
	reqMap["nonce_str"] = wx_notify_req.NonceStr
	reqMap["result_code"] = wx_notify_req.ResultCode
	reqMap["openid"] = wx_notify_req.OpenId
	reqMap["is_subscribe"] = wx_notify_req.IsSubscribe
	reqMap["trade_type"] = wx_notify_req.TradeType
	reqMap["bank_type"] = wx_notify_req.BankType
	reqMap["total_fee"] = wx_notify_req.TotalFee
	reqMap["fee_type"] = wx_notify_req.FeeType
	reqMap["cash_fee"] = wx_notify_req.CashFee
	reqMap["cash_fee_type"] = wx_notify_req.CashFeeType
	reqMap["transaction_id"] = wx_notify_req.TransactionId
	reqMap["out_trade_no"] = wx_notify_req.OutTradeNo
	reqMap["attach"] = wx_notify_req.Attach
	reqMap["time_end"] = wx_notify_req.TimeEnd

	//进行签名校验
	if this.VerifySign(reqMap, wx_notify_req.Sign) {
		record, err := json.Marshal(wx_notify_req)
		if err != nil {
			log.Error(err, "wechat pay marshal err :"+err.Error())
		}
		//TODO 加入你的代码，处理返回值
		fmt.Println(string(record))
		// err = wechat_pay_recoed_producer.Publish("wechat_pay", record)
		if err != nil {
			log.Error(err, "wechat publish record err:"+err.Error())
		}
		c.XML(http.StatusOK, gin.H{
			"return_code": "SUCCESS",
			"return_msg":  "OK",
		})
	} else {
		c.XML(http.StatusOK, gin.H{
			"return_code": "FAIL",
			"return_msg":  "failed to verify sign, please retry!",
		})
	}
	return
}

```

## License

This project is licensed under the MIT License.
