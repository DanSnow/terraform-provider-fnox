package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &FnoxProvider{}

// FnoxProvider shells out to the fnox CLI to expose secrets managed by fnox
// as Terraform data sources.
type FnoxProvider struct {
	version string
}

// fnoxProviderModel maps the provider schema to Go types.
type fnoxProviderModel struct {
	BinaryPath types.String `tfsdk:"binary_path"`
	ConfigPath types.String `tfsdk:"config_path"`
	WorkingDir types.String `tfsdk:"working_dir"`
	Profile    types.String `tfsdk:"profile"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FnoxProvider{version: version}
	}
}

func (p *FnoxProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "fnox"
	resp.Version = p.version
}

func (p *FnoxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interacts with the fnox CLI (https://fnox.jdx.dev) to expose secrets from any fnox-configured backend as Terraform data sources.",
		Attributes: map[string]schema.Attribute{
			"binary_path": schema.StringAttribute{
				Optional:    true,
				Description: "Path to the fnox binary. Defaults to \"fnox\" resolved via $PATH. Can also be set via the FNOX_BINARY_PATH environment variable.",
			},
			"config_path": schema.StringAttribute{
				Optional:    true,
				Description: "Path to the fnox.toml configuration file. Defaults to fnox's own discovery (fnox.toml, searching parent directories). Can also be set via the FNOX_CONFIG environment variable.",
			},
			"working_dir": schema.StringAttribute{
				Optional:    true,
				Description: "Directory to run the fnox CLI in, used for relative configuration discovery.",
			},
			"profile": schema.StringAttribute{
				Optional:    true,
				Description: "Default fnox profile to use for all data sources. Can be overridden per data source. Can also be set via the FNOX_PROFILE environment variable.",
			},
		},
	}
}

func (p *FnoxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data fnoxProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := fnoxClientConfig{
		BinaryPath: data.BinaryPath.ValueString(),
		ConfigPath: data.ConfigPath.ValueString(),
		WorkingDir: data.WorkingDir.ValueString(),
		Profile:    data.Profile.ValueString(),
	}

	if config.BinaryPath == "" {
		config.BinaryPath = os.Getenv("FNOX_BINARY_PATH")
	}
	if config.ConfigPath == "" {
		config.ConfigPath = os.Getenv("FNOX_CONFIG")
	}
	if config.Profile == "" {
		config.Profile = os.Getenv("FNOX_PROFILE")
	}

	resp.DataSourceData = config
}

func (p *FnoxProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *FnoxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSecretDataSource,
	}
}
