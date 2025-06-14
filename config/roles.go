package config

var allRoles = map[string][]string{
	"user":  {},
	"admin": {"getUsers", "manageUsers", "manageAnime"},
	"vip":   {},
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
