package parser

// convertAnytls converts a clash-meta anytls proxy to sing-box anytls outbound format.
func convertAnytls(proxy map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"type": "anytls",
	}

	if name := getString(proxy, "name"); name != "" {
		result["tag"] = name
	}
	if server := getString(proxy, "server"); server != "" {
		result["server"] = server
	}
	if port := getInt(proxy, "port"); port != 0 {
		result["server_port"] = port
	}
	if password := getString(proxy, "password"); password != "" {
		result["password"] = password
	}

	tls := buildTLS(proxy)
	tls["enabled"] = true
	if tls["server_name"] == nil || tls["server_name"] == "" {
		if server := getString(proxy, "server"); server != "" {
			tls["server_name"] = server
		}
	}
	if tls["insecure"] == nil {
		tls["insecure"] = true
	}

	if fp := getString(proxy, "client-fingerprint"); fp != "" {
		tls["utls"] = map[string]interface{}{
			"enabled":     true,
			"fingerprint": fp,
		}
	}

	result["tls"] = tls

	return result, nil
}
