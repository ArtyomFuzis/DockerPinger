package httpserver

import (
	"backend/database"
	"backend/httpserver/transfer"
	"encoding/json"
	"log"
	"net/http"
)

func getServicesInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	repo := database.GetPingRepository()
	services := repo.GetServices()
	w.Header().Set("Content-Type", "application/json")
	res := make([]transfer.PingServiceTransferObject, len(services))
	for i, service := range services {
		res[i] = transfer.PingServiceTransferObject{
			Address:  service.Address,
			LastPing: transfer.PingTransferObject{Date: repo.GetLastPing(service.Address).Date},
		}
	}
	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Unable to convert services list to json: ", err)
	} else {
		_, err = w.Write(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Unable to send services list: ", err)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}

}
func addService(w http.ResponseWriter, r *http.Request) {

}
func addPing(w http.ResponseWriter, r *http.Request) {

}
