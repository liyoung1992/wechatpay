package wechatpay

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"glink/AIYShopWeb/shared/log"
	"crypto/tls"
	"math/rand"
	"net/http"
	"sort"
	"time"
	"strings"
)

type WechatPay struct {
	AppId         string
	MchId         string
	ApiKey        string
	ApiclientCert []byte
	ApiclientKey  []byte
}

func New(appId, mchId, apiKey string, apiclient_cert, apiclient_key []byte) (client *WechatPay) {
	client = &WechatPay{}
	client.AppId = appId
	client.MchId = mchId
	client.ApiKey = apiKey
	client.ApiclientCert = apiclient_cert
	client.ApiclientKey = apiclient_key
	return client
}

//wxpay计算签名的函数
func GetSign(mReq map[string]interface{}, key string) (sign string) {

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
	signCalc := GetSign(needVerifyM, this.ApiKey)
	if sign == signCalc {
		log.Info("wechat verify success!")
		return true
	}
	log.Info("wechat vertify failed!")
	return false
}

//
func withCertBytes(cert, key []byte) *http.Transport {
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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 32; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
