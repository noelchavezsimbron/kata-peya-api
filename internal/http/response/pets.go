package response

import (
	"kata-peya/internal/pet"
)

type Pet struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Vaccines []string `json:"vaccines"`
	Age      string   `json:"age"`
} //@name PetResponse

func Pets(pets []pet.Pet) []Pet {
	r := make([]Pet, len(pets))

	for i, p := range pets {

		vaccines := make([]string, len(p.Vaccines))
		for j, v := range p.Vaccines {
			vaccines[j] = string(v)
		}

		r[i] = Pet{
			Id:       p.Id,
			Name:     p.Name,
			Vaccines: vaccines,
			Age:      p.GetAgeFormatted(),
		}
	}
	return r
}
