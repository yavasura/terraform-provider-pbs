/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package notifications

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
	pbsnotifications "github.com/yavasura/terraform-provider-pbs/pbs/notifications"
)

func shouldDeleteListAttr(plan, state types.List) bool {
	return plan.IsNull() && !state.IsNull() && !state.IsUnknown()
}

func shouldDeleteStringAttr(plan, state types.String) bool {
	return plan.IsNull() && !state.IsNull() && !state.IsUnknown()
}

func setNotificationCommonState(comment string, disable *bool, origin string, commentState *types.String, disableState *types.Bool, originState *types.String) {
	*commentState = tfvalue.StringOrNull(comment)
	*disableState = tfvalue.BoolPtrOrNull(disable)
	*originState = tfvalue.StringOrNull(origin)
}

func stringListState(ctx context.Context, value []string) (types.List, diag.Diagnostics) {
	if len(value) == 0 {
		return types.ListNull(types.StringType), nil
	}
	return types.ListValueFrom(ctx, types.StringType, value)
}

func stringMapState(ctx context.Context, value map[string]string) (types.Map, diag.Diagnostics) {
	if len(value) == 0 {
		return types.MapNull(types.StringType), nil
	}
	return types.MapValueFrom(ctx, types.StringType, value)
}

func decodeStringList(ctx context.Context, value types.List, diags *diag.Diagnostics) []string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	var items []string
	diags.Append(value.ElementsAs(ctx, &items, false)...)
	if diags.HasError() {
		return nil
	}
	return items
}

func buildNotificationMatcher(ctx context.Context, plan, state *notificationMatcherResourceModel, diags *diag.Diagnostics) *pbsnotifications.NotificationMatcher {
	matcher := &pbsnotifications.NotificationMatcher{
		Name: plan.Name.ValueString(),
	}

	if state != nil {
		matcher.Delete = computeMatcherDeletes(plan, state)
	}

	matcher.Targets = decodeStringList(ctx, plan.Targets, diags)
	matcher.MatchSeverity = decodeStringList(ctx, plan.MatchSeverity, diags)
	matcher.MatchField = decodeStringList(ctx, plan.MatchField, diags)
	matcher.MatchCalendar = decodeStringList(ctx, plan.MatchCalendar, diags)
	if diags.HasError() {
		return nil
	}

	if !plan.Mode.IsNull() {
		matcher.Mode = plan.Mode.ValueString()
	}

	if !plan.InvertMatch.IsNull() {
		invertMatch := plan.InvertMatch.ValueBool()
		matcher.InvertMatch = &invertMatch
	}

	if !plan.Comment.IsNull() {
		matcher.Comment = plan.Comment.ValueString()
	}

	if !plan.Disable.IsNull() {
		disable := plan.Disable.ValueBool()
		matcher.Disable = &disable
	}

	return matcher
}

func setNotificationMatcherState(ctx context.Context, matcher *pbsnotifications.NotificationMatcher, state *notificationMatcherResourceModel, diags *diag.Diagnostics) {
	state.Name = types.StringValue(matcher.Name)

	state.Targets, *diags = listStateWithDiags(ctx, matcher.Targets, *diags)
	state.MatchSeverity, *diags = listStateWithDiags(ctx, matcher.MatchSeverity, *diags)
	state.MatchField, *diags = listStateWithDiags(ctx, matcher.MatchField, *diags)
	state.MatchCalendar, *diags = listStateWithDiags(ctx, matcher.MatchCalendar, *diags)

	state.Mode = types.StringValue("all")
	if matcher.Mode != "" {
		state.Mode = types.StringValue(matcher.Mode)
	}

	state.InvertMatch = types.BoolValue(false)
	if matcher.InvertMatch != nil {
		state.InvertMatch = types.BoolValue(*matcher.InvertMatch)
	}

	setNotificationCommonState(matcher.Comment, matcher.Disable, matcher.Origin, &state.Comment, &state.Disable, &state.Origin)
}

func listStateWithDiags(ctx context.Context, value []string, diags diag.Diagnostics) (types.List, diag.Diagnostics) {
	list, listDiags := stringListState(ctx, value)
	diags.Append(listDiags...)
	return list, diags
}
