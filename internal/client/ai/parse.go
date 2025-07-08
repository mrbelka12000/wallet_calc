package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) do(ctx context.Context, in, out any) (err error) {

	var (
		reqUrl   = apiURL + pathCompletions
		httpReq  *http.Request
		httpResp *http.Response
		reqBody  []byte
		respBody []byte
	)

	// Log request details.
	defer func() {
		log := c.log.With(
			"request_method", http.MethodPost,
			"request_url", reqUrl,
			"request_body", string(reqBody),
			"response_body", string(respBody),
		)

		if httpReq != nil {
			log.With("request_headers", httpReq.Header)
		}

		if httpResp != nil {
			log.With("response_headers", httpResp.Header)
			log.With("response_status", httpResp.Status)
		}

		//log.Info("execute request to ai")
	}()

	if in != nil {
		reqBody, _ = json.Marshal(in)
	}

	httpReq, err = http.NewRequest(http.MethodPost, apiURL+pathCompletions, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	httpResp, err = c.hc.Do(httpReq.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to do http request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(respBody, &out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}
