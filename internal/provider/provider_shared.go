package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

var providerInput *ElvidProviderInput

func convertSetToStringArray(set types.Set) []string {
	var result []string
	for _, v := range set.Elements() {
		result = append(result, v.(types.String).ValueString())
	}
	return result
}
