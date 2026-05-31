package controller

import (
	"net/http"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/stretchr/testify/require"
)

type subscriptionBalancePayTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func TestSubscriptionRequestBalancePay_DeductsQuotaAndCreatesSubscription(t *testing.T) {
	setupSubscriptionPaymentControllerTestDB(t)
	configureSubscriptionEpayControllerTest(t)

	user := &model.User{
		Id:       501,
		Username: "subscription-balance-user",
		Password: "password",
		Role:     common.RoleCommonUser,
		Status:   common.UserStatusEnabled,
		Quota:    25 * int(common.QuotaPerUnit),
		Group:    "default",
	}
	require.NoError(t, model.DB.Create(user).Error)

	insertSubscriptionPlanForControllerTest(t, &model.SubscriptionPlan{
		Id:            41,
		Title:         "Balance Plan",
		PriceAmount:   12,
		Currency:      "USD",
		DurationUnit:  model.SubscriptionDurationMonth,
		DurationValue: 1,
		Enabled:       true,
		TotalAmount:   20 * int64(common.QuotaPerUnit),
	})

	c, recorder := buildJSONContext(t, http.MethodPost, `{"plan_id":41}`, user.Id)
	SubscriptionRequestBalancePay(c)

	require.Equal(t, http.StatusOK, recorder.Code)

	var response subscriptionBalancePayTestResponse
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &response))
	require.True(t, response.Success)

	var refreshed model.User
	require.NoError(t, model.DB.First(&refreshed, "id = ?", user.Id).Error)
	require.Equal(t, (25-12)*int(common.QuotaPerUnit), refreshed.Quota)

	var subscriptions []model.UserSubscription
	require.NoError(t, model.DB.Where("user_id = ?", user.Id).Find(&subscriptions).Error)
	require.Len(t, subscriptions, 1)

	var orders []model.SubscriptionOrder
	require.NoError(t, model.DB.Where("user_id = ?", user.Id).Find(&orders).Error)
	require.Len(t, orders, 1)
	require.Equal(t, model.PaymentMethodBalance, orders[0].PaymentMethod)
	require.Equal(t, model.PaymentProviderBalance, orders[0].PaymentProvider)
	require.Equal(t, common.TopUpStatusSuccess, orders[0].Status)
}

func TestSubscriptionRequestBalancePay_RejectsDisallowedBalanceRedemption(t *testing.T) {
	setupSubscriptionPaymentControllerTestDB(t)
	configureSubscriptionEpayControllerTest(t)

	user := &model.User{
		Id:       502,
		Username: "subscription-balance-disabled",
		Password: "password",
		Role:     common.RoleCommonUser,
		Status:   common.UserStatusEnabled,
		Quota:    25 * int(common.QuotaPerUnit),
		Group:    "default",
	}
	require.NoError(t, model.DB.Create(user).Error)

	disabled := false
	insertSubscriptionPlanForControllerTest(t, &model.SubscriptionPlan{
		Id:              42,
		Title:           "No Balance Plan",
		PriceAmount:     12,
		Currency:        "USD",
		DurationUnit:    model.SubscriptionDurationMonth,
		DurationValue:   1,
		Enabled:         true,
		AllowBalancePay: &disabled,
		TotalAmount:     20 * int64(common.QuotaPerUnit),
	})

	c, recorder := buildJSONContext(t, http.MethodPost, `{"plan_id":42}`, user.Id)
	SubscriptionRequestBalancePay(c)

	require.Equal(t, http.StatusOK, recorder.Code)

	var response subscriptionBalancePayTestResponse
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &response))
	require.False(t, response.Success)
	require.Contains(t, response.Message, "余额")

	var refreshed model.User
	require.NoError(t, model.DB.First(&refreshed, "id = ?", user.Id).Error)
	require.Equal(t, 25*int(common.QuotaPerUnit), refreshed.Quota)
}
