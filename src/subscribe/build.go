package subscribe

import (
	"fmt"

	"sing-box-sub/parser"
	"sing-box-sub/template"
)

// BuildSingBoxConfig fetches subscriptions, parses proxies, applies prefixes
// (when multiple subs), loads the template, and returns the final sing-box config.
func BuildSingBoxConfig(subURLs []string, tmplPath string) (map[string]interface{}, error) {
	addPrefix := len(subURLs) > 1
	var allProxies []map[string]interface{}

	for i, url := range subURLs {
		url = trimSpace(url)
		if url == "" {
			continue
		}

		proxyText, err := GetProxyList(url, "")
		if err != nil {
			return nil, fmt.Errorf("failed to fetch subscription %s: %w", url, err)
		}

		proxies, err := parser.ParseProxies(proxyText)
		if err != nil {
			return nil, fmt.Errorf("failed to parse subscription %s: %w", url, err)
		}

		if addPrefix {
			prefix := fmt.Sprintf("%02d-", i+1)
			for _, p := range proxies {
				if tag, ok := p["tag"].(string); ok {
					p["tag"] = prefix + tag
				}
			}
		}

		allProxies = append(allProxies, proxies...)
	}

	if len(allProxies) == 0 {
		return nil, fmt.Errorf("no proxies found in any subscription")
	}

	if tmplPath == "" {
		tmplPath = "templates/momo.json"
	}

	tmpl, err := template.GetTemplate(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	return template.GenerateSingBoxConfigFromTemplate(tmpl, allProxies)
}

func trimSpace(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\r' || s[i] == '\n') {
		i++
	}
	j := len(s) - 1
	for j >= i && (s[j] == ' ' || s[j] == '\t' || s[j] == '\r' || s[j] == '\n') {
		j--
	}
	return s[i : j+1]
}
