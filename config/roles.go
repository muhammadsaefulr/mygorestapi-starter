package config

var allRoles = map[string][]string{
	"user": {"getUserSession"},
	"admin": {"getUsers", "manageUsers", "getUserSession", "manageAnime",
		"createMovieDetails", "updateMovieDetails", "deleteMovieDetails",
		"addMovieEps", "updateMovieEps", "deleteMovieEps",
		"getAllReportError", "postReportError", "getReportErrorByID", "updateReportError", "deleteReportError"},
	"vip": {},
}

var Roles = getKeys(allRoles)
var RoleRights = allRoles

func getKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
