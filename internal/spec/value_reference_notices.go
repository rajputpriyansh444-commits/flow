package spec

import cmnvalue "github.com/rajputpriyansh444-commits/flow/internal/cmn/value"

func reportNamespaceUnavailableStepOutputReferences(sink cmnvalue.ValueReferenceNoticeSink, field, raw string) {
	for _, ref := range cmnvalue.StepOutputReferences(raw) {
		cmnvalue.ReportStepOutputReferenceNotice(
			sink,
			field,
			ref.Expression,
			cmnvalue.ValueReferenceReasonNamespaceUnavailable,
		)
	}
}

func buildNoticeSink(sink cmnvalue.ValueReferenceNoticeSink) cmnvalue.ValueReferenceNoticeSink {
	return cmnvalue.SuppressStepOutputReferenceNotices(sink)
}
