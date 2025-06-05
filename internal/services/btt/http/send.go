package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *client) Send(
	ctx context.Context,
	method string,
	jsonPayload map[string]any,
	extraPayload map[string]string,
) ([]byte, error) {
	var buf bytes.Buffer
	query := url.Values{}
	if jsonPayload != nil {
		if err := json.NewEncoder(&buf).Encode(jsonPayload); err != nil {
			return nil, fmt.Errorf("cannot encode payload: %w", err)
		}
		query.Set("json", buf.String())
	}
	for key, val := range extraPayload {
		query.Set(key, val)
	}

	return c.send(ctx, method, query)
}

func (c *client) send(ctx context.Context, method string, query url.Values) ([]byte, error) {
	u := c.bttURL
	u.Path = "/" + method + "/"
	u.RawQuery = strings.ReplaceAll(query.Encode(), "+", "%20")
	u.ForceQuery = true

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create new request: %w", err)
	}

	defer func() {
		if err != nil {
			c.logger.ErrorContext(ctx, "cannot send request", "error", err, "method", method)
		}
	}()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body: %w", err)
	}

	c.logger.InfoContext(ctx, "sent request", "method", method)

	return body, nil
}
