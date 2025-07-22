package seed

import (
	"log"
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

var allRoles = map[string][]string{
	"user": {"getUserSession", "postReportError", "getBannerApp", "getUserBadge"},
	"admin": {
		"getUsers", "manageUsers", "getUserSession", "manageAnime",
		"getUserRole", "getRolePermissions",
		"createMovieDetails", "updateMovieDetails", "deleteMovieDetails",
		"addMovieEps", "updateMovieEps", "deleteMovieEps",
		"postSubsPlan", "updateSubsPlan", "deleteSubsPlan",
		"userSubscriptionGet", "userSubscriptionPost", "userSubscriptionPut", "userSubscriptionDelete",
		"getAllReportError", "postReportError", "getReportErrorByID",
		"updateReportError", "deleteReportError",
		"createBannerApp", "updateBannerApp", "deleteBannerApp", "getBannerApp",
		"addUserBadge", "updateUserBadge", "deleteUserBadge", "getUserBadge",
		"addBadge", "updateBadge", "deleteBadge", "getBadge",
		"getRequestVIP", "updateRequestVIP", "deleteRequestVIP",
	},
	"owner": {
		"allActions",
	},
}

func SeedRolesAndPermissions(db *gorm.DB) error {
	// 1. Insert all unique permissions
	permSet := make(map[string]struct{})
	for _, perms := range allRoles {
		for _, name := range perms {
			permSet[name] = struct{}{}
		}
	}

	for name := range permSet {
		var count int64
		db.Model(&model.RolePermissions{}).Where("permission_name = ?", name).Count(&count)
		if count == 0 {
			perm := model.RolePermissions{
				PermissionName: name,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			if err := db.Create(&perm).Error; err != nil {
				log.Printf("❌ Gagal insert permission %s: %v", name, err)
			}
		}
	}

	// 2. Load all permissions to a map
	var allPerms []model.RolePermissions
	if err := db.Find(&allPerms).Error; err != nil {
		return err
	}
	permMap := make(map[string]model.RolePermissions)
	for _, p := range allPerms {
		permMap[p.PermissionName] = p
	}

	// 3. Process each role
	for roleName, permNames := range allRoles {
		// Ensure role exists
		var role model.UserRole
		if err := db.Where("role_name = ?", roleName).First(&role).Error; err != nil {
			role = model.UserRole{
				RoleName:  roleName,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := db.Create(&role).Error; err != nil {
				log.Printf("❌ Gagal insert role %s: %v", roleName, err)
				continue
			}
		}

		// Load current permissions for the role
		var currentPerms []model.RolePermissions
		if err := db.Model(&role).Association("Permissions").Find(&currentPerms); err != nil {
			log.Printf("❌ Gagal load permissions untuk role %s: %v", roleName, err)
			continue
		}
		existingPerms := make(map[string]bool)
		for _, p := range currentPerms {
			existingPerms[p.PermissionName] = true
		}

		// Append only new permissions
		for _, permName := range permNames {
			if !existingPerms[permName] {
				if perm, ok := permMap[permName]; ok {
					if err := db.Model(&role).Association("Permissions").Append(&perm); err != nil {
						log.Printf("❌ Gagal assign permission %s ke role %s: %v", permName, roleName, err)
					}
				} else {
					log.Printf("⚠️ Permission %s tidak ditemukan di permMap", permName)
				}
			}
		}
	}

	log.Println("✅ Seeder role & permission selesai.")
	return nil
}
