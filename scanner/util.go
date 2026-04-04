package scanner

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var shaRegex = regexp.MustCompile(`^[0-9a-f]{40}$`)

// isSHA returns true if the string looks like a full git SHA or docker digest.
func isSHA(s string) bool {
	return shaRegex.MatchString(s) || len(s) > 7 && s[:7] == "sha256:"
}

const bearerPrefix = "Bearer "

func mustCompile(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func isYAML(name string) bool {
	return strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml")
}

func slashDir(path string) string {
	return filepath.ToSlash(filepath.Dir(path))
}

func slashBase(path string) string {
	return filepath.Base(path)
}

// doWithRetry executes the request, retrying on 429 (rate limited) or 503
// up to maxRetries times, honouring Retry-After and X-RateLimit-Reset headers.
func doWithRetry(client *http.Client, req *http.Request, maxRetries int) (*http.Response, error) {
	for attempt := range maxRetries + 1 {
		// Clone the request on retries so the body can be re-sent if needed.
		r := req
		if attempt > 0 {
			clone, err := cloneRequest(req)
			if err != nil {
				return nil, err
			}
			r = clone
		}

		resp, err := client.Do(r)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusTooManyRequests && resp.StatusCode != http.StatusServiceUnavailable {
			return resp, nil
		}

		resp.Body.Close()

		if attempt == maxRetries {
			return nil, fmt.Errorf("rate limited after %d retries (HTTP %d)", maxRetries, resp.StatusCode)
		}

		delay := retryDelay(resp)
		time.Sleep(delay)
	}
	// unreachable
	return nil, fmt.Errorf("unexpected retry loop exit")
}

// retryDelay returns how long to wait before the next attempt.
// It reads Retry-After (seconds or HTTP-date) and X-RateLimit-Reset (unix timestamp).
func retryDelay(resp *http.Response) time.Duration {
	const fallback = 60 * time.Second

	if v := resp.Header.Get("Retry-After"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil {
			return time.Duration(secs) * time.Second
		}
		if t, err := http.ParseTime(v); err == nil {
			if d := time.Until(t); d > 0 {
				return d
			}
		}
	}

	if v := resp.Header.Get("X-RateLimit-Reset"); v != "" {
		if unix, err := strconv.ParseInt(v, 10, 64); err == nil {
			if d := time.Until(time.Unix(unix, 0)); d > 0 {
				return d
			}
		}
	}

	return fallback
}

// cloneRequest creates a shallow clone of req suitable for retrying a GET/HEAD
// (no body). For requests with bodies, callers must handle re-reading themselves.
func cloneRequest(req *http.Request) (*http.Request, error) {
	clone := req.Clone(req.Context())
	return clone, nil
}

