package parser

// convertHysteria2 converts a clash-meta hysteria2 proxy to sing-box hysteria2 outbound format.
func convertHysteria2(proxy map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"type": "hysteria2",
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

	password := getString(proxy, "password")
	if password == "" {
		password = getString(proxy, "auth")
	}
	if password != "" {
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
	result["tls"] = tls

	obfsType := getString(proxy, "obfs")
	if obfsType != "" {
		obfs := map[string]interface{}{
			"type": obfsType,
		}
		if obfsPass := getString(proxy, "obfs-password"); obfsPass != "" {
			obfs["password"] = obfsPass
		}
		result["obfs"] = obfs
	}

	return result, nil
}
