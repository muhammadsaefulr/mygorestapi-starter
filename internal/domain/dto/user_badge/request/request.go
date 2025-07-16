package request

type CreateUserBadge struct {
	BadgeName string `json:"badge_name"`
	IconURL   string `json:"icon_url"`
	Color     string `json:"color"`
}

type UpdateUserBadge struct {
	ID        uint   `json:"-"`
	BadgeName string `json:"badge_name"`
	IconURL   string `json:"icon_url"`
	Color     string `json:"color"`
}

type QueryUserBadge struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

type CreateUserBadgeInfo struct {
	UserID    string `json:"user_id"`
	BadgeID   uint   `json:"badge_id"`
	Note      string `json:"note"`
	HandledBy string `json:"-"`
}

type UpdateUserBadgeInfo struct {
	UserID    string `json:"-"`
	BadgeID   uint   `json:"badge_id"`
	Note      string `json:"note"`
	HandledBy string `json:"-"`
}
