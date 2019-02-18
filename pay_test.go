package wechatpay


import (
	"os"
	"fmt"
	"testing"
)
var (
	wechat_cert = "111111111111232121321311"
	wechat_key = "12123222222222223232332323"
	wechat_app_id = "102801212"
	wechat_mch_id = "232312123"
	wechat_api_key = "121212"
)
var wechat_client *WechatPay
func TestMain(m *testing.M) {
	wechat_client = New(wechat_app_id, wechat_mch_id,
		wechat_api_key, []byte(wechat_cert), []byte(wechat_key))
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestWechat_Pay(t *testing.T) {
	var pay_data UnitOrder
	pay_data.NotifyUrl = "http://47.98.87.189"
	pay_data.TradeType = "NATIVE"
	pay_data.Body = "测试支付"
	pay_data.SpbillCreateIp = "47.98.87.189"
	pay_data.TotalFee = 1
	pay_data.OutTradeNo = "123456789"

	fmt.Println(wechat_client.Pay(pay_data))
}