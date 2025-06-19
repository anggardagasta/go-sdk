package vault

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/anggardagasta/go-sdk/vaultremote"
	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// nolint: gochecknoinits
func init() {
	vaultremote.RegisterConfigProvider("vault", NewConfigProvider())
}

// ConfigProvider implements reads configuration from Hashicorp Vault.
type ConfigProvider struct {
	clients map[string]*api.Client
}

// NewConfigProvider returns a new ConfigProvider.
func NewConfigProvider() *ConfigProvider {
	return &ConfigProvider{
		clients: make(map[string]*api.Client),
	}
}

func (p ConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	vaultPath := rp.Path()
	endpoint := rp.Endpoint()
	client, ok := p.clients[endpoint]
	mainPath := os.Getenv("VAULT_MAIN_PATH")
	if !ok {
		vaultRole := os.Getenv("VAULT_ROLE")
		jwtFile := os.Getenv("JWT_TOKEN_PATH")

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}

		httpTransport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}

		clientHttp := &http.Client{
			Transport: httpTransport,
		}

		config := api.DefaultConfig()
		config.HttpClient = clientHttp

		c, err := api.NewClient(config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create vault api client")
		}

		k8sAuth, err := auth.NewKubernetesAuth(
			vaultRole,
			auth.WithServiceAccountTokenPath(jwtFile),
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to initialize Kubernetes auth method")
		}

		authInfo, err := c.Auth().Login(context.Background(), k8sAuth)
		if err != nil {
			return nil, errors.Wrap(err, "unable to log in with Kubernetes auth")
		}
		if authInfo == nil {
			return nil, errors.Wrap(err, "no auth info was returned after login")
		}

		client = c
		p.clients[endpoint] = c
	}
	secret, err := client.KVv2(mainPath).Get(context.Background(), vaultPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read secret")
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.Errorf("source not found: %s", rp.Path())
	}

	// If it's a KV secrets engine, get the data from the "data" field.
	if data, ok := secret.Data["data"]; ok {
		secret.Data = data.(map[string]interface{})
	}

	b, err := json.Marshal(secret.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to json encode secret")
	}

	return bytes.NewReader(b), nil
}

func (p ConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	return nil, errors.New("watch is not implemented for the vault config provider")
}

func (p ConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	panic("watch channel is not implemented for the vault config provider")
}
