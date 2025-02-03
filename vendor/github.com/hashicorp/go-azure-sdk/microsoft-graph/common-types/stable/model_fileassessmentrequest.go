package stable

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/nullable"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ThreatAssessmentRequest = FileAssessmentRequest{}

type FileAssessmentRequest struct {
	// Base64 encoded file content. The file content can't fetch back because it isn't stored.
	ContentData *string `json:"contentData,omitempty"`

	// The file name.
	FileName *string `json:"fileName,omitempty"`

	// Fields inherited from ThreatAssessmentRequest

	Category *ThreatCategory `json:"category,omitempty"`

	// The content type of threat assessment. Possible values are: mail, url, file.
	ContentType *ThreatAssessmentContentType `json:"contentType,omitempty"`

	// The threat assessment request creator.
	CreatedBy IdentitySet `json:"createdBy"`

	// The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example,
	// midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
	CreatedDateTime nullable.Type[string] `json:"createdDateTime,omitempty"`

	ExpectedAssessment *ThreatExpectedAssessment `json:"expectedAssessment,omitempty"`

	// The source of the threat assessment request. Possible values are: administrator.
	RequestSource *ThreatAssessmentRequestSource `json:"requestSource,omitempty"`

	// A collection of threat assessment results. Read-only. By default, a GET /threatAssessmentRequests/{id} does not
	// return this property unless you apply $expand on it.
	Results *[]ThreatAssessmentResult `json:"results,omitempty"`

	// The assessment process status. Possible values are: pending, completed.
	Status *ThreatAssessmentStatus `json:"status,omitempty"`

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

func (s FileAssessmentRequest) ThreatAssessmentRequest() BaseThreatAssessmentRequestImpl {
	return BaseThreatAssessmentRequestImpl{
		Category:           s.Category,
		ContentType:        s.ContentType,
		CreatedBy:          s.CreatedBy,
		CreatedDateTime:    s.CreatedDateTime,
		ExpectedAssessment: s.ExpectedAssessment,
		RequestSource:      s.RequestSource,
		Results:            s.Results,
		Status:             s.Status,
		Id:                 s.Id,
		ODataId:            s.ODataId,
		ODataType:          s.ODataType,
	}
}

func (s FileAssessmentRequest) Entity() BaseEntityImpl {
	return BaseEntityImpl{
		Id:        s.Id,
		ODataId:   s.ODataId,
		ODataType: s.ODataType,
	}
}

var _ json.Marshaler = FileAssessmentRequest{}

func (s FileAssessmentRequest) MarshalJSON() ([]byte, error) {
	type wrapper FileAssessmentRequest
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FileAssessmentRequest: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FileAssessmentRequest: %+v", err)
	}

	if !s.OmitDiscriminatedValue {
		decoded["@odata.type"] = "#microsoft.graph.fileAssessmentRequest"
	}

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FileAssessmentRequest: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &FileAssessmentRequest{}

func (s *FileAssessmentRequest) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ContentData        *string                        `json:"contentData,omitempty"`
		FileName           *string                        `json:"fileName,omitempty"`
		Category           *ThreatCategory                `json:"category,omitempty"`
		ContentType        *ThreatAssessmentContentType   `json:"contentType,omitempty"`
		CreatedDateTime    nullable.Type[string]          `json:"createdDateTime,omitempty"`
		ExpectedAssessment *ThreatExpectedAssessment      `json:"expectedAssessment,omitempty"`
		RequestSource      *ThreatAssessmentRequestSource `json:"requestSource,omitempty"`
		Results            *[]ThreatAssessmentResult      `json:"results,omitempty"`
		Status             *ThreatAssessmentStatus        `json:"status,omitempty"`
		Id                 *string                        `json:"id,omitempty"`
		ODataId            *string                        `json:"@odata.id,omitempty"`
		ODataType          *string                        `json:"@odata.type,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ContentData = decoded.ContentData
	s.FileName = decoded.FileName
	s.Category = decoded.Category
	s.ContentType = decoded.ContentType
	s.CreatedDateTime = decoded.CreatedDateTime
	s.ExpectedAssessment = decoded.ExpectedAssessment
	s.Id = decoded.Id
	s.ODataId = decoded.ODataId
	s.ODataType = decoded.ODataType
	s.RequestSource = decoded.RequestSource
	s.Results = decoded.Results
	s.Status = decoded.Status

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FileAssessmentRequest into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["createdBy"]; ok {
		impl, err := UnmarshalIdentitySetImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CreatedBy' for 'FileAssessmentRequest': %+v", err)
		}
		s.CreatedBy = impl
	}

	return nil
}
