package main

import (
	"log"
	"os"

	"github.com/darkedges/vault-auth-plugin-guacamole/guacamole"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	if err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: guacamole.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	}); err != nil {
		log.Fatal(err)
	}
}
