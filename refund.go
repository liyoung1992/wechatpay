package wechatpay

import (
	"bytes"

	"encoding/xml"

	"glink/AIYShopWeb/shared/log"

	"io/ioutil"

	"net/http"

	"strconv"
	"time"

	"strings"
)

//退款
func (this *WechatPay) Refund(param OrderRefund) bool {

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
	param.Sign = getSign(m, this.ApiKey)

	bytes_req, err := xml.Marshal(param)
	if err != nil {
		log.Error(err, "xml marshal failed")
		return false
	}

	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "Refund", "xml", -1)
	bytes_req = []byte(str_req)

	//发送unified order请求.
	req, err := http.NewRequest("POST", REFUND_URL, bytes.NewReader(bytes_req))
	if err != nil {
		log.Error(err, "new http request failed,err :"+err.Error())
		return false
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	w_req := http.Client{
		Transport: withCertBytes(this.ApiclientCert, this.ApiclientKey),
	}

	resp, _err := w_req.Do(req)
	if _err != nil {
		log.Error(err, "http request failed! err :"+_err.Error())
		return false
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var refund_resp OrderRefundResult

	_err = xml.Unmarshal(body, &refund_resp)
	if _err != nil {
		log.Error(err, "http request failed! err :"+_err.Error())
		return false
	}
	if refund_resp.ResultCode == "SUCCESS" {
		return true
	} else {
		return false
	}

}

//退款查询
func (this *WechatPay) RefundQuery(order_id string) bool {

	var refund_status OrderRefundQuery
	refund_status.AppId = this.AppId
	refund_status.MchId = this.MchId
	refund_status.NonceStr = order_id + strconv.FormatInt(time.Now().Unix(), 10)
	refund_status.OutTradeNo = order_id

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = refund_status.AppId
	m["mch_id"] = refund_status.MchId
	m["out_trade_no"] = refund_status.OutTradeNo
	m["nonce_str"] = refund_status.NonceStr
	refund_status.Sign = getSign(m, this.ApiKey)

	bytes_req, err := xml.Marshal(refund_status)
	if err != nil {
		log.Error(err, "xml marshal failed,err:"+err.Error())
		return false
	}

	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "RefundQuery", "xml", -1)
	bytes_req = []byte(str_req)

	//发送unified order请求.
	req, err := http.NewRequest("POST", REFUND_QUERY_URL, bytes.NewReader(bytes_req))
	if err != nil {
		log.Error(err, "new http request failed,err :"+err.Error())
		return false
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	w_req := http.Client{}
	resp, _err := w_req.Do(req)
	if _err != nil {
		log.Error(err, "http request failed! err :"+_err.Error())
		return false
	}
	var refund_resp OrderRefundQueryResult
	body, _ := ioutil.ReadAll(resp.Body)

	_err = xml.Unmarshal(body, &refund_resp)
	if _err != nil {
		log.Error(err, "http request failed! err :"+_err.Error())
		return false
	}
	if refund_resp.RefundStatus_0 == "SUCCESS" {
		return true
	} else {
		return false
	}
}
