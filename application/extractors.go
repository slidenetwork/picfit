package application

import (
	"fmt"
	"mime"
	"net/url"
	"path/filepath"

	"github.com/slidenetwork/picfit/engines"
	"github.com/slidenetwork/picfit/image"
)

type Extractor func(key string, req *Request) (interface{}, error)

var Operation Extractor = func(key string, req *Request) (interface{}, error) {
	operation, ok := engines.Operations[req.QueryString[key]]

	if !ok {
		operation = engines.Operations["original"]
	}

	return operation, nil
}

var URL Extractor = func(key string, req *Request) (interface{}, error) {
	value, ok := req.QueryString[key]

	if !ok {
		return nil, nil
	}

	url, err := url.Parse(value)

	if err != nil {
		return nil, fmt.Errorf("URL %s is not valid", value)
	}

	mimetype := mime.TypeByExtension(filepath.Ext(value))

	_, ok = image.Extensions[mimetype]

	if !ok {
		return nil, fmt.Errorf("Mimetype %s is not supported", mimetype)
	}

	return url, nil
}

var Path Extractor = func(key string, req *Request) (interface{}, error) {
	return req.QueryString[key], nil
}

var Extractors = map[string]Extractor{
	"op":   Operation,
	"url":  URL,
	"path": Path,
}
