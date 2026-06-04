package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type userControllerTestResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func setupUserControllerTestDB(t *testing.T) *gorm.DB {
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

	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Log{}))

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

func buildUserControllerContext(
	t *testing.T,
	method string,
	body string,
	actorID int,
	actorRole int,
	actorUsername string,
	targetUserID int,
) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(method, "/", strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req
	c.Set("id", actorID)
	c.Set("role", actorRole)
	c.Set("username", actorUsername)
	if targetUserID != 0 {
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(targetUserID)}}
	}
	return c, recorder
}

func seedManagedUser(t *testing.T, user *model.User) {
	t.Helper()
	require.NoError(t, model.DB.Create(user).Error)
}

func decodeUserControllerResponse(t *testing.T, recorder *httptest.ResponseRecorder) userControllerTestResponse {
	t.Helper()

	require.Equal(t, http.StatusOK, recorder.Code)

	var response userControllerTestResponse
	require.NoError(t, common.Unmarshal(recorder.Body.Bytes(), &response))
	return response
}

func TestGetUser_AllowsAdminToReadSelfButBlocksPeerAdmin(t *testing.T) {
	setupUserControllerTestDB(t)

	selfAdmin := &model.User{
		Id:          11,
		Username:    "admin-self-read",
		DisplayName: "Admin Self",
		Password:    "password",
		Role:        common.RoleAdminUser,
		Status:      common.UserStatusEnabled,
		Group:       "default",
		Quota:       2 * int(common.QuotaPerUnit),
		AffCode:     "asr1",
	}
	peerAdmin := &model.User{
		Id:          12,
		Username:    "admin-peer-read",
		DisplayName: "Admin Peer",
		Password:    "password",
		Role:        common.RoleAdminUser,
		Status:      common.UserStatusEnabled,
		Group:       "default",
		AffCode:     "apr1",
	}
	seedManagedUser(t, selfAdmin)
	seedManagedUser(t, peerAdmin)

	t.Run("self admin can read self", func(t *testing.T) {
		c, recorder := buildUserControllerContext(
			t, http.MethodGet, "", selfAdmin.Id, common.RoleAdminUser, selfAdmin.Username, selfAdmin.Id,
		)
		GetUser(c)

		response := decodeUserControllerResponse(t, recorder)
		require.True(t, response.Success)

		var returnedUser model.User
		require.NoError(t, common.Unmarshal(response.Data, &returnedUser))
		require.Equal(t, selfAdmin.Id, returnedUser.Id)
		require.Equal(t, selfAdmin.Username, returnedUser.Username)
	})

	t.Run("admin still cannot read peer admin", func(t *testing.T) {
		c, recorder := buildUserControllerContext(
			t, http.MethodGet, "", selfAdmin.Id, common.RoleAdminUser, selfAdmin.Username, peerAdmin.Id,
		)
		GetUser(c)

		response := decodeUserControllerResponse(t, recorder)
		require.False(t, response.Success)
	})
}

func TestUpdateUser_AllowsAdminToUpdateSelfButBlocksPeerAdmin(t *testing.T) {
	setupUserControllerTestDB(t)

	selfAdmin := &model.User{
		Id:          21,
		Username:    "admin-self-update",
		DisplayName: "Old Name",
		Password:    "password",
		Role:        common.RoleAdminUser,
		Status:      common.UserStatusEnabled,
		Group:       "default",
		Remark:      "old remark",
		AffCode:     "asu1",
	}
	peerAdmin := &model.User{
		Id:          22,
		Username:    "admin-peer-update",
		DisplayName: "Peer Name",
		Password:    "password",
		Role:        common.RoleAdminUser,
		Status:      common.UserStatusEnabled,
		Group:       "default",
		Remark:      "peer remark",
		AffCode:     "apu1",
	}
	seedManagedUser(t, selfAdmin)
	seedManagedUser(t, peerAdmin)

	t.Run("self admin can update own editable fields", func(t *testing.T) {
		body := `{"id":21,"username":"admin-self-update","display_name":"New Name","group":"vip","remark":"new remark"}`
		c, recorder := buildUserControllerContext(
			t, http.MethodPut, body, selfAdmin.Id, common.RoleAdminUser, selfAdmin.Username, 0,
		)
		UpdateUser(c)

		response := decodeUserControllerResponse(t, recorder)
		require.True(t, response.Success)

		refreshed, err := model.GetUserById(selfAdmin.Id, true)
		require.NoError(t, err)
		require.Equal(t, "New Name", refreshed.DisplayName)
		require.Equal(t, "vip", refreshed.Group)
		require.Equal(t, "new remark", refreshed.Remark)
		require.Equal(t, common.RoleAdminUser, refreshed.Role)
	})

	t.Run("admin still cannot update peer admin", func(t *testing.T) {
		body := `{"id":22,"username":"admin-peer-update","display_name":"Blocked Change","group":"blocked","remark":"blocked"}`
		c, recorder := buildUserControllerContext(
			t, http.MethodPut, body, selfAdmin.Id, common.RoleAdminUser, selfAdmin.Username, 0,
		)
		UpdateUser(c)

		response := decodeUserControllerResponse(t, recorder)
		require.False(t, response.Success)

		refreshed, err := model.GetUserById(peerAdmin.Id, true)
		require.NoError(t, err)
		require.Equal(t, "Peer Name", refreshed.DisplayName)
		require.Equal(t, "default", refreshed.Group)
		require.Equal(t, "peer remark", refreshed.Remark)
	})
}

func TestManageUser_AddQuota_AllowsAdminSelfButBlocksPeerAdmin(t *testing.T) {
	setupUserControllerTestDB(t)

	selfAdmin := &model.User{
		Id:          31,
		Username:    "admin-self-quota",
		DisplayName: "Quota Self",
		Password:    "password",
		Role:        common.RoleAdminUser,
		Status:      common.UserStatusEnabled,
		Group:       "default",
		Quota:       10 * int(common.QuotaPerUnit),
		AffCode:     "asq1",
	}
	peerAdmin := &model.User{
		Id:          32,
		Username:    "admin-peer-quota",
		DisplayName: "Quota Peer",
		Password:    "password",
		Role:        common.RoleAdminUser,
		Status:      common.UserStatusEnabled,
		Group:       "default",
		Quota:       6 * int(common.QuotaPerUnit),
		AffCode:     "apq1",
	}
	seedManagedUser(t, selfAdmin)
	seedManagedUser(t, peerAdmin)

	t.Run("self admin add quota", func(t *testing.T) {
		body := fmt.Sprintf(`{"id":31,"action":"add_quota","mode":"add","value":%d}`, int(common.QuotaPerUnit))
		c, recorder := buildUserControllerContext(
			t, http.MethodPost, body, selfAdmin.Id, common.RoleAdminUser, selfAdmin.Username, 0,
		)
		ManageUser(c)

		response := decodeUserControllerResponse(t, recorder)
		require.True(t, response.Success)

		refreshed, err := model.GetUserById(selfAdmin.Id, true)
		require.NoError(t, err)
		require.Equal(t, 11*int(common.QuotaPerUnit), refreshed.Quota)
	})

	t.Run("self admin subtract quota", func(t *testing.T) {
		body := fmt.Sprintf(`{"id":31,"action":"add_quota","mode":"subtract","value":%d}`, 2*int(common.QuotaPerUnit))
		c, recorder := buildUserControllerContext(
			t, http.MethodPost, body, selfAdmin.Id, common.RoleAdminUser, selfAdmin.Username, 0,
		)
		ManageUser(c)

		response := decodeUserControllerResponse(t, recorder)
		require.True(t, response.Success)

		refreshed, err := model.GetUserById(selfAdmin.Id, true)
		require.NoError(t, err)
		require.Equal(t, 9*int(common.QuotaPerUnit), refreshed.Quota)
	})

	t.Run("self admin override quota", func(t *testing.T) {
		body := fmt.Sprintf(`{"id":31,"action":"add_quota","mode":"override","value":%d}`, 3*int(common.QuotaPerUnit))
		c, recorder := buildUserControllerContext(
			t, http.MethodPost, body, selfAdmin.Id, common.RoleAdminUser, selfAdmin.Username, 0,
		)
		ManageUser(c)

		response := decodeUserControllerResponse(t, recorder)
		require.True(t, response.Success)

		refreshed, err := model.GetUserById(selfAdmin.Id, true)
		require.NoError(t, err)
		require.Equal(t, 3*int(common.QuotaPerUnit), refreshed.Quota)
	})

	t.Run("admin still cannot adjust peer admin quota", func(t *testing.T) {
		body := fmt.Sprintf(`{"id":32,"action":"add_quota","mode":"add","value":%d}`, int(common.QuotaPerUnit))
		c, recorder := buildUserControllerContext(
			t, http.MethodPost, body, selfAdmin.Id, common.RoleAdminUser, selfAdmin.Username, 0,
		)
		ManageUser(c)

		response := decodeUserControllerResponse(t, recorder)
		require.False(t, response.Success)

		refreshed, err := model.GetUserById(peerAdmin.Id, true)
		require.NoError(t, err)
		require.Equal(t, 6*int(common.QuotaPerUnit), refreshed.Quota)
	})
}

func TestManageUser_NonQuotaSelfAdminActionsRemainBlocked(t *testing.T) {
	setupUserControllerTestDB(t)

	selfAdmin := &model.User{
		Id:          41,
		Username:    "admin-self-guard",
		DisplayName: "Guard Self",
		Password:    "password",
		Role:        common.RoleAdminUser,
		Status:      common.UserStatusEnabled,
		Group:       "default",
		Quota:       5 * int(common.QuotaPerUnit),
		AffCode:     "asg1",
	}
	seedManagedUser(t, selfAdmin)

	actions := []string{"disable", "demote", "delete"}
	for _, action := range actions {
		t.Run(action, func(t *testing.T) {
			body := fmt.Sprintf(`{"id":41,"action":"%s"}`, action)
			c, recorder := buildUserControllerContext(
				t, http.MethodPost, body, selfAdmin.Id, common.RoleAdminUser, selfAdmin.Username, 0,
			)
			ManageUser(c)

			response := decodeUserControllerResponse(t, recorder)
			require.False(t, response.Success)

			refreshed, err := model.GetUserById(selfAdmin.Id, true)
			require.NoError(t, err)
			require.Equal(t, common.RoleAdminUser, refreshed.Role)
			require.Equal(t, common.UserStatusEnabled, refreshed.Status)
		})
	}
}
