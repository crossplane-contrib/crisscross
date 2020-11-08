/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hasheddan/crisscross/pkg/models"
)

func observe(w http.ResponseWriter, r *http.Request) {
	log.Print("received observation")
	b := models.ObservationResponse{
		External: models.ExternalObservation{
			ResourceExists: true,
		},
	}
	byt, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(byt)
}

func create(w http.ResponseWriter, r *http.Request) {
	log.Print("received create")
	b := models.CreationResponse{
		External: models.ExternalCreation{},
	}
	byt, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(byt)
}

func update(w http.ResponseWriter, r *http.Request) {
	log.Print("received update")
	b := models.UpdateResponse{
		External: models.ExternalUpdate{},
	}
	byt, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(byt)
}

func delete(w http.ResponseWriter, r *http.Request) {
	log.Print("received delete")
	b := models.DeletionResponse{}
	byt, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(byt)
}

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/observe", observe)
	http.HandleFunc("/create", create)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
