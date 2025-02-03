package stable

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId{}

// IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId is a struct representing the Resource ID for a Identity Governance Lifecycle Workflow Deleted Item Workflow Id Version Id Task
type IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId struct {
	WorkflowId                   string
	WorkflowVersionVersionNumber string
	TaskId                       string
}

// NewIdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskID returns a new IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId struct
func NewIdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskID(workflowId string, workflowVersionVersionNumber string, taskId string) IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId {
	return IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId{
		WorkflowId:                   workflowId,
		WorkflowVersionVersionNumber: workflowVersionVersionNumber,
		TaskId:                       taskId,
	}
}

// ParseIdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskID parses 'input' into a IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId
func ParseIdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskID(input string) (*IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseIdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskIDInsensitively parses 'input' case-insensitively into a IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId
// note: this method should only be used for API response data and not user input
func ParseIdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskIDInsensitively(input string) (*IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.WorkflowId, ok = input.Parsed["workflowId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workflowId", input)
	}

	if id.WorkflowVersionVersionNumber, ok = input.Parsed["workflowVersionVersionNumber"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workflowVersionVersionNumber", input)
	}

	if id.TaskId, ok = input.Parsed["taskId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "taskId", input)
	}

	return nil
}

// ValidateIdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskID checks that 'input' can be parsed as a Identity Governance Lifecycle Workflow Deleted Item Workflow Id Version Id Task ID
func ValidateIdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Identity Governance Lifecycle Workflow Deleted Item Workflow Id Version Id Task ID
func (id IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId) ID() string {
	fmtString := "/identityGovernance/lifecycleWorkflows/deletedItems/workflows/%s/versions/%s/tasks/%s"
	return fmt.Sprintf(fmtString, id.WorkflowId, id.WorkflowVersionVersionNumber, id.TaskId)
}

// Segments returns a slice of Resource ID Segments which comprise this Identity Governance Lifecycle Workflow Deleted Item Workflow Id Version Id Task ID
func (id IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("identityGovernance", "identityGovernance", "identityGovernance"),
		resourceids.StaticSegment("lifecycleWorkflows", "lifecycleWorkflows", "lifecycleWorkflows"),
		resourceids.StaticSegment("deletedItems", "deletedItems", "deletedItems"),
		resourceids.StaticSegment("workflows", "workflows", "workflows"),
		resourceids.UserSpecifiedSegment("workflowId", "workflowId"),
		resourceids.StaticSegment("versions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("workflowVersionVersionNumber", "workflowVersionVersionNumber"),
		resourceids.StaticSegment("tasks", "tasks", "tasks"),
		resourceids.UserSpecifiedSegment("taskId", "taskId"),
	}
}

// String returns a human-readable description of this Identity Governance Lifecycle Workflow Deleted Item Workflow Id Version Id Task ID
func (id IdentityGovernanceLifecycleWorkflowDeletedItemWorkflowIdVersionIdTaskId) String() string {
	components := []string{
		fmt.Sprintf("Workflow: %q", id.WorkflowId),
		fmt.Sprintf("Workflow Version Version Number: %q", id.WorkflowVersionVersionNumber),
		fmt.Sprintf("Task: %q", id.TaskId),
	}
	return fmt.Sprintf("Identity Governance Lifecycle Workflow Deleted Item Workflow Id Version Id Task (%s)", strings.Join(components, "\n"))
}
