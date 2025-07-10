package convert_types

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/report_error/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateReportErrorToModel(req *request.CreateReportError) *model.ReportError {
	return &model.ReportError{
		ReportedBy:  uuid.MustParse(req.ReportedBy),
		ProblemDesc: req.ProblemDesc,
		EpisodeId:   req.EpisodeId,
	}
}

func UpdateReportErrorToModel(req *request.UpdateReportError) *model.ReportError {
	return &model.ReportError{
		HandledBy:    uuid.MustParse(req.HandledBy),
		ProblemDesc:  req.ProblemDesc,
		StatusReport: req.StatusReport,
	}
}
