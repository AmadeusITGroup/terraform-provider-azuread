// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identitygovernance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/common-types/stable"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/identitygovernance/stable/privilegedaccessgroupassignmentschedule"
	"github.com/hashicorp/go-azure-sdk/microsoft-graph/identitygovernance/stable/privilegedaccessgroupassignmentschedulerequest"
	"github.com/hashicorp/go-azure-sdk/sdk/nullable"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-provider-azuread/internal/helpers/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azuread/internal/sdk"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/identitygovernance/parse"
)

var _ sdk.ResourceWithUpdate = PrivilegedAccessGroupAssignmentScheduleResource{}

type PrivilegedAccessGroupAssignmentScheduleResource struct{}

func (r PrivilegedAccessGroupAssignmentScheduleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.ValidatePrivilegedAccessGroupScheduleID
}

func (r PrivilegedAccessGroupAssignmentScheduleResource) ResourceType() string {
	return "azuread_privileged_access_group_assignment_schedule"
}

func (r PrivilegedAccessGroupAssignmentScheduleResource) ModelObject() interface{} {
	return &PrivilegedAccessGroupScheduleModel{}
}

func (r PrivilegedAccessGroupAssignmentScheduleResource) Arguments() map[string]*pluginsdk.Schema {
	return privilegedAccessGroupScheduleArguments()
}

func (r PrivilegedAccessGroupAssignmentScheduleResource) Attributes() map[string]*pluginsdk.Schema {
	return privilegedAccessGroupScheduleAttributes()
}

func (r PrivilegedAccessGroupAssignmentScheduleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IdentityGovernance.PrivilegedAccessGroupAssignmentScheduleRequestClient

			var model PrivilegedAccessGroupScheduleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			schedule, err := buildScheduleRequest(&model, &metadata)
			if err != nil {
				return err
			}

			properties := stable.PrivilegedAccessGroupAssignmentScheduleRequest{
				AccessId:      stable.PrivilegedAccessGroupRelationships(model.AssignmentType),
				PrincipalId:   nullable.Value(model.PrincipalId),
				GroupId:       nullable.Value(model.GroupId),
				Action:        pointer.To(stable.ScheduleRequestActions_AdminAssign),
				Justification: nullable.NoZero(model.Justification),
				ScheduleInfo:  schedule,
			}

			if model.TicketNumber != "" || model.TicketSystem != "" {
				properties.TicketInfo = &stable.TicketInfo{
					TicketNumber: nullable.NoZero(model.TicketNumber),
					TicketSystem: nullable.NoZero(model.TicketSystem),
				}
			}

			resp, err := client.CreatePrivilegedAccessGroupAssignmentScheduleRequest(ctx, properties, privilegedaccessgroupassignmentschedulerequest.DefaultCreatePrivilegedAccessGroupAssignmentScheduleRequestOperationOptions())
			if err != nil {
				return fmt.Errorf("creating assignment schedule request: %v", err)
			}

			request := resp.Model
			if request == nil {
				return fmt.Errorf("creating assignment schedule request: model was nil")
			}
			if request.Id == nil || *request.Id == "" {
				return fmt.Errorf("creating assignment schedule request: ID returned for request is nil/empty")
			}

			if pointer.From(request.Status) == PrivilegedAccessGroupScheduleRequestStatusFailed {
				return fmt.Errorf("creating assignment schedule request: request is in a failed state")
			}

			resourceId, err := parse.ParsePrivilegedAccessGroupScheduleID(request.TargetScheduleId.GetOrZero())
			if err != nil {
				return err
			}

			metadata.SetID(resourceId)

			return nil
		},
	}
}

func (r PrivilegedAccessGroupAssignmentScheduleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			scheduleClient := metadata.Client.IdentityGovernance.PrivilegedAccessGroupAssignmentScheduleClient
			requestsClient := metadata.Client.IdentityGovernance.PrivilegedAccessGroupAssignmentScheduleRequestClient

			resourceId, err := parse.ParsePrivilegedAccessGroupScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivilegedAccessGroupScheduleModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := stable.NewIdentityGovernancePrivilegedAccessGroupAssignmentScheduleID(resourceId.ID())

			scheduleResp, err := scheduleClient.GetPrivilegedAccessGroupAssignmentSchedule(ctx, id, privilegedaccessgroupassignmentschedule.DefaultGetPrivilegedAccessGroupAssignmentScheduleOperationOptions())
			if err != nil && !response.WasNotFound(scheduleResp.HttpResponse) {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			schedule := scheduleResp.Model

			// Some details are only available on the request which is used for the create/update of the schedule.
			// Schedule requests are never deleted. New ones are created when changes are made.
			// Therefore, on read, we need to find the latest version of the request.
			// This is to cater for changes being made outside of Terraform.
			options := privilegedaccessgroupassignmentschedulerequest.ListPrivilegedAccessGroupAssignmentScheduleRequestsOperationOptions{
				Filter: pointer.To(fmt.Sprintf("groupId eq '%s' and targetScheduleId eq '%s'", resourceId.GroupId, resourceId.ID())),
				OrderBy: pointer.To(odata.OrderBy{
					Field:     "createdDateTime",
					Direction: odata.Descending,
				}),
			}
			requestsResp, err := requestsClient.ListPrivilegedAccessGroupAssignmentScheduleRequests(ctx, options)
			if err != nil {
				return fmt.Errorf("listing requests: %v", err)
			}

			var request *stable.PrivilegedAccessGroupAssignmentScheduleRequest

			requests := requestsResp.Model
			if requests == nil || len(*requests) == 0 {
				if response.WasNotFound(scheduleResp.HttpResponse) {
					// No request and no schedule was found
					return metadata.MarkAsGone(resourceId)
				}
			} else {
				request = pointer.To((*requests)[0])
			}

			var scheduleInfo *stable.RequestSchedule

			if request != nil {
				// The request is still present, populate from the request
				scheduleInfo = request.ScheduleInfo

				model.AssignmentType = string(request.AccessId)
				model.GroupId = request.GroupId.GetOrZero()
				model.Justification = request.Justification.GetOrZero()
				model.PrincipalId = request.PrincipalId.GetOrZero()
				model.Status = pointer.From(request.Status)

				if ticketInfo := request.TicketInfo; ticketInfo != nil {
					model.TicketNumber = ticketInfo.TicketNumber.GetOrZero()
					model.TicketSystem = ticketInfo.TicketSystem.GetOrZero()
				}
			} else if schedule != nil {
				// The request has likely expired, so populate from the schedule
				scheduleInfo = &schedule.ScheduleInfo

				model.AssignmentType = string(schedule.AccessId)
				model.GroupId = schedule.GroupId.GetOrZero()
				model.PrincipalId = schedule.PrincipalId.GetOrZero()
				model.Status = schedule.Status.GetOrZero()
			}

			if scheduleInfo != nil {
				model.StartDate = scheduleInfo.StartDateTime.GetOrZero()

				if expiration := scheduleInfo.Expiration; expiration != nil {
					model.Duration = expiration.Duration.GetOrZero()
					model.ExpirationDate = expiration.EndDateTime.GetOrZero()

					if expiration.Type != nil {
						model.PermanentAssignment = *expiration.Type == stable.ExpirationPatternType_NoExpiration
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r PrivilegedAccessGroupAssignmentScheduleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IdentityGovernance.PrivilegedAccessGroupAssignmentScheduleRequestClient

			resourceId, err := parse.ParsePrivilegedAccessGroupScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivilegedAccessGroupScheduleModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			schedule, err := buildScheduleRequest(&model, &metadata)
			if err != nil {
				return err
			}

			properties := stable.PrivilegedAccessGroupAssignmentScheduleRequest{
				AccessId:      stable.PrivilegedAccessGroupRelationships(resourceId.Relationship),
				PrincipalId:   nullable.Value(model.PrincipalId),
				GroupId:       nullable.Value(resourceId.GroupId),
				Action:        pointer.To(stable.ScheduleRequestActions_AdminAssign),
				Justification: nullable.NoZero(model.Justification),
				ScheduleInfo:  schedule,
			}

			if model.TicketNumber != "" || model.TicketSystem != "" {
				properties.TicketInfo = &stable.TicketInfo{
					TicketNumber: nullable.NoZero(model.TicketNumber),
					TicketSystem: nullable.NoZero(model.TicketSystem),
				}
			}

			resp, err := client.CreatePrivilegedAccessGroupAssignmentScheduleRequest(ctx, properties, privilegedaccessgroupassignmentschedulerequest.DefaultCreatePrivilegedAccessGroupAssignmentScheduleRequestOperationOptions())
			if err != nil {
				return fmt.Errorf("creating updated assignment schedule request: %v", err)
			}

			request := resp.Model
			if request == nil {
				return fmt.Errorf("creating updated assignment schedule request: model was nil")
			}
			if request.Id == nil || *request.Id == "" {
				return fmt.Errorf("creating updated assignment schedule request: ID returned for request is nil/empty")
			}

			if pointer.From(request.Status) == PrivilegedAccessGroupScheduleRequestStatusFailed {
				return fmt.Errorf("creating updated assignment schedule request: request is in a failed state")
			}

			return nil
		},
	}
}

func (r PrivilegedAccessGroupAssignmentScheduleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IdentityGovernance.PrivilegedAccessGroupAssignmentScheduleRequestClient

			resourceId, err := parse.ParsePrivilegedAccessGroupScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivilegedAccessGroupScheduleModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			switch model.Status {
			case PrivilegedAccessGroupScheduleRequestStatusDenied,
				PrivilegedAccessGroupScheduleRequestStatusFailed,
				PrivilegedAccessGroupScheduleRequestStatusGranted,
				PrivilegedAccessGroupScheduleRequestStatusPendingAdminDecision,
				PrivilegedAccessGroupScheduleRequestStatusPendingApproval,
				PrivilegedAccessGroupScheduleRequestStatusPendingProvisioning,
				PrivilegedAccessGroupScheduleRequestStatusPendingScheduleCreation:
				return cancelAssignmentRequest(ctx, metadata, client, stable.NewIdentityGovernancePrivilegedAccessGroupAssignmentScheduleRequestID(resourceId.ID()))
			case PrivilegedAccessGroupScheduleRequestStatusProvisioned,
				PrivilegedAccessGroupScheduleRequestStatusScheduleCreated:
				return revokeAssignmentRequest(ctx, metadata, client, *resourceId, model)
			case PrivilegedAccessGroupScheduleRequestStatusCanceled,
				PrivilegedAccessGroupScheduleRequestStatusRevoked:
				return metadata.MarkAsGone(resourceId)
			}

			return fmt.Errorf("unable to destroy due to unknown status: %s", model.Status)
		},
	}
}

func cancelAssignmentRequest(ctx context.Context, metadata sdk.ResourceMetaData, client *privilegedaccessgroupassignmentschedulerequest.PrivilegedAccessGroupAssignmentScheduleRequestClient, id stable.IdentityGovernancePrivilegedAccessGroupAssignmentScheduleRequestId) error {
	if resp, err := client.CancelPrivilegedAccessGroupAssignmentScheduleRequest(ctx, id, privilegedaccessgroupassignmentschedulerequest.DefaultCancelPrivilegedAccessGroupAssignmentScheduleRequestOperationOptions()); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return metadata.MarkAsGone(id)
		}

		return fmt.Errorf("canceling %s: %v", id, err)
	}

	return nil
}

func revokeAssignmentRequest(ctx context.Context, metadata sdk.ResourceMetaData, client *privilegedaccessgroupassignmentschedulerequest.PrivilegedAccessGroupAssignmentScheduleRequestClient, id parse.PrivilegedAccessGroupScheduleId, model PrivilegedAccessGroupScheduleModel) error {
	request := stable.PrivilegedAccessGroupAssignmentScheduleRequest{
		AccessId:    stable.PrivilegedAccessGroupRelationships(id.Relationship),
		PrincipalId: nullable.Value(model.PrincipalId),
		GroupId:     nullable.Value(id.GroupId),
		Action:      pointer.To(stable.ScheduleRequestActions_AdminRemove),
	}

	if resp, err := client.CreatePrivilegedAccessGroupAssignmentScheduleRequest(ctx, request, privilegedaccessgroupassignmentschedulerequest.DefaultCreatePrivilegedAccessGroupAssignmentScheduleRequestOperationOptions()); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return metadata.MarkAsGone(&id)
		}

		return fmt.Errorf("creating schedule removal request: %v", err)
	}

	return nil
}
