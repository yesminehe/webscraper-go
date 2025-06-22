package main

import "net/url"

func IsValidURL(u string) bool {
	parsed, err := url.ParseRequestURI(u)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
} 