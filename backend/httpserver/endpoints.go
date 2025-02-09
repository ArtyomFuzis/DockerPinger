package httpserver

import (
	messaging "backend/amqp"
	"backend/database"
	"backend/transfer"
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
		lastPing := repo.GetLastPing(service.Address)
		res[i] = transfer.PingServiceTransferObject{
			Address:     service.Address,
			State:       lastPing.State,
			LastPing:    transfer.PingTransferObject{Date: lastPing.Date},
			LastSuccess: transfer.PingTransferObject{Date: repo.GetLastSuccessPing(service.Address).Date},
		}
	}
	data, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Unable to convert services list to json:", err)
	} else {
		_, err = w.Write(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Unable to send services list:", err)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}

}
func servicesAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Unable to parse form:", err)
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte("Unable to parse form"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Unable to send response msg:", err)
		}
		return
	}
	address := r.Form.Get("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Not found \"address\" in request body")
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte("Not found \"address\" in request body"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Unable to send response msg:", err)
		}
		return
	} else {
		repo := database.GetPingRepository()
		repo.AddService(address)
		messaging.SendToAddService(address, transfer.Add)
		w.WriteHeader(http.StatusOK)
	}
}
func servicesDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Unable to parse form:", err)
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte("Unable to parse form"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Unable to send response msg:", err)
		}
		return
	}
	address := r.Form.Get("address")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Not found \"address\" in request body")
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte("Not found \"address\" in request body"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Unable to send response msg:", err)
		}
		return
	} else {
		repo := database.GetPingRepository()
		repo.DeleteService(address)
		messaging.SendToAddService(address, transfer.Delete)
		w.WriteHeader(http.StatusOK)
	}
}
