package guacamole

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
)

// Factory creates a new usable instance of this auth method.
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(c)
	if err := b.Setup(ctx, c); err != nil {
		return nil, errors.Wrapf(err, "failed to create factory")
	}
	return b, nil
}

type backend struct {
	*framework.Backend
}

// Backend creates a new backend, mapping the proper paths, help information,
// and required callbacks.
func Backend(c *logical.BackendConfig) *backend {
	var b backend

	b.Backend = &framework.Backend{
		BackendType: logical.TypeCredential,
		AuthRenew:   b.pathAuthRenew,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{"login"},
		},
		Paths: func() []*framework.Path {
			var paths []*framework.Path
			paths = append(paths, &framework.Path{
				Pattern:      "info",
				HelpSynopsis: "Display information about the plugin",
				HelpDescription: `
Displays information about the plugin, such as the plugin version and where to
get help.
`,
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: b.pathInfoRead,
				},
			})

			paths = append(paths, &framework.Path{
				Pattern: "login",
				Fields: map[string]*framework.FieldSchema{
					"password": &framework.FieldSchema{
						Type: framework.TypeString,
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.pathAuthLogin,
				},
			})

			paths = append(paths, &framework.Path{
				Pattern:      "config",
				HelpSynopsis: "Configuration such the team and ttls",
				HelpDescription: `
Read or writer configuration to Vault's storage backend such as OAuth
information, team, behavior configuration tunables, and TTLs. For example:
    $ vault write auth/slack/config \
        access_token="xoxp-2164918114..." \
        teams="HashiCorp"
For more information and examples, please see the online documentation.
`,

				Fields: map[string]*framework.FieldSchema{
					"access_token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Slack OAuth access token for your Slack application.",
					},

					"teams": &framework.FieldSchema{
						Type: framework.TypeCommaStringSlice,
						Description: "Comma-separated list of permitted Slack teams. The " +
							"user must be a member of at least one of these teams to " +
							"authenticate.",
					},

					"allow_bot_users": &framework.FieldSchema{
						Type:        framework.TypeBool,
						Description: "Allow bot users to authenticate.",
					},

					"allow_non_2fa": &framework.FieldSchema{
						Type: framework.TypeBool,
						Description: "Allow users to not have 2FA enabled on their Slack " +
							"account to authenticate.",
						Default: true,
					},

					"allow_restricted_users": &framework.FieldSchema{
						Type: framework.TypeBool,
						Description: "Allow restricted users (multi-channel guests) to " +
							"authenticate.",
					},

					"allow_ultra_restricted_users": &framework.FieldSchema{
						Type: framework.TypeBool,
						Description: "Allow ultra restricted users (single-channel " +
							"guests) to authenticate.",
					},

					"anyone_policies": &framework.FieldSchema{
						Type: framework.TypeCommaStringSlice,
						Description: "Comma-separated list of policies to apply to " +
							"everyone, even unmapped users.",
					},

					"ttl": &framework.FieldSchema{
						Type:        framework.TypeDurationSecond,
						Description: "Duration after which authentication will expire.",
					},

					"max_ttl": &framework.FieldSchema{
						Type:        framework.TypeDurationSecond,
						Description: "Maximum duration after which authentication will expire.",
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.pathConfigWrite,
					logical.ReadOperation:   b.pathConfigRead,
				},
			})
			return paths
		}(),
	}

	return &b
}
