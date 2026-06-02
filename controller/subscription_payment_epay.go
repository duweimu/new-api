package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Calcium-Ion/go-epay/epay"
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type SubscriptionEpayPayRequest struct {
	PlanId        int    `json:"plan_id"`
	PaymentMethod string `json:"payment_method"`
}

type SubscriptionEpayAmountRequest struct {
	PlanId int `json:"plan_id"`
}

type subscriptionEpayPurchaseArgs struct {
	Type           string
	ServiceTradeNo string
	Name           string
	Money          string
	Device         epay.DeviceType
	NotifyUrl      *url.URL
	ReturnUrl      *url.URL
}

var subscriptionEpayPurchase = func(baseURL string, args subscriptionEpayPurchaseArgs) (string, map[string]string, error) {
	client, err := epay.NewClient(&epay.Config{
		PartnerID: operation_setting.EpayId,
		Key:       operation_setting.EpayKey,
	}, baseURL)
	if err != nil {
		return "", nil, err
	}
	return client.Purchase(&epay.PurchaseArgs{
		Type:           args.Type,
		ServiceTradeNo: args.ServiceTradeNo,
		Name:           args.Name,
		Money:          args.Money,
		Device:         args.Device,
		NotifyUrl:      args.NotifyUrl,
		ReturnUrl:      args.ReturnUrl,
	})
}

func getSubscriptionPlanForEpay(planID int) (*model.SubscriptionPlan, error) {
	if planID <= 0 {
		return nil, errors.New("invalid subscription plan")
	}

	plan, err := model.GetSubscriptionPlanById(planID)
	if err != nil {
		return nil, err
	}
	if !plan.Enabled {
		return nil, errors.New("subscription plan is disabled")
	}
	if plan.PriceAmount < 0.01 {
		return nil, errors.New("subscription plan amount is too low")
	}
	return plan, nil
}

func calculateSubscriptionEpayMoney(plan *model.SubscriptionPlan) (float64, error) {
	if plan == nil {
		return 0, errors.New("subscription plan is required")
	}
	if plan.PriceAmount < 0.01 {
		return 0, errors.New("subscription plan amount is too low")
	}

	payMoney := decimal.NewFromFloat(plan.PriceAmount).
		Mul(decimal.NewFromFloat(operation_setting.Price)).
		Round(2)
	if payMoney.LessThan(decimal.RequireFromString("0.01")) {
		return 0, errors.New("payment amount is too low")
	}
	return payMoney.InexactFloat64(), nil
}

func SubscriptionRequestEpayAmount(c *gin.Context) {
	if !requirePaymentCompliance(c) {
		return
	}

	var req SubscriptionEpayAmountRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.PlanId <= 0 {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid parameters"})
		return
	}

	plan, err := getSubscriptionPlanForEpay(req.PlanId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": err.Error()})
		return
	}
	if GetEpayClient() == nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "epay is not configured"})
		return
	}

	payMoney, err := calculateSubscriptionEpayMoney(plan)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": strconv.FormatFloat(payMoney, 'f', 2, 64)})
}

func SubscriptionRequestEpay(c *gin.Context) {
	if !requirePaymentCompliance(c) {
		return
	}

	var req SubscriptionEpayPayRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.PlanId <= 0 {
		common.ApiErrorMsg(c, "invalid parameters")
		return
	}

	plan, err := getSubscriptionPlanForEpay(req.PlanId)
	if err != nil {
		common.ApiErrorMsg(c, err.Error())
		return
	}
	if !operation_setting.ContainsPayMethod(req.PaymentMethod) {
		common.ApiErrorMsg(c, "payment method does not exist")
		return
	}

	userId := c.GetInt("id")
	if plan.MaxPurchasePerUser > 0 {
		count, err := model.CountUserSubscriptionsByPlan(userId, plan.Id)
		if err != nil {
			common.ApiError(c, err)
			return
		}
		if count >= int64(plan.MaxPurchasePerUser) {
			common.ApiErrorMsg(c, "subscription purchase limit reached")
			return
		}
	}

	callBackAddress := service.GetCallbackAddress()
	returnUrl, err := url.Parse(callBackAddress + "/api/subscription/epay/return")
	if err != nil {
		common.ApiErrorMsg(c, "invalid payment callback address")
		return
	}
	notifyUrl, err := url.Parse(callBackAddress + "/api/subscription/epay/notify")
	if err != nil {
		common.ApiErrorMsg(c, "invalid payment callback address")
		return
	}

	tradeNo := fmt.Sprintf("%s%d", common.GetRandomString(6), time.Now().Unix())
	tradeNo = fmt.Sprintf("SUBUSR%dNO%s", userId, tradeNo)

	if GetEpayClient() == nil {
		common.ApiErrorMsg(c, "epay is not configured")
		return
	}

	payMoney, err := calculateSubscriptionEpayMoney(plan)
	if err != nil {
		common.ApiErrorMsg(c, err.Error())
		return
	}
	payMoneyText := strconv.FormatFloat(payMoney, 'f', 2, 64)

	order := &model.SubscriptionOrder{
		UserId:          userId,
		PlanId:          plan.Id,
		Money:           payMoney,
		TradeNo:         tradeNo,
		PaymentMethod:   req.PaymentMethod,
		PaymentProvider: model.PaymentProviderEpay,
		CreateTime:      time.Now().Unix(),
		Status:          common.TopUpStatusPending,
	}
	if err := order.Insert(); err != nil {
		common.ApiErrorMsg(c, "failed to create subscription order")
		return
	}

	uri, params, err := subscriptionEpayPurchase(operation_setting.PayAddress, subscriptionEpayPurchaseArgs{
		Type:           req.PaymentMethod,
		ServiceTradeNo: tradeNo,
		Name:           fmt.Sprintf("SUB:%s", plan.Title),
		Money:          payMoneyText,
		Device:         epay.PC,
		NotifyUrl:      notifyUrl,
		ReturnUrl:      returnUrl,
	})
	if err != nil {
		_ = model.ExpireSubscriptionOrder(tradeNo, model.PaymentProviderEpay)
		common.ApiErrorMsg(c, "failed to initiate payment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": params, "url": uri})
}

func SubscriptionEpayNotify(c *gin.Context) {
	var params map[string]string

	if c.Request.Method == "POST" {
		if err := c.Request.ParseForm(); err != nil {
			_, _ = c.Writer.Write([]byte("fail"))
			return
		}
		params = lo.Reduce(lo.Keys(c.Request.PostForm), func(r map[string]string, t string, i int) map[string]string {
			r[t] = c.Request.PostForm.Get(t)
			return r
		}, map[string]string{})
	} else {
		params = lo.Reduce(lo.Keys(c.Request.URL.Query()), func(r map[string]string, t string, i int) map[string]string {
			r[t] = c.Request.URL.Query().Get(t)
			return r
		}, map[string]string{})
	}

	if len(params) == 0 {
		_, _ = c.Writer.Write([]byte("fail"))
		return
	}

	client := GetEpayClient()
	if client == nil {
		_, _ = c.Writer.Write([]byte("fail"))
		return
	}
	verifyInfo, err := client.Verify(params)
	if err != nil || !verifyInfo.VerifyStatus {
		_, _ = c.Writer.Write([]byte("fail"))
		return
	}

	if verifyInfo.TradeStatus != epay.StatusTradeSuccess {
		_, _ = c.Writer.Write([]byte("fail"))
		return
	}

	LockOrder(verifyInfo.ServiceTradeNo)
	defer UnlockOrder(verifyInfo.ServiceTradeNo)

	if err := model.CompleteSubscriptionOrder(verifyInfo.ServiceTradeNo, common.GetJsonString(verifyInfo), model.PaymentProviderEpay, verifyInfo.Type); err != nil {
		_, _ = c.Writer.Write([]byte("fail"))
		return
	}

	_, _ = c.Writer.Write([]byte("success"))
}

// SubscriptionEpayReturn handles browser return after payment.
// It verifies the payload and completes the order, then redirects to console.
func SubscriptionEpayReturn(c *gin.Context) {
	var params map[string]string

	if c.Request.Method == "POST" {
		if err := c.Request.ParseForm(); err != nil {
			c.Redirect(http.StatusFound, paymentReturnPath("/console/topup?pay=fail"))
			return
		}
		params = lo.Reduce(lo.Keys(c.Request.PostForm), func(r map[string]string, t string, i int) map[string]string {
			r[t] = c.Request.PostForm.Get(t)
			return r
		}, map[string]string{})
	} else {
		params = lo.Reduce(lo.Keys(c.Request.URL.Query()), func(r map[string]string, t string, i int) map[string]string {
			r[t] = c.Request.URL.Query().Get(t)
			return r
		}, map[string]string{})
	}

	if len(params) == 0 {
		c.Redirect(http.StatusFound, paymentReturnPath("/console/topup?pay=fail"))
		return
	}

	client := GetEpayClient()
	if client == nil {
		c.Redirect(http.StatusFound, paymentReturnPath("/console/topup?pay=fail"))
		return
	}
	verifyInfo, err := client.Verify(params)
	if err != nil || !verifyInfo.VerifyStatus {
		c.Redirect(http.StatusFound, paymentReturnPath("/console/topup?pay=fail"))
		return
	}
	if verifyInfo.TradeStatus == epay.StatusTradeSuccess {
		LockOrder(verifyInfo.ServiceTradeNo)
		defer UnlockOrder(verifyInfo.ServiceTradeNo)
		if err := model.CompleteSubscriptionOrder(verifyInfo.ServiceTradeNo, common.GetJsonString(verifyInfo), model.PaymentProviderEpay, verifyInfo.Type); err != nil {
			c.Redirect(http.StatusFound, paymentReturnPath("/console/topup?pay=fail"))
			return
		}
		c.Redirect(http.StatusFound, paymentReturnPath("/console/topup?pay=success"))
		return
	}
	c.Redirect(http.StatusFound, paymentReturnPath("/console/topup?pay=pending"))
}
