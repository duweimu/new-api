package relay

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
)

type fakeResponsesCompatAdaptor struct {
	capturedResponsesReq *dto.OpenAIResponsesRequest
}

func (a *fakeResponsesCompatAdaptor) Init(info *relaycommon.RelayInfo) {}

func (a *fakeResponsesCompatAdaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	return "https://example.com/backend-api/codex/responses", nil
}

func (a *fakeResponsesCompatAdaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error {
	return nil
}

func (a *fakeResponsesCompatAdaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	return nil, errors.New("unexpected ConvertOpenAIRequest call")
}

func (a *fakeResponsesCompatAdaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error) {
	return nil, errors.New("unexpected ConvertRerankRequest call")
}

func (a *fakeResponsesCompatAdaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error) {
	return nil, errors.New("unexpected ConvertEmbeddingRequest call")
}

func (a *fakeResponsesCompatAdaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	return nil, errors.New("unexpected ConvertAudioRequest call")
}

func (a *fakeResponsesCompatAdaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error) {
	return nil, errors.New("unexpected ConvertImageRequest call")
}

func (a *fakeResponsesCompatAdaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error) {
	reqCopy := request
	a.capturedResponsesReq = &reqCopy
	return reqCopy, nil
}

func (a *fakeResponsesCompatAdaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error) {
	if a.capturedResponsesReq == nil || !a.capturedResponsesReq.IsStream(c) || !info.IsStream {
		return nil, errors.New("Stream must be set to true")
	}

	body := strings.Join([]string{
		`data: {"type":"response.created","response":{"id":"resp_1","model":"gpt-5.5","created_at":1}}`,
		"",
		`data: {"type":"response.output_text.delta","delta":"hi"}`,
		"",
		`data: {"type":"response.completed","response":{"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2}}}`,
		"",
		"data: [DONE]",
		"",
	}, "\n")

	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}, nil
}

func (a *fakeResponsesCompatAdaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError) {
	return nil, types.NewError(errors.New("unexpected DoResponse call"), types.ErrorCodeBadResponse)
}

func (a *fakeResponsesCompatAdaptor) GetModelList() []string {
	return []string{"gpt-5.5"}
}

func (a *fakeResponsesCompatAdaptor) GetChannelName() string {
	return "fake"
}

func (a *fakeResponsesCompatAdaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.ClaudeRequest) (any, error) {
	return nil, errors.New("unexpected ConvertClaudeRequest call")
}

func (a *fakeResponsesCompatAdaptor) ConvertGeminiRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeminiChatRequest) (any, error) {
	return nil, errors.New("unexpected ConvertGeminiRequest call")
}

func TestChatCompletionsViaResponsesForcesCodexStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)
	constant.StreamingTimeout = 30

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader(`{"model":"gpt-5.5","messages":[{"role":"user","content":"hi"}],"stream":false}`))
	c.Request.Header.Set("Content-Type", "application/json")
	common.SetContextKey(c, constant.ContextKeyIsStream, false)

	info := &relaycommon.RelayInfo{
		RelayMode:       relayconstant.RelayModeChatCompletions,
		RelayFormat:     types.RelayFormatOpenAI,
		OriginModelName: "gpt-5.5",
		IsStream:        false,
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiType:     constant.APITypeCodex,
			ChannelType: constant.ChannelTypeCodex,
			UpstreamModelName: "gpt-5.5",
		},
	}

	request := &dto.GeneralOpenAIRequest{
		Model:  "gpt-5.5",
		Stream: common.GetPointer(false),
		Messages: []dto.Message{
			{
				Role:    "user",
				Content: "hi",
			},
		},
	}

	adaptor := &fakeResponsesCompatAdaptor{}
	usage, newAPIError := chatCompletionsViaResponses(c, info, adaptor, request)
	if newAPIError != nil {
		t.Fatalf("chatCompletionsViaResponses returned error: %v", newAPIError)
	}
	if usage == nil {
		t.Fatal("expected usage, got nil")
	}
	if adaptor.capturedResponsesReq == nil {
		t.Fatal("expected captured responses request")
	}
	if !adaptor.capturedResponsesReq.IsStream(c) {
		t.Fatal("expected converted responses request to force stream=true")
	}
	if !info.IsStream {
		t.Fatal("expected relay info to be marked as stream")
	}
	if !common.GetContextKeyBool(c, constant.ContextKeyIsStream) {
		t.Fatal("expected request context stream flag to be true")
	}
}

func TestForceCodexResponsesStreamingUpdatesRequestAndContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", strings.NewReader(`{}`))
	common.SetContextKey(c, constant.ContextKeyIsStream, false)

	info := &relaycommon.RelayInfo{
		RelayMode: relayconstant.RelayModeResponses,
		IsStream:  false,
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiType:     constant.APITypeCodex,
			ChannelType: constant.ChannelTypeCodex,
		},
	}
	request := &dto.OpenAIResponsesRequest{Model: "gpt-5.5", Stream: common.GetPointer(false)}

	if !forceCodexResponsesStreaming(c, info, request) {
		t.Fatal("expected forceCodexResponsesStreaming to report a mutation")
	}
	if !request.IsStream(c) {
		t.Fatal("expected request stream=true after forcing")
	}
	if !info.IsStream {
		t.Fatal("expected relay info to be marked as stream")
	}
	if !common.GetContextKeyBool(c, constant.ContextKeyIsStream) {
		t.Fatal("expected context stream flag to be true")
	}
}

func TestNormalizeCodexResponsesPassthroughBodyForcesStream(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(
		http.MethodPost,
		"/v1/responses",
		strings.NewReader(`{"model":"gpt-5.5","input":[{"role":"user","content":[{"type":"input_text","text":"hi"}]}],"stream":false}`),
	)
	c.Request.Header.Set("Content-Type", "application/json")
	common.SetContextKey(c, constant.ContextKeyIsStream, false)

	info := &relaycommon.RelayInfo{
		RelayMode: relayconstant.RelayModeResponses,
		IsStream:  false,
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiType:     constant.APITypeCodex,
			ChannelType: constant.ChannelTypeCodex,
		},
	}

	if err := normalizeCodexResponsesPassthroughBody(c, info); err != nil {
		t.Fatalf("normalizeCodexResponsesPassthroughBody returned error: %v", err)
	}

	storage, err := common.GetBodyStorage(c)
	if err != nil {
		t.Fatalf("GetBodyStorage returned error: %v", err)
	}
	bodyBytes, err := storage.Bytes()
	if err != nil {
		t.Fatalf("body storage bytes returned error: %v", err)
	}
	var request dto.OpenAIResponsesRequest
	if err := common.Unmarshal(bodyBytes, &request); err != nil {
		t.Fatalf("failed to unmarshal normalized body: %v", err)
	}
	if !request.IsStream(c) {
		t.Fatal("expected passthrough request body to force stream=true")
	}
	if !info.IsStream {
		t.Fatal("expected relay info to be marked as stream")
	}
	if !common.GetContextKeyBool(c, constant.ContextKeyIsStream) {
		t.Fatal("expected context stream flag to be true")
	}
	if info.UpstreamRequestBodySize != storage.Size() {
		t.Fatalf("expected upstream request size %d, got %d", storage.Size(), info.UpstreamRequestBodySize)
	}
}
