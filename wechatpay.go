package wechatpay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"

	"encoding/xml"
	"fmt"
	"glink/AIYShopWeb/shared/log"
	// "glink/AIYShopWeb/shared/beegorm"
	"crypto/tls"

	"io/ioutil"
	"net/http"
	"math/rand"
	"time"
	"sort"
	// "strconv"
	"strings"
	"errors"

)

type WechatPay struct {
	AppId string 
	MchId string 
	ApiKey string 
	ApiclientCert []byte 
	ApiclientKey []byte 
}

func New(appId,mchId,apiKey string,apiclient_cert,apiclient_key []byte) (client *WechatPay) {
	client = &WechatPay{}
	client.AppId = appId
	client.MchId = mchId
	client.ApiKey = apiKey
	client.ApiclientCert = apiclient_cert
	client.ApiclientKey = apiclient_key
	return client	
}
//wxpay计算签名的函数
func getSign(mReq map[string]interface{}, key string) (sign string) {
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	var signStrings string
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	if key != "" {
		signStrings = signStrings + "key=" + key
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

//微信支付签名验证函数
func (this *WechatPay) verifySign(needVerifyM map[string]interface{}, sign string) bool {
	signCalc := getSign(needVerifyM, this.ApiKey)
	if sign == signCalc {
		log.Info("wechat verify success!")
		return true
	}
	log.Info("wechat vertify failed!")
	return false
}

//
func withCertBytes(cert,key []byte) *http.Transport {
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil
	}
	conf := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}
	trans := &http.Transport{
		TLSClientConfig: conf,
	}
	return trans
}

func randomNonceStr() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r :=  rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 32; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//统一下单
func (this *WechatPay) WechatOrderPay(param UnitOrder) (string, error){
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
	param.Sign = getSign(m, this.ApiKey)

	bytes_req, err := xml.Marshal(param)
	if err != nil {
		return "",err
	}
	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "UnitOrder", "xml", -1)

	req, err := http.NewRequest("POST", UNIT_ORDER_URL, bytes.NewReader([]byte(str_req)))
	if err != nil {
		return "",err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	if param.TradeType == "MWEB" {
	   req.Header.Set("Referer",param.Referer)
	}

	w_req := http.Client{}
	resp, err := w_req.Do(req)
	if err != nil {
		return "",err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var pay_result UnifyOrderResult
	err = xml.Unmarshal(body, &pay_result)
	if pay_result.ReturnCode == "SUCCESS" && pay_result.ResultCode == "SUCCESS" {
		if pay_result.TradeType == "MWEB" {
			return pay_result.MwebUrl,nil
		}else {
			return pay_result.CodeUrl,nil
		}
	}
	return "",errors.New(pay_result.ReturnMsg)
}
