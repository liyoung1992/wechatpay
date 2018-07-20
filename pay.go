package wechatpay

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

//统一下单
func (this *WechatPay) Pay(param UnitOrder) (*UnifyOrderResult, error) {
	param.AppId = this.AppId
	param.MchId = this.MchId
	param.NonceStr = randomNonceStr()

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = param.AppId
	m["body"] = param.Body
	m["mch_id"] = param.MchId
	m["notify_url"] = param.NotifyUrl
	m["trade_type"] = param.TradeType
	m["spbill_create_ip"] = param.SpbillCreateIp
	m["total_fee"] = param.TotalFee
	m["out_trade_no"] = param.OutTradeNo
	m["nonce_str"] = param.NonceStr
	if param.TradeType == "MWEB" {
		m["scene_info"] = param.SceneInfo
	}
	if param.TradeType == "JSAPI" {
		m["openid"] = param.Openid
	}
	param.Sign = GetSign(m, this.ApiKey)

	bytes_req, err := xml.Marshal(param)
	if err != nil {
		return nil, err
	}
	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "UnitOrder", "xml", -1)
	req, err := http.NewRequest("POST", UNIT_ORDER_URL, bytes.NewReader([]byte(str_req)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	if param.TradeType == "MWEB" {
		req.Header.Set("Referer", param.Referer)
	}

	w_req := http.Client{}
	resp, err := w_req.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var pay_result UnifyOrderResult
	err = xml.Unmarshal(body, &pay_result)
	if err != nil {
		return nil, err
	}
	return &pay_result, nil
}

//
//微信扫码回调地址(gin框架)
// func (this *WechatPay) PayNotifyUrl(c *gin.Context) {

// 	body, err := ioutil.ReadAll(c.Request.Body)
// 	if err != nil {
// 		log.Error(err, "read notify body error")
// 	}

// 	var wx_notify_req PayNotifyResult
// 	err = xml.Unmarshal(body, &wx_notify_req)
// 	if err != nil {
// 		log.Error(err, "read http body xml failed! err :"+err.Error())
// 	}
// 	var reqMap map[string]interface{}
// 	reqMap = make(map[string]interface{}, 0)

// 	reqMap["return_code"] = wx_notify_req.ReturnCode
// 	reqMap["return_msg"] = wx_notify_req.ReturnMsg
// 	reqMap["appid"] = wx_notify_req.AppId
// 	reqMap["mch_id"] = wx_notify_req.MchId
// 	reqMap["nonce_str"] = wx_notify_req.NonceStr
// 	reqMap["result_code"] = wx_notify_req.ResultCode
// 	reqMap["openid"] = wx_notify_req.OpenId
// 	reqMap["is_subscribe"] = wx_notify_req.IsSubscribe
// 	reqMap["trade_type"] = wx_notify_req.TradeType
// 	reqMap["bank_type"] = wx_notify_req.BankType
// 	reqMap["total_fee"] = wx_notify_req.TotalFee
// 	reqMap["fee_type"] = wx_notify_req.FeeType
// 	reqMap["cash_fee"] = wx_notify_req.CashFee
// 	reqMap["cash_fee_type"] = wx_notify_req.CashFeeType
// 	reqMap["transaction_id"] = wx_notify_req.TransactionId
// 	reqMap["out_trade_no"] = wx_notify_req.OutTradeNo
// 	reqMap["attach"] = wx_notify_req.Attach
// 	reqMap["time_end"] = wx_notify_req.TimeEnd

// 	//进行签名校验
// 	if this.VerifySign(reqMap, wx_notify_req.Sign) {
// 		record, err := json.Marshal(wx_notify_req)
// 		if err != nil {
// 			log.Error(err, "wechat pay marshal err :"+err.Error())
// 		}
// 		//TODO 加入你的代码，处理返回值
// 		fmt.Println(string(record))
// 		// err = wechat_pay_recoed_producer.Publish("wechat_pay", record)
// 		if err != nil {
// 			log.Error(err, "wechat publish record err:"+err.Error())
// 		}
// 		c.XML(http.StatusOK, gin.H{
// 			"return_code": "SUCCESS",
// 			"return_msg":  "OK",
// 		})
// 	} else {
// 		c.XML(http.StatusOK, gin.H{
// 			"return_code": "FAIL",
// 			"return_msg":  "failed to verify sign, please retry!",
// 		})
// 	}
// 	return
// }
