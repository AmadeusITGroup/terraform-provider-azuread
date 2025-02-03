package user

import (
	"github.com/hashicorp/go-azure-sdk/sdk/nullable"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WipeManagedAppRegistrationsByAzureAdDeviceIdRequest struct {
	AzureAdDeviceId nullable.Type[string] `json:"azureAdDeviceId,omitempty"`
}
