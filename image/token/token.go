package token

import (
	"context"
	"github.com/google/go-containerregistry/pkg/authn"

	"github.com/wadeling/registry-client/types"
)

var (
	registries []Registry
)

func init() {
	//RegisterRegistry(&gcr.GCR{})
}

type Registry interface {
	CheckOptions(domain string, option types.DockerOption) error
	GetCredential(ctx context.Context) (string, string, error)
}

func RegisterRegistry(registry Registry) {
	registries = append(registries, registry)
}

func GetToken(ctx context.Context, domain string, opt types.DockerOption) (auth authn.Basic) {
	if opt.UserName != "" || opt.Password != "" {
		return authn.Basic{Username: opt.UserName, Password: opt.Password}
	}

	// check registry which particular to get credential
	for _, registry := range registries {
		err := registry.CheckOptions(domain, opt)
		if err != nil {
			continue
		}
		username, password, err := registry.GetCredential(ctx)
		if err != nil {
			// only skip check registry if error occurred
			break
		}
		return authn.Basic{Username: username, Password: password}
	}
	return authn.Basic{}
}

