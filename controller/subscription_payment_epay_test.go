package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type subscriptionAmountTestResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

type subscriptionPayTestResponse struct {
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
	URL     string            `json:"url"`
}

func setupSubscriptionPaymentControllerTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	gin.SetMode(gin.TestMode)

	originalDB := model.DB
	originalLogDB := model.LOG_DB
	originalUsingSQLite := common.UsingSQLite
	originalUsingMySQL := common.UsingMySQL
	originalUsingPostgreSQL := common.UsingPostgreSQL
	originalRedisEnabled := common.RedisEnabled

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	model.DB = db
	model.LOG_DB = db
	common.UsingSQLite = true
	common.UsingMySQL = false
	common.UsingPostgreSQL = false
	common.RedisEnabled = false

	require.NoError(t, db.AutoMigrate(
		&model.User{},
		&model.SubscriptionPlan{},
		&model.SubscriptionOrder{},
		&model.UserSubscription{},
		&model.TopUp{},
	))

	t.Cleanup(func() {
		model.DB = originalDB
		model.LOG_DB = originalLogDB
		common.UsingSQLite = originalUsingSQLite
		common.UsingMySQL = originalUsingMySQL
		common.UsingPostgreSQL = originalUsingPostgreSQL
		common.RedisEnabled = originalRedisEnabled

		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db
}

func configureSubscriptionEpayControllerTest(t *testing.T) {
	t.Helper()

	paymentSetting := operation_setting.GetPaymentSetting()
	originalComplianceConfirmed := paymentSetting.ComplianceConfirmed
	originalComplianceTermsVersion := paymentSetting.ComplianceTermsVersion
	originalPrice := operation_setting.Price
	originalPayAddress := operation_setting.PayAddress
	originalEpayID := operation_setting.EpayId
	originalEpayKey := operation_setting.EpayKey
	originalPayMethods := operation_setting.PayMethods
	originalCallbackAddress := operation_setting.CustomCallbackAddress
	originalServerAddress := system_setting.ServerAddress

	t.Cleanup(func() {
		paymentSetting.ComplianceConfirmed = originalComplianceConfirmed
		paymentSetting.ComplianceTermsVersion = originalComplianceTermsVersion
		operation_setting.Price = originalPrice
		operation_setting.PayAddress = originalPayAddress
		operation_setting.EpayId = originalEpayID
		operation_setting.EpayKey = originalEpayKey
		operation_setting.PayMethods = originalPayMethods
		operation_setting.CustomCallbackAddress = originalCallbackAddress
		system_setting.ServerAddress = originalServerAddress
	})

	paymentSetting.ComplianceConfirmed = true
	paymentSetting.ComplianceTermsVersion = operation_setting.CurrentComplianceTermsVersion
	operation_setting.Price = 6.79
	operation_setting.PayAddress = "https://pay.example.com"
	operation_setting.EpayId = "partner-id"
	operation_setting.EpayKey = "secret-key"
	operation_setting.PayMethods = []map[string]string{
		{"type": "alipay", "name": "支付宝"},
	}
	operation_setting.CustomCallbackAddress = ""
	system_setting.ServerAddress = "https://console.example.com"
}

func insertSubscriptionPlanForControllerTest(t *testing.T, plan *model.SubscriptionPlan) {
	t.Helper()
	require.NoError(t, model.DB.Create(plan).Error)
	model.InvalidateSubscriptionPlanCache(plan.Id)
}

func buildJSONContext(t *testing.T, method string, body string, userID int) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(method, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Set("id", userID)
	return c, recorder
}

func TestCalculateSubscriptionEpayMoney_UsesLocalPrice(t *testing.T) {
	configureSubscriptionEpayControllerTest(t)

	payMoney, err := calculateSubscriptionEpayMoney(&model.SubscriptionPlan{
		Id:          1,
		Title:       "Starter",
		PriceAmount: 10,
		Enabled:     true,
	})
	require.NoError(t, err)
	require.True(t, decimal.NewFromFloat(payMoney).Equal(decimal.RequireFromString("67.90")))
}

func TestSubscriptionRequestEpayAmount_ReturnsConvertedAmount(t *testing.T) {
	setupSubscriptionPaymentControllerTestDB(t)
	configureSubscriptionEpayControllerTest(t)

	insertSubscriptionPlanForControllerTest(t, &model.SubscriptionPlan{
		Id:            11,
		Title:         "Starter",
		PriceAmount:   10,
		Currency:      "USD",
		DurationUnit:  model.SubscriptionDurationMonth,
		DurationValue: 1,
		Enabled:       true,
	})

	c, recorder := buildJSONContext(t, http.MethodPost, `{"plan_id":11}`, 101)
	SubscriptionRequestEpayAmount(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	var response subscriptionAmountTestResponse
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &response))
	require.Equal(t, "success", response.Message)
	require.Equal(t, "67.90", response.Data)
}

func TestSubscriptionRequestEpay_StoresConvertedMoneyAndUsesConvertedPurchaseAmount(t *testing.T) {
	setupSubscriptionPaymentControllerTestDB(t)
	configureSubscriptionEpayControllerTest(t)

	insertSubscriptionPlanForControllerTest(t, &model.SubscriptionPlan{
		Id:            21,
		Title:         "Growth",
		PriceAmount:   10,
		Currency:      "USD",
		DurationUnit:  model.SubscriptionDurationMonth,
		DurationValue: 1,
		Enabled:       true,
	})

	originalPurchase := subscriptionEpayPurchase
	var capturedMoney string
	subscriptionEpayPurchase = func(_ string, args subscriptionEpayPurchaseArgs) (string, map[string]string, error) {
		capturedMoney = args.Money
		return "https://pay.example.com/submit.php", map[string]string{"money": args.Money}, nil
	}
	t.Cleanup(func() {
		subscriptionEpayPurchase = originalPurchase
	})

	c, recorder := buildJSONContext(t, http.MethodPost, `{"plan_id":21,"payment_method":"alipay"}`, 202)
	SubscriptionRequestEpay(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	var response subscriptionPayTestResponse
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &response))
	require.Equal(t, "success", response.Message)
	require.Equal(t, "67.90", capturedMoney)
	require.Equal(t, "67.90", response.Data["money"])

	var order model.SubscriptionOrder
	require.NoError(t, model.DB.First(&order).Error)
	require.True(t, decimal.NewFromFloat(order.Money).Equal(decimal.RequireFromString("67.90")))
}

func TestSubscriptionRequestEpayAmount_RejectsDisabledPlan(t *testing.T) {
	setupSubscriptionPaymentControllerTestDB(t)
	configureSubscriptionEpayControllerTest(t)

	plan := &model.SubscriptionPlan{
		Id:            31,
		Title:         "Disabled",
		PriceAmount:   10,
		Currency:      "USD",
		DurationUnit:  model.SubscriptionDurationMonth,
		DurationValue: 1,
		Enabled:       false,
	}
	insertSubscriptionPlanForControllerTest(t, plan)
	require.NoError(t, model.DB.Model(&model.SubscriptionPlan{}).Where("id = ?", plan.Id).Update("enabled", false).Error)
	model.InvalidateSubscriptionPlanCache(plan.Id)

	c, recorder := buildJSONContext(t, http.MethodPost, `{"plan_id":31}`, 303)
	SubscriptionRequestEpayAmount(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	var response subscriptionAmountTestResponse
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &response))
	require.Equal(t, "error", response.Message)
	require.NotEmpty(t, response.Data)
}
