package beta

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PrivilegeEscalationFinding = PrivilegeEscalationUserFinding{}

type PrivilegeEscalationUserFinding struct {

	// Fields inherited from PrivilegeEscalationFinding

	Identity *AuthorizationSystemIdentity `json:"identity,omitempty"`

	// An identity's information details. Inherited from finding.
	IdentityDetails *IdentityDetails `json:"identityDetails,omitempty"`

	PermissionsCreepIndex *PermissionsCreepIndex `json:"permissionsCreepIndex,omitempty"`

	// The list of escalations that the identity is capable of performing.
	PrivilegeEscalationDetails *[]PrivilegeEscalation `json:"privilegeEscalationDetails,omitempty"`

	// Fields inherited from Finding

	// Defines when the finding was created.
	CreatedDateTime *string `json:"createdDateTime,omitempty"`

	// Fields inherited from Entity

	// The unique identifier for an entity. Read-only.
	Id *string `json:"id,omitempty"`

	// The OData ID of this entity
	ODataId *string `json:"@odata.id,omitempty"`

	// The OData Type of this entity
	ODataType *string `json:"@odata.type,omitempty"`

	// Model Behaviors
	OmitDiscriminatedValue bool `json:"-"`
}

func (s PrivilegeEscalationUserFinding) PrivilegeEscalationFinding() BasePrivilegeEscalationFindingImpl {
	return BasePrivilegeEscalationFindingImpl{
		Identity:                   s.Identity,
		IdentityDetails:            s.IdentityDetails,
		PermissionsCreepIndex:      s.PermissionsCreepIndex,
		PrivilegeEscalationDetails: s.PrivilegeEscalationDetails,
		CreatedDateTime:            s.CreatedDateTime,
		Id:                         s.Id,
		ODataId:                    s.ODataId,
		ODataType:                  s.ODataType,
	}
}

func (s PrivilegeEscalationUserFinding) Finding() BaseFindingImpl {
	return BaseFindingImpl{
		CreatedDateTime: s.CreatedDateTime,
		Id:              s.Id,
		ODataId:         s.ODataId,
		ODataType:       s.ODataType,
	}
}

func (s PrivilegeEscalationUserFinding) Entity() BaseEntityImpl {
	return BaseEntityImpl{
		Id:        s.Id,
		ODataId:   s.ODataId,
		ODataType: s.ODataType,
	}
}

var _ json.Marshaler = PrivilegeEscalationUserFinding{}

func (s PrivilegeEscalationUserFinding) MarshalJSON() ([]byte, error) {
	type wrapper PrivilegeEscalationUserFinding
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PrivilegeEscalationUserFinding: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PrivilegeEscalationUserFinding: %+v", err)
	}

	if !s.OmitDiscriminatedValue {
		decoded["@odata.type"] = "#microsoft.graph.privilegeEscalationUserFinding"
	}

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PrivilegeEscalationUserFinding: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &PrivilegeEscalationUserFinding{}

func (s *PrivilegeEscalationUserFinding) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		IdentityDetails            *IdentityDetails       `json:"identityDetails,omitempty"`
		PermissionsCreepIndex      *PermissionsCreepIndex `json:"permissionsCreepIndex,omitempty"`
		PrivilegeEscalationDetails *[]PrivilegeEscalation `json:"privilegeEscalationDetails,omitempty"`
		CreatedDateTime            *string                `json:"createdDateTime,omitempty"`
		Id                         *string                `json:"id,omitempty"`
		ODataId                    *string                `json:"@odata.id,omitempty"`
		ODataType                  *string                `json:"@odata.type,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.CreatedDateTime = decoded.CreatedDateTime
	s.Id = decoded.Id
	s.IdentityDetails = decoded.IdentityDetails
	s.ODataId = decoded.ODataId
	s.ODataType = decoded.ODataType
	s.PermissionsCreepIndex = decoded.PermissionsCreepIndex
	s.PrivilegeEscalationDetails = decoded.PrivilegeEscalationDetails

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling PrivilegeEscalationUserFinding into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["identity"]; ok {
		impl, err := UnmarshalAuthorizationSystemIdentityImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Identity' for 'PrivilegeEscalationUserFinding': %+v", err)
		}
		s.Identity = &impl
	}

	return nil
}
