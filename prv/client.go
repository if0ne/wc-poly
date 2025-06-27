package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/auth"
	"github.com/ollama/ollama/envconfig"
	"github.com/ollama/ollama/format"
	"github.com/ollama/ollama/openai"
	"github.com/ollama/ollama/version"
)

type Client struct {
	base *url.URL
	http *http.Client
}

func checkError(resp *http.Response, body []byte) error {
	if resp.StatusCode < http.StatusBadRequest {
		return nil
	}

	apiError := api.StatusError{StatusCode: resp.StatusCode}

	err := json.Unmarshal(body, &apiError)
	if err != nil {
		// Use the full body as the message if we fail to decode a response.
		apiError.ErrorMessage = string(body)
	}

	return apiError
}

// ClientFromEnvironment creates a new [Client] using configuration from the
// environment variable OLLAMA_HOST, which points to the network host and
// port on which the ollama service is listening. The format of this variable
// is:
//
//	<scheme>://<host>:<port>
//
// If the variable is not specified, a default ollama host and port will be
// used.
func ClientFromEnvironment() (*Client, error) {
	return &Client{
		base: envconfig.Host(),
		http: http.DefaultClient,
	}, nil
}

func NewClient(base *url.URL, http *http.Client) *Client {
	return &Client{
		base: base,
		http: http,
	}
}

func getAuthorizationToken(ctx context.Context, challenge string) (string, error) {
	token, err := auth.Sign(ctx, []byte(challenge))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *Client) do(ctx context.Context, method, path string, reqData, respData any) error {
	var reqBody io.Reader
	var data []byte
	var err error

	switch reqData := reqData.(type) {
	case io.Reader:
		// reqData is already an io.Reader
		reqBody = reqData
	case nil:
		// noop
	default:
		data, err = json.Marshal(reqData)
		if err != nil {
			return err
		}

		reqBody = bytes.NewReader(data)
	}

	requestURL := c.base.JoinPath(path)

	var token string
	if envconfig.UseAuth() || c.base.Hostname() == "ollama.com" {
		now := strconv.FormatInt(time.Now().Unix(), 10)
		chal := fmt.Sprintf("%s,%s?ts=%s", method, path, now)
		token, err = getAuthorizationToken(ctx, chal)
		if err != nil {
			return err
		}

		q := requestURL.Query()
		q.Set("ts", now)
		requestURL.RawQuery = q.Encode()
	}

	request, err := http.NewRequestWithContext(ctx, method, requestURL.String(), reqBody)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", fmt.Sprintf("ollama/%s (%s %s) Go/%s", version.Version, runtime.GOARCH, runtime.GOOS, runtime.Version()))

	if token != "" {
		request.Header.Set("Authorization", token)
	}

	respObj, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer respObj.Body.Close()

	respBody, err := io.ReadAll(respObj.Body)
	if err != nil {
		return err
	}

	if err := checkError(respObj, respBody); err != nil {
		return err
	}

	if len(respBody) > 0 && respData != nil {
		if err := json.Unmarshal(respBody, respData); err != nil {
			return err
		}
	}
	return nil
}

const maxBufferSize = 512 * format.KiloByte

func (c *Client) stream(ctx context.Context, method, path string, data any, fn func([]byte) error) error {
	var buf io.Reader
	if data != nil {
		bts, err := json.Marshal(data)
		if err != nil {
			return err
		}

		buf = bytes.NewBuffer(bts)
	}

	requestURL := c.base.JoinPath(path)

	var token string
	if envconfig.UseAuth() || c.base.Hostname() == "ollama.com" {
		var err error
		now := strconv.FormatInt(time.Now().Unix(), 10)
		chal := fmt.Sprintf("%s,%s?ts=%s", method, path, now)
		token, err = getAuthorizationToken(ctx, chal)
		if err != nil {
			return err
		}

		q := requestURL.Query()
		q.Set("ts", now)
		requestURL.RawQuery = q.Encode()
	}

	request, err := http.NewRequestWithContext(ctx, method, requestURL.String(), buf)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/x-ndjson")
	request.Header.Set("User-Agent", fmt.Sprintf("ollama/%s (%s %s) Go/%s", version.Version, runtime.GOARCH, runtime.GOOS, runtime.Version()))

	if token != "" {
		request.Header.Set("Authorization", token)
	}

	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	// increase the buffer size to avoid running out of space
	scanBuf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(scanBuf, maxBufferSize)
	for scanner.Scan() {
		var errorResponse struct {
			Error string `json:"error,omitempty"`
		}

		bts := scanner.Bytes()
		if err := json.Unmarshal(bts, &errorResponse); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		if errorResponse.Error != "" {
			return errors.New(errorResponse.Error)
		}

		if response.StatusCode >= http.StatusBadRequest {
			return api.StatusError{
				StatusCode:   response.StatusCode,
				Status:       response.Status,
				ErrorMessage: errorResponse.Error,
			}
		}

		if err := fn(bts); err != nil {
			return err
		}
	}

	return nil
}

// GenerateResponseFunc is a function that [Client.Generate] invokes every time
// a response is received from the service. If this function returns an error,
// [Client.Generate] will stop generating and return this error.
type GenerateResponseFunc func(api.GenerateResponse) error

// Generate generates a response for a given prompt. The req parameter should
// be populated with prompt details. fn is called for each response (there may
// be multiple responses, e.g. in case streaming is enabled).
func (c *Client) Generate(ctx context.Context, req *api.GenerateRequest, fn GenerateResponseFunc) error {
	return c.stream(ctx, http.MethodPost, "/api/generate", req, func(bts []byte) error {
		var resp api.GenerateResponse
		if err := json.Unmarshal(bts, &resp); err != nil {
			return err
		}

		return fn(resp)
	})
}

// ChatResponseFunc is a function that [Client.Chat] invokes every time
// a response is received from the service. If this function returns an error,
// [Client.Chat] will stop generating and return this error.
type ChatResponseFunc func(openai.ChatCompletion) error

// Chat generates the next message in a chat. [ChatRequest] may contain a
// sequence of messages which can be used to maintain chat history with a model.
// fn is called for each response (there may be multiple responses, e.g. if case
// streaming is enabled).
func (c *Client) Chat(ctx context.Context, req *openai.ChatCompletionRequest, fn ChatResponseFunc) error {
	return c.stream(ctx, http.MethodPost, "/v1/chat/completions", req, func(bts []byte) error {
		var resp openai.ChatCompletion
		if err := json.Unmarshal(bts, &resp); err != nil {
			return err
		}

		return fn(resp)
	})
}

// PullProgressFunc is a function that [Client.Pull] invokes every time there
// is progress with a "pull" request sent to the service. If this function
// returns an error, [Client.Pull] will stop the process and return this error.
type PullProgressFunc func(api.ProgressResponse) error

// Pull downloads a model from the ollama library. fn is called each time
// progress is made on the request and can be used to display a progress bar,
// etc.
func (c *Client) Pull(ctx context.Context, req *api.PullRequest, fn PullProgressFunc) error {
	return c.stream(ctx, http.MethodPost, "/api/pull", req, func(bts []byte) error {
		var resp api.ProgressResponse
		if err := json.Unmarshal(bts, &resp); err != nil {
			return err
		}

		return fn(resp)
	})
}

// List lists models that are available locally.
func (c *Client) List(ctx context.Context) (*api.ListResponse, error) {
	var lr api.ListResponse
	if err := c.do(ctx, http.MethodGet, "/api/tags", nil, &lr); err != nil {
		return nil, err
	}
	return &lr, nil
}

// Delete deletes a model and its data.
func (c *Client) Delete(ctx context.Context, req *api.DeleteRequest) error {
	if err := c.do(ctx, http.MethodDelete, "/api/delete", req, nil); err != nil {
		return err
	}
	return nil
}
