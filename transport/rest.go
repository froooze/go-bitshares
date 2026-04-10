package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// RESTClient performs REST-style requests against a base URL.
type RESTClient struct {
	baseURL *url.URL
	client  *http.Client
}

// NewRESTClient creates a REST client rooted at the provided base URL.
func NewRESTClient(endpoint string) (*RESTClient, error) {
	return NewRESTClientWithHTTPClient(endpoint, nil)
}

// NewRESTClientWithHTTPClient creates a REST client with a custom HTTP client.
func NewRESTClientWithHTTPClient(endpoint string, httpClient *http.Client) (*RESTClient, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("empty endpoint")
	}

	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, fmt.Errorf("invalid endpoint %q", endpoint)
	}

	return &RESTClient{
		baseURL: baseURL,
		client:  normalizeHTTPClient(httpClient),
	}, nil
}

func (c *RESTClient) Close() error { return nil }

// Do issues a REST-style request using an arbitrary HTTP method and path.
func (c *RESTClient) Do(ctx context.Context, method, resourcePath string, body any, reply any) error {
	if err := requireContext(ctx); err != nil {
		return err
	}
	if method == "" {
		return fmt.Errorf("empty method")
	}

	rawBody, err := encodeRequestBody(body)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, c.resolvePath(resourcePath), bytes.NewReader(rawBody))
	if err != nil {
		return err
	}
	if len(rawBody) > 0 {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		if len(snippet) == 0 {
			return fmt.Errorf("rest request failed: %s", resp.Status)
		}
		return fmt.Errorf("rest request failed: %s: %s", resp.Status, strings.TrimSpace(string(snippet)))
	}

	if reply == nil {
		return nil
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(reply); err != nil {
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

// Get issues a GET request.
func (c *RESTClient) Get(ctx context.Context, resourcePath string, reply any) error {
	return c.Do(ctx, http.MethodGet, resourcePath, nil, reply)
}

// Post issues a POST request.
func (c *RESTClient) Post(ctx context.Context, resourcePath string, body any, reply any) error {
	return c.Do(ctx, http.MethodPost, resourcePath, body, reply)
}

// Put issues a PUT request.
func (c *RESTClient) Put(ctx context.Context, resourcePath string, body any, reply any) error {
	return c.Do(ctx, http.MethodPut, resourcePath, body, reply)
}

// Patch issues a PATCH request.
func (c *RESTClient) Patch(ctx context.Context, resourcePath string, body any, reply any) error {
	return c.Do(ctx, http.MethodPatch, resourcePath, body, reply)
}

// Delete issues a DELETE request.
func (c *RESTClient) Delete(ctx context.Context, resourcePath string, reply any) error {
	return c.Do(ctx, http.MethodDelete, resourcePath, nil, reply)
}

// Resource returns a path-oriented helper rooted at the provided segments.
func (c *RESTClient) Resource(parts ...string) *Resource {
	return &Resource{
		client: c,
		path:   joinPath(parts...),
	}
}

func (c *RESTClient) resolvePath(resourcePath string) string {
	relative := strings.TrimSpace(resourcePath)
	if relative == "" {
		return c.baseURL.String()
	}

	parsed, err := url.Parse(relative)
	if err != nil {
		parsed = &url.URL{Path: relative}
	}

	out := *c.baseURL
	out.Path = path.Join(out.Path, parsed.Path)
	if strings.HasSuffix(relative, "/") && !strings.HasSuffix(out.Path, "/") {
		out.Path += "/"
	}
	out.RawQuery = parsed.RawQuery
	out.Fragment = parsed.Fragment
	return out.String()
}

func encodeRequestBody(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	switch v := body.(type) {
	case []byte:
		return append([]byte(nil), v...), nil
	case json.RawMessage:
		return append(json.RawMessage(nil), v...), nil
	case string:
		return []byte(v), nil
	case io.Reader:
		return io.ReadAll(v)
	default:
		return json.Marshal(body)
	}
}

func normalizeHTTPClient(client *http.Client) *http.Client {
	if client != nil {
		return client
	}
	return &http.Client{
		Timeout: 15 * time.Second,
	}
}

// Resource is a path-scoped REST helper.
type Resource struct {
	client *RESTClient
	path   string
}

// Child returns a nested resource under the current path.
func (r *Resource) Child(parts ...string) *Resource {
	if r == nil {
		return nil
	}

	return &Resource{
		client: r.client,
		path:   joinPath(append([]string{r.path}, parts...)...),
	}
}

// Do issues a request against the resource path.
func (r *Resource) Do(ctx context.Context, method string, body any, reply any) error {
	if r == nil || r.client == nil {
		return fmt.Errorf("resource client is not configured")
	}
	return r.client.Do(ctx, method, r.path, body, reply)
}

// Get issues a GET request against the resource path.
func (r *Resource) Get(ctx context.Context, reply any) error {
	return r.Do(ctx, http.MethodGet, nil, reply)
}

// Post issues a POST request against the resource path.
func (r *Resource) Post(ctx context.Context, body any, reply any) error {
	return r.Do(ctx, http.MethodPost, body, reply)
}

// Put issues a PUT request against the resource path.
func (r *Resource) Put(ctx context.Context, body any, reply any) error {
	return r.Do(ctx, http.MethodPut, body, reply)
}

// Patch issues a PATCH request against the resource path.
func (r *Resource) Patch(ctx context.Context, body any, reply any) error {
	return r.Do(ctx, http.MethodPatch, body, reply)
}

// Delete issues a DELETE request against the resource path.
func (r *Resource) Delete(ctx context.Context, reply any) error {
	return r.Do(ctx, http.MethodDelete, nil, reply)
}

func joinPath(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}

	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			segments = append(segments, strings.Trim(trimmed, "/"))
		}
	}
	if len(segments) == 0 {
		return ""
	}

	return "/" + path.Join(segments...)
}
