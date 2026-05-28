package relay

import (
	"io"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/gin-gonic/gin"
)

func shouldForceCodexResponsesStreaming(info *relaycommon.RelayInfo) bool {
	if info == nil || info.ChannelMeta == nil {
		return false
	}
	if info.RelayMode != relayconstant.RelayModeResponses {
		return false
	}
	if info.ApiType != constant.APITypeCodex && info.ChannelType != constant.ChannelTypeCodex {
		return false
	}
	return !info.IsStream
}

func setRelayStream(c *gin.Context, info *relaycommon.RelayInfo) {
	if info == nil {
		return
	}
	info.IsStream = true
	if c != nil {
		common.SetContextKey(c, constant.ContextKeyIsStream, true)
	}
}

func forceCodexResponsesStreaming(c *gin.Context, info *relaycommon.RelayInfo, request *dto.OpenAIResponsesRequest) bool {
	if !shouldForceCodexResponsesStreaming(info) {
		return false
	}
	setRelayStream(c, info)
	if request != nil {
		request.Stream = common.GetPointer(true)
	}
	return true
}

func normalizeCodexResponsesPassthroughBody(c *gin.Context, info *relaycommon.RelayInfo) error {
	if !shouldForceCodexResponsesStreaming(info) {
		return nil
	}
	storage, err := common.GetBodyStorage(c)
	if err != nil {
		return err
	}
	requestBody, err := storage.Bytes()
	if err != nil {
		return err
	}

	request := &dto.OpenAIResponsesRequest{}
	if err := common.Unmarshal(requestBody, request); err != nil {
		return err
	}
	forceCodexResponsesStreaming(c, info, request)

	jsonData, err := common.Marshal(request)
	if err != nil {
		return err
	}

	newStorage, err := common.CreateBodyStorage(jsonData)
	if err != nil {
		return err
	}
	if _, err := newStorage.Seek(0, io.SeekStart); err != nil {
		_ = newStorage.Close()
		return err
	}

	_ = storage.Close()
	c.Set(common.KeyBodyStorage, newStorage)
	c.Request.Body = io.NopCloser(newStorage)
	info.UpstreamRequestBodySize = newStorage.Size()
	return nil
}
