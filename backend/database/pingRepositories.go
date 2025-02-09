package database

import (
	"gorm.io/gorm"
	"time"
)

type PingRepository struct {
	pingConnection *gorm.DB
}

func (pingRepository *PingRepository) AddService(address string) {
	query := conn.Unscoped().Model(&PingedServices{}).Where("address = ?", address)
	if query.RowsAffected == 0 {
		query.Update("deleted_at", nil)
	} else {
		newService := PingedServices{Address: address}
		pingRepository.pingConnection.Create(&newService)
	}
}
func (pingRepository *PingRepository) DeleteService(address string) {
	var service PingedServices
	pingRepository.pingConnection.Where("address = ?", address).Take(&service)
	pingRepository.pingConnection.Delete(&service)
}
func (pingRepository *PingRepository) AddPingByAddress(address string, date time.Time, state bool) {
	var service PingedServices
	pingRepository.pingConnection.Where("address = ?", address).Take(&service)
	newPing := Ping{Date: date, ServiceId: service.ID, State: state}
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

func (pingRepository *PingRepository) GetLastSuccessPing(address string) Ping {
	var service PingedServices
	pingRepository.pingConnection.Where("address = ?", address).Take(&service)
	var ping Ping
	query := pingRepository.pingConnection.Where("service_id = ? AND state = true", service.ID).Order("date desc")
	if query.RowsAffected == 0 {
		return ping
	}
	query.First(&ping)
	return ping
}
