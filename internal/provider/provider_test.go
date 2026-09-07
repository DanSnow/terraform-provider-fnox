package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate the provider
// during acceptance testing. The factory function is called for each
// Terraform CLI command executed to create a provider server that the CLI
// can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"fnox": providerserver.NewProtocol6WithError(New("test")()),
}
