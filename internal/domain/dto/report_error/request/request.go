package request

type CreateReportError struct {
	ReportedBy  string `json:"-"`
	TypeMovie   string `json:"movie_type"`
	ProblemDesc string `json:"problem_desc"`
	EpisodeId   string `json:"episode_id"`
}

type UpdateReportError struct {
	HandledBy    string `json:"-"`
	TypeMovie    string `json:"movie_type"`
	ProblemDesc  string `json:"problem_desc"`
	StatusReport string `json:"status_report"`
}

type QueryReportError struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Type   string `query:"type"`
	Search string `query:"search"`
}
