package service

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
)

func TestStreamResponseOpenAI2Claude_EmitsStopForFinishOnlyChunkWithCachedUsage(t *testing.T) {
	info := &relaycommon.RelayInfo{
		ClaudeConvertInfo: &relaycommon.ClaudeConvertInfo{
			LastMessagesType: relaycommon.LastMessageTypeText,
			Index:            0,
			Usage: &dto.Usage{
				PromptTokens:     3,
				CompletionTokens: 5,
				TotalTokens:      8,
			},
		},
	}

	stopChunk := &dto.ChatCompletionsStreamResponse{
		Choices: []dto.ChatCompletionsStreamResponseChoice{
			{
				FinishReason: common.GetPointer("stop"),
			},
		},
	}

	responses := StreamResponseOpenAI2Claude(stopChunk, info)
	if len(responses) != 3 {
		t.Fatalf("expected 3 Claude events, got %d", len(responses))
	}
	if responses[0].Type != "content_block_stop" {
		t.Fatalf("expected first event to be content_block_stop, got %q", responses[0].Type)
	}
	if responses[1].Type != "message_delta" {
		t.Fatalf("expected second event to be message_delta, got %q", responses[1].Type)
	}
	if responses[1].Usage == nil || responses[1].Usage.OutputTokens != 5 {
		t.Fatalf("expected message_delta usage output_tokens=5, got %+v", responses[1].Usage)
	}
	if responses[2].Type != "message_stop" {
		t.Fatalf("expected third event to be message_stop, got %q", responses[2].Type)
	}
}
