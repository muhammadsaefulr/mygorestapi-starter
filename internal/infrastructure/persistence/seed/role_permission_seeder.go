package seed

import (
	"log"
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var allRoles = map[string][]string{
	"user": {"getUserSession"},
	"admin": {
		"getUsers", "manageUsers", "getUserSession", "manageAnime",
		"getUserRole", "getRolePermissions",
		"createMovieDetails", "updateMovieDetails", "deleteMovieDetails",
		"addMovieEps", "updateMovieEps", "deleteMovieEps",
		"postSubsPlan", "updateSubsPlan", "deleteSubsPlan",
		"userSubscriptionGet", "userSubscriptionPost", "userSubscriptionPut", "userSubscriptionDelete",
		"getAllReportError", "postReportError", "getReportErrorByID",
		"updateReportError", "deleteReportError",
	},
	"owner": {
		"allActions",
	},
}

func SeedRolesAndPermissions(db *gorm.DB) error {
	// 1. Seed all unique permissions
	permMap := make(map[string]model.RolePermissions)
	for _, perms := range allRoles {
		for _, name := range perms {
			permMap[name] = model.RolePermissions{
				PermissionName: name,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
		}
	}
	var permissions []model.RolePermissions
	for _, p := range permMap {
		permissions = append(permissions, p)
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&permissions).Error; err != nil {
		return err
	}

	// 2. Reload permissions from DB into map[name] => record
	var permissionRecords []model.RolePermissions
	if err := db.Find(&permissionRecords).Error; err != nil {
		return err
	}
	permLookup := make(map[string]model.RolePermissions)
	for _, p := range permissionRecords {
		permLookup[p.PermissionName] = p
	}

	// 3. Seed roles one by one, then assign permissions
	for roleName, permNames := range allRoles {
		// Insert role (ON CONFLICT DO NOTHING)
		role := model.UserRole{
			RoleName:  roleName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error; err != nil {
			log.Printf("❌ Gagal insert role %s: %v", roleName, err)
			continue
		}

		// Fetch role (guaranteed exists now)
		var dbRole model.UserRole
		if err := db.Where("role_name = ?", roleName).First(&dbRole).Error; err != nil {
			log.Printf("❌ Gagal fetch role %s: %v", roleName, err)
			continue
		}

		// Collect permission records
		var rolePerms []model.RolePermissions
		for _, pname := range permNames {
			if perm, ok := permLookup[pname]; ok {
				rolePerms = append(rolePerms, perm)
			} else {
				log.Printf("⚠️ Permission %s tidak ditemukan di DB", pname)
			}
		}

		// Assign permissions via many2many
		if err := db.Model(&dbRole).Association("Permissions").Replace(rolePerms); err != nil {
			log.Printf("❌ Gagal assign permission untuk role %s: %v", roleName, err)
		}
	}

	log.Println("✅ Seeder role & permission berhasil.")
	return nil
}
