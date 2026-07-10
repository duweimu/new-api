package relay

import (
	"testing"

	"github.com/QuantumNous/new-api/constant"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/setting/model_setting"
)

func TestShouldUseResponsesCompatibility_UsesConfiguredPolicy(t *testing.T) {
	original := model_setting.GetGlobalSettings().ChatCompletionsToResponsesPolicy
	model_setting.GetGlobalSettings().ChatCompletionsToResponsesPolicy = model_setting.ChatCompletionsToResponsesPolicy{
		Enabled:       true,
		ChannelIDs:    []int{9},
		ModelPatterns: []string{`^gpt-5\.4$`},
	}
	defer func() {
		model_setting.GetGlobalSettings().ChatCompletionsToResponsesPolicy = original
	}()

	info := &relaycommon.RelayInfo{
		OriginModelName: "gpt-5.4",
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelId:   9,
			ApiType:     constant.APITypeOpenAI,
			ChannelType: constant.ChannelTypeOpenAI,
		},
	}

	if !shouldUseResponsesCompatibility(info) {
		t.Fatal("expected configured compatibility policy to enable responses mode")
	}
}

func TestShouldUseResponsesCompatibility_CodexForcesCompatibilityWithoutPolicy(t *testing.T) {
	original := model_setting.GetGlobalSettings().ChatCompletionsToResponsesPolicy
	model_setting.GetGlobalSettings().ChatCompletionsToResponsesPolicy = model_setting.ChatCompletionsToResponsesPolicy{}
	defer func() {
		model_setting.GetGlobalSettings().ChatCompletionsToResponsesPolicy = original
	}()

	info := &relaycommon.RelayInfo{
		OriginModelName: "gpt-5.4",
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelId:   3,
			ApiType:     constant.APITypeCodex,
			ChannelType: constant.ChannelTypeCodex,
		},
	}

	if !shouldUseResponsesCompatibility(info) {
		t.Fatal("expected codex channels to force responses compatibility by default")
	}
}

func TestShouldUseResponsesCompatibility_NonCodexWithoutPolicyFalse(t *testing.T) {
	original := model_setting.GetGlobalSettings().ChatCompletionsToResponsesPolicy
	model_setting.GetGlobalSettings().ChatCompletionsToResponsesPolicy = model_setting.ChatCompletionsToResponsesPolicy{}
	defer func() {
		model_setting.GetGlobalSettings().ChatCompletionsToResponsesPolicy = original
	}()

	info := &relaycommon.RelayInfo{
		OriginModelName: "gpt-5.4",
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelId:   3,
			ApiType:     constant.APITypeOpenAI,
			ChannelType: constant.ChannelTypeOpenAI,
		},
	}

	if shouldUseResponsesCompatibility(info) {
		t.Fatal("expected non-codex channels without policy to keep default request handling")
	}
}
