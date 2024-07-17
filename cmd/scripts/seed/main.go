package main

import (
	"encoding/json"
	"future-app/models"
	"future-app/store"
	"log"
	"os"
	"time"
)

func main() {
	byteValue, err := os.ReadFile("appointments.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var appointments []models.Appointment
	if err = json.Unmarshal(byteValue, &appointments); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	dbStore, err := store.NewStore()
	if err != nil {
		log.Fatalf("Error creating store: %v", err)
	}
	defer dbStore.Close()
	if err := dbStore.Init(); err != nil {
		log.Fatalf("Error initializing store: %v", err)
	}

	for _, appointment := range appointments {
		query := `
        INSERT INTO appointments (id, user_id, trainer_id, started_at, ended_at)
        VALUES ($1, $2, $3, $4, $5)
        `

		_, err := dbStore.DB.Exec(
			query,
			appointment.ID,
			appointment.UserID,
			appointment.TrainerID,
			(appointment.StartedAt).Format(time.RFC3339),
			(appointment.EndedAt).Format(time.RFC3339),
		)

		if err != nil {
			log.Fatalf("Error creating appointment: %v", err)
		}
	}
}
