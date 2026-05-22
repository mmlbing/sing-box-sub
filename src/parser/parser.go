package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// NormalizeProxies parses the raw proxies text (in clash-meta YAML format) into a list of maps.
// Supports both inline flow format (- { key: value, ... }) and block format.
func NormalizeProxies(proxies string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := yaml.Unmarshal([]byte(proxies), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxies: %w", err)
	}
	return result, nil
}

// ParseProxies converts a list of clash-meta proxies into sing-box outbound nodes.
func ParseProxies(proxies string) ([]map[string]interface{}, error) {
	normalized, err := NormalizeProxies(proxies)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, proxy := range normalized {
		proxyType := getString(proxy, "type")
		var converted map[string]interface{}
		var convErr error

		switch proxyType {
		case "hysteria2":
			converted, convErr = convertHysteria2(proxy)
		case "anytls":
			converted, convErr = convertAnytls(proxy)
		case "vmess":
			converted, convErr = convertVmess(proxy)
		default:
			continue
		}

		if convErr != nil {
			return nil, fmt.Errorf("failed to convert proxy %s: %w", getString(proxy, "name"), convErr)
		}
		result = append(result, converted)
	}

	return result, nil
}

func getString(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case int:
		return fmt.Sprintf("%d", val)
	default:
		return ""
	}
}

func getInt(m map[string]interface{}, key string) int {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	default:
		return 0
	}
}

func getBool(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	b, _ := v.(bool)
	return b
}

// buildTLS constructs the sing-box TLS object from clash-meta proxy fields.
func buildTLS(proxy map[string]interface{}) map[string]interface{} {
	tls := map[string]interface{}{}

	if sni := getString(proxy, "sni"); sni != "" {
		tls["server_name"] = sni
	}
	if alpn := getList(proxy, "alpn"); len(alpn) > 0 {
		tls["alpn"] = alpn
	}
	if getBool(proxy, "skip-cert-verify") {
		tls["insecure"] = true
	}

	return tls
}

func getList(m map[string]interface{}, key string) []interface{} {
	v, ok := m[key]
	if !ok {
		return nil
	}
	l, _ := v.([]interface{})
	return l
}
