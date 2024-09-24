package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var providerInput *ElvidProviderInput

func convertSetToStringArray(set types.Set) []string {
	elements := set.Elements()
	if len(elements) == 0 {
		return []string{}
	}

	var result []string
	for _, v := range elements {
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

func ScopeApprovalWarning(id string) string {
	return "Missing approval for scopes for this client. Please ask in #core-henvendelser for approval. Check here for client details: " + providerInput.ElvIDAuthority + "/Configuration/ClientDetails?client_Id=" + id
}
