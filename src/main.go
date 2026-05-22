package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	server "sing-box-sub/http"
	"sing-box-sub/subscribe"
)

var buildVersion = "dev"

func main() {
	subURL := flag.String("u", "", "Subscription URL(s), comma-separated (max 10)")
	subURLShort := flag.String("url", "", "Subscription URL(s), comma-separated (max 10)")
	tmplPath := flag.String("t", "", "Template file path (local or URL)")
	tmplPathShort := flag.String("template", "", "Template file path (local or URL)")
	output := flag.String("o", "", "Output file path (default: stdout)")
	outputShort := flag.String("output", "", "Output file path (default: stdout)")
	daemon := flag.Bool("d", false, "Run as HTTP server")
	help := flag.Bool("h", false, "Show help")
	helpShort := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help || *helpShort {
		printHelp()
		return
	}

	if *daemon {
		RunAsDaemon()
		return
	}

	// Resolve short/long flag variants
	urls := *subURL
	if urls == "" {
		urls = *subURLShort
	}
	tmpl := *tmplPath
	if tmpl == "" {
		tmpl = *tmplPathShort
	}
	out := *output
	if out == "" {
		out = *outputShort
	}

	if urls == "" && flag.NArg() > 0 {
		urls = flag.Arg(0)
	}

	if urls == "" {
		fmt.Fprintln(os.Stderr, "Error: subscription URL is required")
		printHelp()
		os.Exit(1)
	}

	RunAsParser(urls, tmpl, out)
}

func RunAsParser(subURLs string, tmplPath string, outputPath string) {
	urlList := strings.Split(subURLs, ",")
	if len(urlList) > 10 {
		log("warning: maximum 10 subscription URLs allowed, using first 10")
		urlList = urlList[:10]
	}

	log("building config from %d subscription(s)...", len(urlList))
	res, err := subscribe.BuildSingBoxConfig(urlList, tmplPath)
	if err != nil {
		log("error: %v", err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			log("error creating output file: %v", err)
			os.Exit(1)
		}
		defer f.Close()
		encoder = json.NewEncoder(f)
		encoder.SetIndent("", "  ")
		encoder.Encode(res)
		log("config written to %s", outputPath)
	} else {
		encoder.Encode(res)
	}
}

func RunAsDaemon() {
	listenAddr := "0.0.0.0:40533"
	if env := os.Getenv("LISTEN_ADDR"); env != "" {
		listenAddr = env
	}
	log("starting HTTP server on %s", listenAddr)
	if err := server.StartHttpServer(listenAddr); err != nil {
		log("server error: %v", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("sing-box-sub - sing-box subscription converter")
	fmt.Println()
	fmt.Println("Usage 1 (CLI):")
	fmt.Println("  sing-box-sub -u <subscription URL> -t <template path> -o <output file>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println("    -u, -url <string>       Subscription URL(s), comma-separated (max 10)")
	fmt.Println("    -t, -template <string>  Template file path (local file or URL)")
	fmt.Println("    -o, -output <string>    Output file path (default: stdout)")
	fmt.Println("    -h, -help               Show help")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    sing-box-sub -u https://example.com/sub -t https://example.com/template.json -o config.json")
	fmt.Println("    sing-box-sub -u https://example.com/sub -t /templates/default.json")
	fmt.Println()
	fmt.Println("Usage 2 (HTTP server):")
	fmt.Println("  sing-box-sub -d")
	fmt.Println()
	fmt.Println("  Endpoints:")
	fmt.Println("    GET /convert?sub=<url1>&sub=<url2>&tpl=<template>")
	fmt.Println("    GET /health")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    http://127.0.0.1:40533/convert?sub=https://sub1.com&sub=https://sub2.com&tpl=https://tpl.json")
	fmt.Println("    http://127.0.0.1:40533/convert?sub=https://sub1.com?profile=simple&sub=https://sub2.com")
	fmt.Println()
}

func log(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[sing-box-sub] %s\n", fmt.Sprintf(format, args...))
}
