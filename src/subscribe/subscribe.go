package subscribe

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func FetchSub(src string, userAgent string) ([]byte, error) {
	if !strings.HasPrefix(src, "https://") && !strings.HasPrefix(src, "http://") {
		return nil, fmt.Errorf("unsupported subscribe path format: %s", src)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", src, nil)
	if err != nil {
		return nil, err
	}
	if userAgent == "" {
		userAgent = "clashmeta"
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// GetProxyListFromSubConfig extracts the raw "proxies:" section from a clash-meta YAML config.
// It returns the text content of all proxy entries (excluding the "proxies:" key line itself).
func GetProxyListFromSubConfig(sub string) (string, error) {
	lines := strings.Split(sub, "\n")
	inProxies := false
	var proxyLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		isTopLevel := len(line) > 0 && line[0] != ' ' && line[0] != '\t' && strings.HasSuffix(trimmed, ":")

		if isTopLevel {
			if trimmed == "proxies:" || strings.HasPrefix(trimmed, "proxies") {
				inProxies = true
				continue
			} else if inProxies {
				break
			}
		}

		if inProxies {
			proxyLines = append(proxyLines, line)
		}
	}

	// Remove empty leading/trailing lines
	for len(proxyLines) > 0 && strings.TrimSpace(proxyLines[0]) == "" {
		proxyLines = proxyLines[1:]
	}
	for len(proxyLines) > 0 && strings.TrimSpace(proxyLines[len(proxyLines)-1]) == "" {
		proxyLines = proxyLines[:len(proxyLines)-1]
	}

	if len(proxyLines) == 0 {
		return "", fmt.Errorf("no proxies section found in subscription config")
	}

	// Detect inline flow format: first non-empty line starts with whitespace + "- {"
	// In this case, strip the common leading whitespace so yaml.v3 can parse it.
	firstNonEmpty := ""
	for _, line := range proxyLines {
		if strings.TrimSpace(line) != "" {
			firstNonEmpty = line
			break
		}
	}

	if isInlineFlowFormat(firstNonEmpty) {
		// Strip the leading whitespace indent from all lines
		indent := len(firstNonEmpty) - len(strings.TrimLeft(firstNonEmpty, " \t"))
		for i, line := range proxyLines {
			if line == "" {
				continue
			}
			trimLen := indent
			if len(line) < trimLen {
				trimLen = len(line)
			}
			proxyLines[i] = line[trimLen:]
		}
	}

	result := strings.TrimSpace(strings.Join(proxyLines, "\n"))
	if result == "" {
		return "", fmt.Errorf("no proxies section found in subscription config")
	}

	return result, nil
}

func isInlineFlowFormat(firstLine string) bool {
	trimmed := strings.TrimSpace(firstLine)
	return strings.HasPrefix(trimmed, "- {") || strings.HasPrefix(trimmed, "-{")
}

// GetProxyList fetches a subscription and returns the raw proxies text.
func GetProxyList(src string, userAgent string) (string, error) {
	body, err := FetchSub(src, userAgent)
	if err != nil {
		return "", err
	}

	proxyList, err := GetProxyListFromSubConfig(string(body))
	if err != nil {
		return "", err
	}

	return proxyList, nil
}
