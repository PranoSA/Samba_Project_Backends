package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PranoSA/samba_share_backend/web_server/models"
	"github.com/julienschmidt/httprouter"
)

type CreateSpaceRequest struct {
	Megabytes int64
}

func (ar AppRouter) CreateSpace(w http.ResponseWriter, r *http.Request, pa httprouter.Params) {

	w.Header().Add("Content-Type", "application/json")

	user, ok := r.Context().Value("Authorization").(string)

	if !ok {
		log.Fatal("Invalid Type Casting at Spaces.go /Create Space")
	}

	var request CreateSpaceRequest

	e := json.NewDecoder(r.Body).Decode(&request)

	//megabytes, e := strconv.Atoi(r.URL.Query().Get("megabytes"))

	if e != nil {
		json.NewEncoder(w)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	megabytes := request.Megabytes

	if megabytes < 10 || megabytes > 100_000 {

		var ErrorResponse BadRequestResponse = BadRequestResponse{
			{
				ParameterName: "megabytes",
				Param_Type:    "Query",
				Value_Type:    "int64",
				Message:       "Values Between 10 and 100,000 needed ",
			},
		}

		json.NewEncoder(w).Encode(&ErrorResponse)
		return
	}

	res, err := ar.Models.Spaces.CreateSpace(models.SpaceRequest{
		Owner:     user,
		Megabytes: int64(megabytes),
	})

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	json.NewEncoder(w).Encode(&res)

}

func (ar AppRouter) DeleteSpace(w http.ResponseWriter, r *http.Request, ap httprouter.Params) {
	spaceid := ap.ByName("spaceid")
	email := r.Context().Value("Authorization")
	if email == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	res, err := ar.Models.Spaces.GetSpaceById(models.DeleteSpaceRequest{
		Owner:    email.(string),
		Space_id: spaceid,
	})
	if err == models.ErrorEntryDoesNotExist {
		w.WriteHeader(http.StatusForbidden)
	}

	json.NewEncoder(w).Encode(res)
}

func (ar AppRouter) GetMySpaces(w http.ResponseWriter, r *http.Request, ap httprouter.Params) {

	email := r.Context().Value("Authorization")
	res, err := ar.Models.Spaces.GetSpaceByOwner(email.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	json.NewEncoder(w).Encode(res)
}
