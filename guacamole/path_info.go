package guacamole

import (
	"context"

	"github.com/darkedges/vault-auth-plugin-guacamole/version"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathInfoRead corresponds to READ auth/guacamole/info.
func (b *backend) pathInfoRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{
		Data: map[string]interface{}{
			"name":    version.Name,
			"version": version.Version,
		},
	}, nil
}
