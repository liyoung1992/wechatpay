package wechatpay

import (
	"fmt"
	"testing"
)
func TestWechat_Refund(t *testing.T) {
	var refund_data OrderRefund

	refund_data.TotalFee = 1
	refund_data.OutTradeNo = "1234567"
	refund_data.OutRefundNo = "r122121"
	refund_data.RefundFee = 1
	fmt.Println(wechat_client.Refund(refund_data))
}