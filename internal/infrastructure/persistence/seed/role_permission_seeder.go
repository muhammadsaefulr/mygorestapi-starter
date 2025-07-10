package seed

import (
	"time"

	"log"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var allRoles = map[string][]string{
	"user": {"getUserSession"},
	"admin": {
		"getUsers", "manageUsers", "getUserSession", "manageAnime",
		"createMovieDetails", "updateMovieDetails", "deleteMovieDetails",
		"addMovieEps", "updateMovieEps", "deleteMovieEps",
		"getAllReportError", "postReportError", "getReportErrorByID",
		"updateReportError", "deleteReportError",
	},
}

func SeedRolesAndPermissions(db *gorm.DB) error {
	// Step 1: Insert semua permission unik
	permMap := make(map[string]model.RolePermissions)

	for _, perms := range allRoles {
		for _, perm := range perms {
			permMap[perm] = model.RolePermissions{
				PermissionName: perm,
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

	var allPermissionRecords []model.RolePermissions
	if err := db.Find(&allPermissionRecords).Error; err != nil {
		return err
	}
	permLookup := make(map[string]model.RolePermissions)
	for _, p := range allPermissionRecords {
		permLookup[p.PermissionName] = p
	}

	for roleName, permNames := range allRoles {
		var relatedPerms []model.RolePermissions
		for _, pname := range permNames {
			if p, exists := permLookup[pname]; exists {
				relatedPerms = append(relatedPerms, p)
			}
		}

		role := model.UserRole{
			RoleName:   roleName,
			Permission: relatedPerms,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error
		if err != nil {
			log.Printf("Gagal insert role %s: %v", roleName, err)
			continue
		}

		var existingRole model.UserRole
		if err := db.Preload("Permissions").Where("role_name = ?", roleName).First(&existingRole).Error; err == nil {
			if err := db.Model(&existingRole).Association("Permissions").Replace(relatedPerms); err != nil {
				log.Printf("Gagal assign permission untuk role %s: %v", roleName, err)
			}
		}
	}

	log.Println("✅ Seeder role & permission berhasil.")
	return nil
}
