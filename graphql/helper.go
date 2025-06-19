package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// SchemaRegistry struct configuration for schema registry
type SchemaRegistry struct {
	ServiceName        string
	ServiceAddr        string
	Schema             string
	SchemaRegistryAddr string
	Version            string
}

// RegisterSchema push schema to registry schema server
func RegisterSchema(ctx context.Context, schemaCfg SchemaRegistry) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	if schemaCfg.SchemaRegistryAddr == "" {
		return errors.New("schema registry address is empty")
	}

	if schemaCfg.ServiceName == "" {
		return errors.New("service name is empty")
	}

	if schemaCfg.ServiceAddr == "" {
		return errors.New("service address is empty")
	}

	if schemaCfg.Schema == "" {
		return errors.New("schema is empty")
	}

	graphqlSchema := schemaCfg.Schema

	if schemaCfg.Version == "" {
		schemaCfg.Version = RandomStringUnsafe(5)
	}

	requestBody, err := json.Marshal(map[string]string{
		"name":      schemaCfg.ServiceName,
		"version":   schemaCfg.Version,
		"type_defs": graphqlSchema,
		"url":       schemaCfg.ServiceAddr,
	})

	if err != nil {
		return errors.Wrap(err, "error when marshal request body")
	}

	response, err := client.Post(
		fmt.Sprintf("%v/schema/push", schemaCfg.SchemaRegistryAddr),
		"application/json",
		bytes.NewBuffer(requestBody),
	)

	if err != nil {
		return errors.Wrap(err, "error when send post registry schema")
	}
	var res map[string]interface{}

	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		return errors.Wrap(err, "error when decode response body")
	}

	return nil
}

func RandomStringUnsafe(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
