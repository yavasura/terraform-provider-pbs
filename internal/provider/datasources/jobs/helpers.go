package jobs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/yavasura/terraform-provider-pbs/internal/provider/tfvalue"
	pbsjobs "github.com/yavasura/terraform-provider-pbs/pbs/jobs"
)

func stringValueOrNull(value string) types.String {
	return tfvalue.StringOrNull(value)
}

func int64ValueOrNull(value *int) types.Int64 {
	return tfvalue.IntPtrOrNull(value)
}

func boolValueOrNull(value *bool) types.Bool {
	return tfvalue.BoolPtrOrNull(value)
}

func stringListValueOrNull(ctx context.Context, value []string) (types.List, diag.Diagnostics) {
	return tfvalue.StringListOrNull(ctx, value)
}

func setPruneJobDataSourceState(job *pbsjobs.PruneJob, state *pruneJobDataSourceModel) {
	state.ID = types.StringValue(job.ID)
	state.Store = types.StringValue(job.Store)
	state.Schedule = types.StringValue(job.Schedule)
	state.KeepLast = int64ValueOrNull(job.KeepLast)
	state.KeepHourly = int64ValueOrNull(job.KeepHourly)
	state.KeepDaily = int64ValueOrNull(job.KeepDaily)
	state.KeepWeekly = int64ValueOrNull(job.KeepWeekly)
	state.KeepMonthly = int64ValueOrNull(job.KeepMonthly)
	state.KeepYearly = int64ValueOrNull(job.KeepYearly)
	state.MaxDepth = int64ValueOrNull(job.MaxDepth)
	state.Namespace = stringValueOrNull(job.Namespace)
	state.Comment = stringValueOrNull(job.Comment)
	state.Disable = boolValueOrNull(job.Disable)
	state.Digest = stringValueOrNull(job.Digest)
}

func pruneJobToState(job *pbsjobs.PruneJob, state *pruneJobDataSourceModel) {
	setPruneJobDataSourceState(job, state)
}

func setVerifyJobDataSourceState(job *pbsjobs.VerifyJob, state *verifyJobDataSourceModel) {
	state.ID = types.StringValue(job.ID)
	state.Store = types.StringValue(job.Store)
	state.Schedule = types.StringValue(job.Schedule)
	state.IgnoreVerified = boolValueOrNull(job.IgnoreVerified)
	state.OutdatedAfter = int64ValueOrNull(job.OutdatedAfter)
	state.Namespace = stringValueOrNull(job.Namespace)
	state.MaxDepth = int64ValueOrNull(job.MaxDepth)
	state.Comment = stringValueOrNull(job.Comment)
	state.Digest = stringValueOrNull(job.Digest)
}

func verifyJobToState(job *pbsjobs.VerifyJob, state *verifyJobDataSourceModel) {
	setVerifyJobDataSourceState(job, state)
}

func setSyncJobDataSourceState(ctx context.Context, job *pbsjobs.SyncJob, state *syncJobDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.ID = types.StringValue(job.ID)
	state.Store = types.StringValue(job.Store)
	state.Schedule = types.StringValue(job.Schedule)
	state.Remote = types.StringValue(job.Remote)
	state.RemoteStore = types.StringValue(job.RemoteStore)
	state.RemoteNamespace = stringValueOrNull(job.RemoteNamespace)
	state.Namespace = stringValueOrNull(job.Namespace)
	state.MaxDepth = int64ValueOrNull(job.MaxDepth)
	state.RemoveVanished = boolValueOrNull(job.RemoveVanished)
	state.Comment = stringValueOrNull(job.Comment)
	state.Digest = stringValueOrNull(job.Digest)

	groupFilter, listDiags := stringListValueOrNull(ctx, job.GroupFilter)
	diags.Append(listDiags...)
	state.GroupFilter = groupFilter

	return diags
}
