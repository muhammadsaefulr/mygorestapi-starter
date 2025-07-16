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
		TypeMovie:   req.TypeMovie,
		EpisodeId:   req.EpisodeId,
	}
}

func UpdateReportErrorToModel(req *request.UpdateReportError) *model.ReportError {

	var handledBy *uuid.UUID
	if req.HandledBy != "" {
		parsed := uuid.MustParse(req.HandledBy)
		handledBy = &parsed
	}

	return &model.ReportError{
		HandledBy:    handledBy,
		ProblemDesc:  req.ProblemDesc,
		TypeMovie:    req.TypeMovie,
		StatusReport: req.StatusReport,
	}
}
