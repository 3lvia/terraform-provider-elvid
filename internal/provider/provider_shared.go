package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var providerInput *ElvidProviderInput

func convertSetToStringArray(set types.Set) []string {
	var result []string
	for _, v := range set.Elements() {
		result = append(result, v.(types.String).ValueString())
	}
	return result
}

func convertStringArrayToSet(array []string) types.Set {
	// Create a slice to hold the elements of the set
	elements := make([]attr.Value, len(array))

	// Convert each string in the array to a types.String and add it to the set
	for i, str := range array {
		elements[i] = types.StringValue(str)
	}

	// Return the Set containing the elements
	set, diag := types.SetValue(types.StringType, elements)
	if diag.HasError() {
		// Handle any diagnostic errors, possibly return an empty Set or log
		return types.Set{} // Return an empty set in case of an error
	}

	return set
}
