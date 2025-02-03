package beta

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/nullable"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EducationResource = EducationFileResource{}

type EducationFileResource struct {
	// Location on disk of the file resource.
	FileUrl nullable.Type[string] `json:"fileUrl,omitempty"`

	// Fields inherited from EducationResource

	// Who created the resource?
	CreatedBy *IdentitySet `json:"createdBy,omitempty"`

	// The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example,
	// midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
	CreatedDateTime nullable.Type[string] `json:"createdDateTime,omitempty"`

	// Display name of resource.
	DisplayName nullable.Type[string] `json:"displayName,omitempty"`

	// Who was the last user to modify the resource?
	LastModifiedBy *IdentitySet `json:"lastModifiedBy,omitempty"`

	// Moment in time when the resource was last modified. The Timestamp type represents date and time information using ISO
	// 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
	LastModifiedDateTime nullable.Type[string] `json:"lastModifiedDateTime,omitempty"`

	// The OData ID of this entity
	ODataId *string `json:"@odata.id,omitempty"`

	// The OData Type of this entity
	ODataType *string `json:"@odata.type,omitempty"`

	// Model Behaviors
	OmitDiscriminatedValue bool `json:"-"`
}

func (s EducationFileResource) EducationResource() BaseEducationResourceImpl {
	return BaseEducationResourceImpl{
		CreatedBy:            s.CreatedBy,
		CreatedDateTime:      s.CreatedDateTime,
		DisplayName:          s.DisplayName,
		LastModifiedBy:       s.LastModifiedBy,
		LastModifiedDateTime: s.LastModifiedDateTime,
		ODataId:              s.ODataId,
		ODataType:            s.ODataType,
	}
}

var _ json.Marshaler = EducationFileResource{}

func (s EducationFileResource) MarshalJSON() ([]byte, error) {
	type wrapper EducationFileResource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EducationFileResource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EducationFileResource: %+v", err)
	}

	if !s.OmitDiscriminatedValue {
		decoded["@odata.type"] = "#microsoft.graph.educationFileResource"
	}

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EducationFileResource: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &EducationFileResource{}

func (s *EducationFileResource) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		FileUrl              nullable.Type[string] `json:"fileUrl,omitempty"`
		CreatedDateTime      nullable.Type[string] `json:"createdDateTime,omitempty"`
		DisplayName          nullable.Type[string] `json:"displayName,omitempty"`
		LastModifiedDateTime nullable.Type[string] `json:"lastModifiedDateTime,omitempty"`
		ODataId              *string               `json:"@odata.id,omitempty"`
		ODataType            *string               `json:"@odata.type,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.FileUrl = decoded.FileUrl
	s.CreatedDateTime = decoded.CreatedDateTime
	s.DisplayName = decoded.DisplayName
	s.LastModifiedDateTime = decoded.LastModifiedDateTime
	s.ODataId = decoded.ODataId
	s.ODataType = decoded.ODataType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling EducationFileResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["createdBy"]; ok {
		impl, err := UnmarshalIdentitySetImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CreatedBy' for 'EducationFileResource': %+v", err)
		}
		s.CreatedBy = &impl
	}

	if v, ok := temp["lastModifiedBy"]; ok {
		impl, err := UnmarshalIdentitySetImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'LastModifiedBy' for 'EducationFileResource': %+v", err)
		}
		s.LastModifiedBy = &impl
	}

	return nil
}
