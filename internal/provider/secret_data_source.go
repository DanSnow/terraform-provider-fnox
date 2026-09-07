package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &SecretDataSource{}
var _ datasource.DataSourceWithConfigure = &SecretDataSource{}

func NewSecretDataSource() datasource.DataSource {
	return &SecretDataSource{}
}

// SecretDataSource fetches a single secret via `fnox get`.
type SecretDataSource struct {
	client fnoxClientConfig
}

type secretDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Key          types.String `tfsdk:"key"`
	Profile      types.String `tfsdk:"profile"`
	Base64Decode types.Bool   `tfsdk:"base64_decode"`
	Value        types.String `tfsdk:"value"`
}

func (d *SecretDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (d *SecretDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single secret from fnox via `fnox get`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source, equal to key.",
			},
			"key": schema.StringAttribute{
				Required:    true,
				Description: "The fnox secret key to retrieve.",
			},
			"profile": schema.StringAttribute{
				Optional:    true,
				Description: "fnox profile to use for this secret, overriding the provider-level default profile.",
			},
			"base64_decode": schema.BoolAttribute{
				Optional:    true,
				Description: "Base64-decode the secret value after retrieval (passes --base64-decode to fnox get).",
			},
			"value": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The retrieved secret value.",
			},
		},
	}
}

func (d *SecretDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(fnoxClientConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected provider.fnoxClientConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *SecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data secretDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	value, err := d.client.getSecret(ctx, getSecretParams{
		Key:          data.Key.ValueString(),
		Profile:      data.Profile.ValueString(),
		Base64Decode: data.Base64Decode.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read fnox Secret", err.Error())
		return
	}

	data.ID = data.Key
	data.Value = types.StringValue(value)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
