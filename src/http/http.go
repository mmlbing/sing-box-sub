package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sing-box-sub/subscribe"
)

func StartHttpServer(listenAddr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", handleConvert)
	mux.HandleFunc("/convert/", handleConvert)
	mux.HandleFunc("/health", handleHealth)

	return http.ListenAndServe(listenAddr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /convert?sub=<url1>&sub=<url2>&tpl=<template>
//
// Examples:
//
//	http://127.0.0.1:40533/convert?sub=https://sub1.com&sub=https://sub2.com&tpl=https://tpl.json
//	http://127.0.0.1:40533/convert?sub=https://sub1.com?profile=simple&sub=https://sub2.com
//
// No URL encoding needed — Go's query parser handles ? and = within parameter values.
func handleConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "use GET", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	subURLs := query["sub"]
	if len(subURLs) == 0 {
		http.Error(w, "sub is required. Usage: /convert?sub=<url1>&sub=<url2>&tpl=<template>", http.StatusBadRequest)
		return
	}
	if len(subURLs) > 10 {
		http.Error(w, "maximum 10 subscription URLs", http.StatusBadRequest)
		return
	}

	tmplPath := query.Get("tpl")

	res, err := subscribe.BuildSingBoxConfig(subURLs, tmplPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(res)
}
