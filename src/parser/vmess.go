package parser

// convertVmess converts a clash-meta vmess proxy to sing-box vmess outbound format.
func convertVmess(proxy map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"type": "vmess",
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
	if uuid := getString(proxy, "uuid"); uuid != "" {
		result["uuid"] = uuid
	}
	if alterID := getInt(proxy, "alterId"); alterID != 0 {
		result["alter_id"] = alterID
	}
	if cipher := getString(proxy, "cipher"); cipher != "" {
		result["security"] = cipher
	}

	network := getString(proxy, "network")
	if network != "" && network != "tcp" {
		result["transport"] = map[string]interface{}{
			"type": network,
		}
	}

	tlsEnabled := false
	if tlsVal, ok := proxy["tls"]; ok {
		switch v := tlsVal.(type) {
		case bool:
			tlsEnabled = v
		case string:
			tlsEnabled = v == "true" || v == "1"
		}
	}

	if tlsEnabled {
		tls := map[string]interface{}{
			"enabled": true,
		}
		if sni := getString(proxy, "sni"); sni != "" {
			tls["server_name"] = sni
		}
		if alpn := getList(proxy, "alpn"); len(alpn) > 0 {
			tls["alpn"] = alpn
		}
		if getBool(proxy, "skip-cert-verify") {
			tls["insecure"] = true
		}
		if fp := getString(proxy, "client-fingerprint"); fp != "" {
			tls["utls"] = map[string]interface{}{
				"enabled":     true,
				"fingerprint": fp,
			}
		}
		result["tls"] = tls
	}

	return result, nil
}
