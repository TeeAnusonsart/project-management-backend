package gorm

import (
	"log"
	"project-home-iot/internal/auth"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) error {
	SeedRooms(db)
	SeedCapabilities(db)
	SeedUsers(db)
	SeedDevices(db)
	log.Println("Seeding completed!")
	return nil
}

func SeedRooms(db *gorm.DB) {
	rooms := []Room{
		{Name: "Living Room"},
		{Name: "Bedroom"},
		{Name: "Kitchen"},
	}
	for _, r := range rooms {
		db.FirstOrCreate(&r, Room{Name: r.Name})
	}
}
func SeedCapabilities(db *gorm.DB) {
	capabilities := []Capability{
		{CapabilityType: "toggle", ControlType: "light-switch"},
		{CapabilityType: "mode", ControlType: "light-system"},
		{CapabilityType: "sensor", ControlType: "light-intensity"},
		{CapabilityType: "toggle", ControlType: "air-power"},
		{CapabilityType: "adjust", ControlType: "air-temperature"},
		{CapabilityType: "toggle", ControlType: "fan-power"},
		{CapabilityType: "level", ControlType: "fan-speed"},
		{CapabilityType: "sensor", ControlType: "temperature"},
	}
	for _, c := range capabilities {
		db.FirstOrCreate(&c, Capability{
			CapabilityType: c.CapabilityType,
			ControlType:    c.ControlType,
		})
	}
}

func SeedUsers(db *gorm.DB) error {
	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := auth.UserAccount{Username: "testuser", Password: string(hashedPassword)}
	db.FirstOrCreate(&user, auth.UserAccount{Username: user.Username})

	sensorUser := User{Email: "Sensor",Name: "Sensor",}
	db.FirstOrCreate(&sensorUser)
	return nil
}

func SeedDevices(db *gorm.DB) {

	// var livingRoom Room
	// db.Where("room_name = ?", "Living Room").First(&livingRoom)

	// var bedRoom Room
	// db.Where("room_name = ?", "Bedroom").First(&bedRoom)

	devices := []Device{
		{
			DeviceID:      "Light-001",
			DeviceName:    "Main Light",
			DeviceType:    "light",
			RoomID:        nil,
			LastHeartbeat: time.Now(),
		},
		{
			DeviceID:      "Air-001",
			DeviceName:    "Air Conditioner",
			DeviceType:    "air",
			RoomID:        nil,
			LastHeartbeat: time.Now(),
		},
		{
			DeviceID:      "Sensor-001",
			DeviceName:    "Tem Sensor",
			DeviceType:    "temperature-sensor",
			RoomID:        nil,
			LastHeartbeat: time.Now(),
		},
		{
			DeviceID:      "Sensor-001",
			DeviceName:    "Unpaired Sensor",
			DeviceType:    "Sensor",
			RoomID:        nil,
			LastHeartbeat: time.Now(),
		},
	}

	for _, d := range devices {
		err := db.FirstOrCreate(&d, Device{DeviceID: d.DeviceID}).Error
		if err != nil {
			log.Printf("Could not seed device %s: %v", d.DeviceID, err)
		}
	}
}
