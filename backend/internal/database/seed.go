package database

import (
	"time"

	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

// SeedDummyData populates the database with a small set of representative
// Organisations, Users, Items, and LoanRequests for local development and
// demos. Organisations, Users, and Items are looked up by their unique field
// (name, email, inventory number) and only created if missing; LoanRequests
// are only seeded once (skipped entirely if any already exist). Calling
// SeedDummyData multiple times, or alongside BootstrapAdmin, is safe.
func SeedDummyData(db *gorm.DB) error {
	informatik, err := seedOrganisation(db, "Fachbereich Informatik")
	if err != nil {
		return err
	}
	elektrotechnik, err := seedOrganisation(db, "Fachbereich Elektrotechnik")
	if err != nil {
		return err
	}

	if _, err := seedUser(db, "logistics@penates.local", "logistics-pw-123", models.RoleLogistics, informatik.ID); err != nil {
		return err
	}
	student1, err := seedUser(db, "student1@penates.local", "student-pw-123", models.RoleStudent, informatik.ID)
	if err != nil {
		return err
	}
	student2, err := seedUser(db, "student2@penates.local", "student-pw-123", models.RoleStudent, elektrotechnik.ID)
	if err != nil {
		return err
	}

	for _, item := range dummyItems() {
		if err := seedItem(db, item); err != nil {
			return err
		}
	}

	return seedLoanRequests(db, student1.ID, student2.ID)
}

func seedOrganisation(db *gorm.DB, name string) (models.Organisation, error) {
	var org models.Organisation
	err := db.Where("name = ?", name).First(&org).Error
	if err == nil {
		return org, nil
	}
	if err != gorm.ErrRecordNotFound {
		return models.Organisation{}, err
	}
	org = models.Organisation{Name: name}
	if err := db.Create(&org).Error; err != nil {
		return models.Organisation{}, err
	}
	return org, nil
}

func seedUser(db *gorm.DB, email, password string, role models.Role, organisationID uint) (models.User, error) {
	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	if err == nil {
		return user, nil
	}
	if err != gorm.ErrRecordNotFound {
		return models.User{}, err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return models.User{}, err
	}
	user = models.User{
		Email:          email,
		PasswordHash:   hash,
		Role:           role,
		OrganisationID: organisationID,
	}
	if err := db.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedItem(db *gorm.DB, item models.Item) error {
	var existing models.Item
	err := db.Where("inventory_number = ?", item.InventoryNumber).First(&existing).Error
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return db.Create(&item).Error
}

func seedLoanRequests(db *gorm.DB, student1ID, student2ID uint) error {
	var existingCount int64
	if err := db.Model(&models.LoanRequest{}).Count(&existingCount).Error; err != nil {
		return err
	}
	if existingCount > 0 {
		return nil
	}

	repo := repository.NewGormLoanRequestRepository(db)
	for _, seed := range dummyLoanRequests(student1ID, student2ID) {
		loanRequest := seed.loanRequest
		if err := repo.Create(&loanRequest, seed.itemInventoryNumbers); err != nil {
			return err
		}
	}
	return nil
}

func dummyItems() []models.Item {
	return []models.Item{
		{
			InventoryNumber:              "100001",
			Name:                         "Laptop Dell Latitude 5420",
			Description:                  `14" business laptop, i5, 16GB RAM`,
			Category:                     "Electronics",
			Location:                     "Room 101",
			Amount:                       5,
			Manufacturer:                 "Dell",
			SerialNumber:                 "DL5420-0001",
			Note:                         "Comes with charger",
			ElectricalAppliance:          true,
			LastTechnicalInspectionDate:  datePtr(daysAgo(90)),
			LastElectricalInspectionDate: datePtr(daysAgo(90)),
			DateOfPurchase:               datePtr(daysAgo(400)),
			Price:                        899.99,
			ResolutionNumber:             "RES-2024-001",
		},
		{
			InventoryNumber:              "100002",
			Name:                         "Beamer Epson EB-2250U",
			Description:                  "Full HD business projector",
			Category:                     "Electronics",
			Location:                     "Room 204",
			Amount:                       2,
			Manufacturer:                 "Epson",
			SerialNumber:                 "EP2250-0002",
			ElectricalAppliance:          true,
			LastTechnicalInspectionDate:  datePtr(daysAgo(60)),
			LastElectricalInspectionDate: datePtr(daysAgo(60)),
			DateOfPurchase:               datePtr(daysAgo(250)),
			Price:                        749.00,
			ResolutionNumber:             "RES-2024-002",
		},
		{
			InventoryNumber:     "100003",
			Name:                "Werkzeugkoffer",
			Description:         "Standard-Werkzeugsatz mit 120 Teilen",
			Category:            "Tools",
			Location:            "Werkstatt",
			Amount:              3,
			Manufacturer:        "Bosch",
			SerialNumber:        "WK-0003",
			ElectricalAppliance: false,
			DateOfPurchase:      datePtr(daysAgo(600)),
			Price:               129.50,
			ResolutionNumber:    "RES-2023-014",
		},
		{
			InventoryNumber:              "100004",
			Name:                         "Multimeter Fluke 87V",
			Description:                  "Digital-Multimeter, True-RMS",
			Category:                     "Electronics",
			Location:                     "Room 101",
			Amount:                       4,
			Manufacturer:                 "Fluke",
			SerialNumber:                 "FL87V-0004",
			ElectricalAppliance:          true,
			LastTechnicalInspectionDate:  datePtr(daysAgo(30)),
			LastElectricalInspectionDate: datePtr(daysAgo(30)),
			DateOfPurchase:               datePtr(daysAgo(180)),
			Price:                        399.00,
			ResolutionNumber:             "RES-2024-009",
		},
		{
			InventoryNumber:     "100005",
			Name:                "Klapptisch",
			Description:         "Faltbarer Veranstaltungstisch, 180x80cm",
			Category:            "Furniture",
			Location:            "Lager",
			Amount:              10,
			Manufacturer:        "Lifetime",
			SerialNumber:        "KT-0005",
			ElectricalAppliance: false,
			DateOfPurchase:      datePtr(daysAgo(800)),
			Price:               89.90,
			ResolutionNumber:    "RES-2022-031",
		},
		{
			InventoryNumber:              "100006",
			Name:                         "Lötkolben-Set Weller",
			Description:                  "Temperaturgeregelte Lötstation mit Zubehör",
			Category:                     "Tools",
			Location:                     "Werkstatt",
			Amount:                       6,
			Manufacturer:                 "Weller",
			SerialNumber:                 "WL-0006",
			ElectricalAppliance:          true,
			LastTechnicalInspectionDate:  datePtr(daysAgo(45)),
			LastElectricalInspectionDate: datePtr(daysAgo(45)),
			DateOfPurchase:               datePtr(daysAgo(120)),
			Price:                        159.00,
			ResolutionNumber:             "RES-2024-011",
		},
	}
}

type dummyLoanRequest struct {
	loanRequest          models.LoanRequest
	itemInventoryNumbers []string
}

func dummyLoanRequests(student1ID, student2ID uint) []dummyLoanRequest {
	return []dummyLoanRequest{
		{
			loanRequest: models.LoanRequest{
				RequestingUserID: student1ID,
				DateOfLending:    datePtr(daysFromNow(2)),
				DateOfReturn:     datePtr(daysFromNow(9)),
				LocationOfItems:  "Lab 2",
				Status:           "pending",
			},
			itemInventoryNumbers: []string{"100001", "100004"},
		},
		{
			loanRequest: models.LoanRequest{
				RequestingUserID: student2ID,
				DateOfLending:    datePtr(daysAgo(3)),
				DateOfReturn:     datePtr(daysFromNow(4)),
				LocationOfItems:  "Hoersaal A",
				Status:           "approved",
			},
			itemInventoryNumbers: []string{"100002"},
		},
		{
			loanRequest: models.LoanRequest{
				RequestingUserID: student1ID,
				DateOfLending:    datePtr(daysAgo(20)),
				DateOfReturn:     datePtr(daysAgo(13)),
				LocationOfItems:  "Werkstatt",
				Status:           "returned",
			},
			itemInventoryNumbers: []string{"100003", "100006"},
		},
	}
}

func datePtr(t time.Time) *time.Time {
	return &t
}

func daysAgo(days int) time.Time {
	return time.Now().AddDate(0, 0, -days)
}

func daysFromNow(days int) time.Time {
	return time.Now().AddDate(0, 0, days)
}
