package wechatpay

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

//退款
func (this *WechatPay) Refund(param OrderRefund) (*OrderRefundResult, error) {

	param.AppId = this.AppId
	param.MchId = this.MchId
	param.NonceStr = randomNonceStr()

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = param.AppId
	m["mch_id"] = param.MchId
	m["total_fee"] = param.TotalFee
	m["out_trade_no"] = param.OutTradeNo
	m["nonce_str"] = param.NonceStr
	m["refund_fee"] = param.RefundFee
	m["out_refund_no"] = param.OutRefundNo
	param.Sign = GetSign(m, this.ApiKey)

	bytes_req, err := xml.Marshal(param)
	if err != nil {
		return nil, err
	}

	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "Refund", "xml", -1)
	bytes_req = []byte(str_req)

	//发送unified order请求.
	req, err := http.NewRequest("POST", REFUND_URL, bytes.NewReader(bytes_req))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	w_req := http.Client{
		Transport: WithCertBytes(this.ApiclientCert, this.ApiclientKey),
	}

	resp, _err := w_req.Do(req)
	if _err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var refund_resp OrderRefundResult

	_err = xml.Unmarshal(body, &refund_resp)
	if _err != nil {
		return nil, err
	}
	return &refund_resp, nil
}

//退款查询
func (this *WechatPay) RefundQuery(refund_status OrderRefundQuery) (*OrderRefundQueryResult, error) {

	refund_status.AppId = this.AppId
	refund_status.MchId = this.MchId
	refund_status.NonceStr = randomNonceStr()

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = refund_status.AppId
	m["mch_id"] = refund_status.MchId
	m["out_trade_no"] = refund_status.OutTradeNo
	m["nonce_str"] = refund_status.NonceStr
	refund_status.Sign = GetSign(m, this.ApiKey)

	bytes_req, err := xml.Marshal(refund_status)
	if err != nil {
		return nil, err
	}

	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "RefundQuery", "xml", -1)
	bytes_req = []byte(str_req)

	req, err := http.NewRequest("POST", REFUND_QUERY_URL, bytes.NewReader(bytes_req))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	w_req := http.Client{}
	resp, _err := w_req.Do(req)
	if _err != nil {
		return nil, err
	}
	var refund_resp OrderRefundQueryResult
	body, _ := ioutil.ReadAll(resp.Body)

	_err = xml.Unmarshal(body, &refund_resp)
	if _err != nil {
		return nil, err
	}
	return &refund_resp, nil
}
