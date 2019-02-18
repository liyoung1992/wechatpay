package wechatpay

import (
	"bytes"
	"encoding/xml"
	"errors"
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
	if param.FeeType != "" {
		m["fee_type"] = param.FeeType
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
	} else if pay_result.ReturnCode != "SUCCESS" {
		return nil, errors.New(pay_result.ReturnMsg)
	}
	return &pay_result, nil
}
