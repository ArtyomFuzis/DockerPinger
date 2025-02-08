package database

import (
	"gorm.io/gorm"
	"time"
)

type PingRepository struct {
	pingConnection *gorm.DB
}

func (pingRepository *PingRepository) AddService(address string) {
	newService := PingedServices{Address: address}
	pingRepository.pingConnection.Create(&newService)
}
func (pingRepository *PingRepository) AddPingByAddress(address string, date time.Time) {
	var service PingedServices
	pingRepository.pingConnection.Where("address = ?", address).Take(&service)
	newPing := Ping{Date: date, ServiceId: service.ID}
	pingRepository.pingConnection.Create(&newPing)
}
func (pingRepository *PingRepository) GetLastPing(address string) Ping {
	var service PingedServices
	pingRepository.pingConnection.Where("address = ?", address).Take(&service)
	var ping Ping
	pingRepository.pingConnection.Where("service_id = ?", service.ID).Order("date desc").First(&ping)
	return ping
}

func (pingRepository *PingRepository) GetServices() []PingedServices {
	var services []PingedServices
	pingRepository.pingConnection.Find(&services)
	return services
}
