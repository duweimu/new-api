package relay

import (
	"github.com/QuantumNous/new-api/constant"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/service"
)

func shouldUseResponsesCompatibility(info *relaycommon.RelayInfo) bool {
	if info == nil {
		return false
	}

	channelID := 0
	channelType := 0
	apiType := 0
	if info.ChannelMeta != nil {
		channelID = info.ChannelId
		channelType = info.ChannelType
		apiType = info.ApiType
	}

	if apiType == constant.APITypeCodex || channelType == constant.ChannelTypeCodex {
		return true
	}

	return service.ShouldChatCompletionsUseResponsesGlobal(channelID, channelType, info.OriginModelName)
}
