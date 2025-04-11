package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/khostya/pvz/internal/config"
	"github.com/swaggo/swag"
	"os"
	"strings"
)

func GetOpenapiV3(ctx context.Context) (*openapi3.T, error) {
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	return loader.LoadFromFile("./api/v1/swagger/swagger.yaml")
}

const defaultApiHost = "localhost"

func init() {
	apiHost := os.Getenv("API_HOST")
	if apiHost == "" {
		apiHost = defaultApiHost
	}

	cfg := config.MustNewConfig()

	openapi, err := GetOpenapiV3(context.Background())
	if err != nil {
		panic(err)
	}

	openapi.AddServer(
		&openapi3.Server{
			URL: fmt.Sprintf("http://%v:%v", apiHost, cfg.API.HTTP.Port),
		},
	)

	file, err := openapi.MarshalJSON()
	if err != nil {
		panic(err)
	}

	template := string(file)

	m := make(map[string]any)
	err = json.NewDecoder(strings.NewReader(template)).Decode(&m)
	if err != nil {
		panic(err)
	}

	m["basePath"] = "/"
	m["host"] = fmt.Sprintf("http://localhost:%v", cfg.API.HTTP.Port)

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(m)
	if err != nil {
		panic(err)
	}

	swagger := &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  buf.String(),
	}

	swag.Register(swagger.InstanceName(), swagger)
}
