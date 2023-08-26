package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/discord-gophers/goapi-gen/examples/petstore-expanded/api"
)

func main() {
	port := flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	client := api.Client{
		BaseURL: "http://localhost:" + strconv.Itoa(*port),
		Client:  http.DefaultClient, // in case we need to do authentication, we can have an oauth2.Client for example
	}

	spotTag := "TagOfSpot"
	// Create a new pet
	newPet := api.NewPet{
		Name: "Spot",
		Tag:  &spotTag,
	}

	// Add the pet
	pet, apiErr, err := client.AddPet(context.Background(), newPet)
	if err != nil {
		panic(err)
	}
	if apiErr != nil {
		panic(apiErr)
	}

	log.Println("Added pet", pet.ID)

	// Get the pet
	pet, apiErr, err = client.FindPetByID(context.Background(), pet.ID)
	if err != nil {
		panic(err)
	}
	if apiErr != nil {
		panic(apiErr)
	}

	log.Println("Found pet", pet.ID, "named", pet.Name)

	pets, apiErr, err := client.FindPets(context.Background(), api.FindPetsParams{
		Tags: []string{spotTag},
	})
	if err != nil {
		panic(err)
	}
	if apiErr != nil {
		panic(apiErr)
	}

	log.Println("Found", len(*pets), "pets")

	for _, pet := range *pets {
		log.Println("Pet", pet.ID, "named", pet.Name)

		// Delete the pet
		apiErr, err = client.DeletePet(context.Background(), pet.ID)
		if err != nil {
			panic(err)
		}

		if apiErr != nil {
			panic(apiErr)
		}

		log.Println("Deleted pet", pet.ID)
	}
}
