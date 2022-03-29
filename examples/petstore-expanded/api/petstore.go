//go:generate go run github.com/discord-gophers/goapi-gen --package=api --generate types,server,spec -o petstore.gen.go ../petstore-expanded.yaml

package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/render"
)

type PetStore struct {
	Lock   sync.Mutex
	Pets   map[int64]Pet
	NextId int64
}

// Make sure we conform to ServerInterface

var _ ServerInterface = (*PetStore)(nil)

func NewPetStore() *PetStore {
	return &PetStore{
		Pets:   make(map[int64]Pet),
		NextId: 1000,
	}
}

const petNotFoundMsg = "Could not find pet with ID %d"

// Here, we implement all of the handlers in the ServerInterface
func (p *PetStore) FindPets(w http.ResponseWriter, r *http.Request, params FindPetsParams) *Response {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	var result []Pet

	for _, pet := range p.Pets {
		if params.Tags != nil {
			// If we have tags,  filter pets by tag
			for _, t := range params.Tags {
				if pet.Tag != nil && (*pet.Tag == t) {
					result = append(result, pet)
				}
			}
		} else {
			// Add all pets if we're not filtering
			result = append(result, pet)
		}

		if params.Limit != nil {
			l := int(*params.Limit)
			if len(result) >= l {
				// We're at the limit
				break
			}
		}
	}

	return FindPetsJSON200Response(result)
}

func (p *PetStore) AddPet(w http.ResponseWriter, r *http.Request) *Response {
	// We expect a NewPet object in the request body.
	var newPet AddPetJSONRequestBody
	if err := render.Bind(r, &newPet); err != nil {
		return AddPetJSONDefaultResponse(Error{"Invalid format for NewPet"}).Status(http.StatusBadRequest)
	}

	// We now have a pet, let's add it to our "database".

	// We're always asynchronous, so lock unsafe operations below
	p.Lock.Lock()
	defer p.Lock.Unlock()

	// We handle pets, not NewPets, which have an additional ID field
	var pet Pet
	pet.Name = newPet.Name
	pet.Tag = newPet.Tag
	pet.ID = p.NextId
	p.NextId = p.NextId + 1

	// Insert into map
	p.Pets[pet.ID] = pet

	// Now, we have to return the NewPet
	return AddPetJSON201Response(pet)
}

func (p *PetStore) FindPetByID(w http.ResponseWriter, r *http.Request, id int64) *Response {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	pet, found := p.Pets[id]
	if !found {
		return FindPetByIDJSONDefaultResponse(Error{fmt.Sprintf(petNotFoundMsg, id)}).Status(http.StatusNotFound)
	}

	return FindPetByIDJSON200Response(pet)
}

func (p *PetStore) DeletePet(w http.ResponseWriter, r *http.Request, id int64) *Response {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, found := p.Pets[id]
	if !found {
		return DeletePetJSONDefaultResponse(Error{fmt.Sprintf(petNotFoundMsg, id)}).Status(http.StatusNotFound)
	}
	delete(p.Pets, id)

	return &Response{Code: http.StatusNoContent}
}
