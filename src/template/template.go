package template

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// GetTemplate loads a sing-box template from a local file path or a remote URL.
func GetTemplate(tmplPath string) (map[string]interface{}, error) {
	var data []byte
	var err error

	if strings.HasPrefix(tmplPath, "http://") || strings.HasPrefix(tmplPath, "https://") {
		data, err = fetchTemplate(tmplPath)
	} else {
		data, err = os.ReadFile(tmplPath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse template JSON: %w", err)
	}

	return result, nil
}

func fetchTemplate(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// GenerateSingBoxConfigFromTemplate inserts the parsed proxy list into the template
// and returns the final sing-box configuration.
func GenerateSingBoxConfigFromTemplate(tmpl map[string]interface{}, proxies []map[string]interface{}) (map[string]interface{}, error) {
	// Apply global nodeFilter to exclude unwanted nodes entirely (e.g. promo/fake nodes).
	if nodeFilterVal, ok := tmpl["nodeFilter"]; ok {
		if filterRules, ok := nodeFilterVal.([]interface{}); ok {
			proxies = applyNodeFilter(proxies, filterRules)
		}
	}
	delete(tmpl, "nodeFilter")

	allTags := make([]string, 0, len(proxies))
	for _, p := range proxies {
		tag, ok := p["tag"].(string)
		if !ok || tag == "" {
			continue
		}
		allTags = append(allTags, tag)
	}

	outbounds, ok := tmpl["outbounds"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("template outbounds not found or invalid")
	}

	// Process each outbound: expand {all} placeholder and apply filters.
	for i, ob := range outbounds {
		obMap, ok := ob.(map[string]interface{})
		if !ok {
			continue
		}

		if isAllPlaceholder(obMap["outbounds"]) {
			filteredTags := allTags
			if filterVal, ok := obMap["filter"]; ok {
				if filterRules, ok := filterVal.([]interface{}); ok {
					filteredTags = applyFilter(allTags, filterRules)
				}
			}
			tagInterfaces := make([]interface{}, len(filteredTags))
			for j, t := range filteredTags {
				tagInterfaces[j] = t
			}
			obMap["outbounds"] = tagInterfaces
		}

		outbounds[i] = obMap
	}

	// Remove "filter" from all outbounds — filter is a template-only field,
	// not recognized by sing-box and will cause a fatal config error.
	for i, ob := range outbounds {
		if obMap, ok := ob.(map[string]interface{}); ok {
			delete(obMap, "filter")
			outbounds[i] = obMap
		}
	}

	// Append each proxy as a new outbound node.
	for _, p := range proxies {
		outbounds = append(outbounds, p)
	}
	tmpl["outbounds"] = outbounds

	return tmpl, nil
}

// applyNodeFilter filters proxy nodes by their "tag" field using include/exclude rules.
// Nodes that don't pass the filter are removed entirely from the proxies list.
func applyNodeFilter(proxies []map[string]interface{}, rules []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(proxies))
	copy(result, proxies)

	for _, rule := range rules {
		ruleMap, ok := rule.(map[string]interface{})
		if !ok {
			continue
		}

		action := getString(ruleMap, "action")
		keywords := getStringList(ruleMap, "keywords")
		if len(keywords) == 0 {
			continue
		}

		pattern := strings.Join(keywords, "|")
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		var filtered []map[string]interface{}
		for _, p := range result {
			tag, _ := p["tag"].(string)
			matched := re.MatchString(tag)
			switch action {
			case "include":
				if matched {
					filtered = append(filtered, p)
				}
			case "exclude":
				if !matched {
					filtered = append(filtered, p)
				}
			}
		}
		result = filtered
	}

	return result
}

// applyFilter filters a list of proxy tag strings based on include/exclude rules.
// Rules are executed in order. Each rule has:
//   - "action": "include" or "exclude"
//   - "keywords": []string, joined with "|" to form a regex pattern
//   - "for" (optional): list of subscription tags to scope the rule to (ignored in single-sub mode)
func applyFilter(tags []string, rules []interface{}) []string {
	result := make([]string, len(tags))
	copy(result, tags)

	for _, rule := range rules {
		ruleMap, ok := rule.(map[string]interface{})
		if !ok {
			continue
		}

		action := getString(ruleMap, "action")
		keywords := getStringList(ruleMap, "keywords")
		if len(keywords) == 0 {
			continue
		}

		pattern := strings.Join(keywords, "|")
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		var filtered []string
		for _, tag := range result {
			matched := re.MatchString(tag)
			switch action {
			case "include":
				if matched {
					filtered = append(filtered, tag)
				}
			case "exclude":
				if !matched {
					filtered = append(filtered, tag)
				}
			}
		}
		result = filtered
	}

	return result
}

func isAllPlaceholder(v interface{}) bool {
	arr, ok := v.([]interface{})
	if !ok || len(arr) != 1 {
		return false
	}
	s, ok := arr[0].(string)
	return ok && s == "{all}"
}

func getString(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func getStringList(m map[string]interface{}, key string) []string {
	v, ok := m[key]
	if !ok {
		return nil
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	var result []string
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
