# wechatpay

微信支付SDK for Go！包括微信商户支付的全部功能！集成简单！





## 安装

`go get -u github.com/liyoung1992/wechatpay`


## 帮助

如果在集成过程中遇到问题，请联系：liyoung_1992@163.com

## 目前实现的接口

- 扫码支付（NATIVE ）

- H5支付 （MWEB）

- 公众号支付 （JSAPI ）

- APP支付 （APP）

- 小程序支付 (JSAPI)

- 退款

- 退款查询



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



### 统一支付(其他支付改对应的tradetype即可)



```go

	var pay_data wechatpay.UnitOrder
	pay_data.NotifyUrl ="your notify_url"
	pay_data.TradeType = "NATIVE"
	pay_data.Body = payweb.Subject
	pay_data.SpbillCreateIp =  c.ClientIP()
	pay_data.TotalFee = 1
	pay_data.OutTradeNo = payweb.OrderId
	result ,err:= wechat_client.Pay(pay_data)

```

APP支付、公众号支付、小程序支付都是先返回：预支付交易单，然后用预支付交易码在进行支付操作；
小程序和公众号支付都需要获取用户的openid。具体的获取方式:[获取openid](https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_4)

### 发起退款



```go


	var refund_data wechatpay.OrderRefund
	refund_data.TotalFee = 1 
	refund_data.OutTradeNo = order_id
	refund_data.OutRefundNo = "r" + order_id
	refund_data.RefundFee = 1 
	result, err := wechat_client.Refund(refund_data)

```



### 退款查询



```go


	var refund_status wechatpay.OrderRefundQuery
	refund_status.OutTradeNo = order_id
	result, err := wechat_client.RefundQuery(refund_status)


```



## 异步通知


支付结果的异步通知，这里使用gin框架的例子，其他框架原理一样（也可以不用框架），接受到微信通知，返回给微信数据即可！
修复

```go

//微信扫码回调地址
func (pay WechatPayController) WechatWebNotifyUrl(c *gin.Context) {
	log.Info("wechat notify start")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err, "read http body failed！error msg:"+err.Error())
	}
	log.Info("wechat pay notify body :" + string(body))

	var wx_notify_req wechatpay.PayNotifyResult

	err = xml.Unmarshal(body, &wx_notify_req)
	if err != nil {
		log.Error(err, "read http body xml failed! err :"+err.Error())
	}
	mv, err := mxj.NewMapXml(body)
	if err != nil {
		log.Error(err,err.Error())
	}
	
	//进行签名校验
	if wechat_client.VerifySign(mv["xml"].(map[string]interface{}), mv["xml"].(map[string]interface{})["sign"].(string)) {
		record, err := json.Marshal(wx_notify_req)
		if err != nil {
			log.Error(err, "wechat pay marshal err :"+err.Error())
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

## 其他支付方式（gosdk）



- 支付宝 [支付宝sdk](https://github.com/smartwalle/alipay)

- stripe [stripe国际支付](https://github.com/stripe/stripe-go)

- paypal [paypal国际支付](https://github.com/logpacker/PayPal-Go-SDK)


以上几种支付sdk，都亲测可用！且已经开发上线。有需要的可以直接联系我！



## License



This project is licensed under the MIT License.

