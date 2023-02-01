package handlers

import (
	"encoding/json"
	filmsdto "finaltask/dto/films"
	dto "finaltask/dto/result"
	"finaltask/models"
	"finaltask/repositories"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"

	"github.com/gorilla/mux"
)

type handlerFilms struct {
	FilmsRepository repositories.FilmsRepository
}

func HandlerFilms(FilmsRepository repositories.FilmsRepository) *handlerFilms {
	return &handlerFilms{FilmsRepository}
}

func convertResponseFilms(u models.Films) models.Films {
	return models.Films{
		ID:    u.ID,
		Title: u.Title,
		Desc:  u.Desc,
		Price: u.Price,
		Genre: u.Genre,
		Image: u.Image,
		YTID:  u.YTID,
		Token: u.Token,
	}
}

func (h *handlerFilms) MakeFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	images := r.Context().Value("cloudImage")
	filename := images.(string)

	price, _ := strconv.Atoi(r.FormValue("price"))
	genreid, _ := strconv.Atoi(r.FormValue("genre_id"))
	request := filmsdto.FilmsRequest{
		Title:   r.FormValue("title"),
		Desc:    r.FormValue("desc"),
		Price:   price,
		YTID:    r.FormValue("ytid"),
		FullUrl: r.FormValue("full_url"),
		Status:  r.FormValue("status"),
		GenreID: genreid,
		Image:   filename,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "Form not fully filled!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	films := models.Films{
		Title:   request.Title,
		Desc:    request.Desc,
		Price:   request.Price,
		GenreID: request.GenreID,
		YTID:    request.YTID,
		FullUrl: request.FullUrl,
		Status:  request.Status,
		Image:   filename,
	}

	films, err = h.FilmsRepository.AddFilms(films)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseFilms(films)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFilms) AllFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	films, err := h.FilmsRepository.GetFilms()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: films}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFilms) GetFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	films, err := h.FilmsRepository.FindAFilm(id, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: films}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFilms) GetPublicFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	films, err := h.FilmsRepository.FindOneFilm(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseFilms(films)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFilms) ForBanner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	films, err := h.FilmsRepository.GetTopFilms()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: films}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFilms) EditFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	images := r.Context().Value("cloudImage")
	filename := images.(string)
	fmt.Println(images, "file", filename)

	price, _ := strconv.Atoi(r.FormValue("price"))
	genreid, _ := strconv.Atoi(r.FormValue("genre_id"))
	request := filmsdto.FilmsRequest{
		Title:   r.FormValue("title"),
		Desc:    r.FormValue("desc"),
		Price:   price,
		YTID:    r.FormValue("ytid"),
		FullUrl: r.FormValue("full_url"),
		Status:  r.FormValue("status"),
		GenreID: genreid,
		Image:   filename,
	}

	ID, _ := strconv.Atoi(mux.Vars(r)["id"])

	films, err := h.FilmsRepository.FindOneFilm(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	// films := models.Films{}

	if request.Title != "" {
		films.Title = request.Title
	}

	if request.Desc != "" {
		films.Desc = request.Desc
	}

	if request.YTID != "" {
		films.YTID = request.YTID
	}

	if request.FullUrl != "" {
		films.FullUrl = request.FullUrl
	}

	if request.Price != 0 {
		films.Price = request.Price
	}

	if request.Status != "" {
		films.Status = request.Status
	}

	if request.Image != "" {
		films.Image = request.Image
	}

	if request.GenreID != 0 {
		films.GenreID = request.GenreID
	}

	data, err := h.FilmsRepository.EditFilms(films, ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, _ = h.FilmsRepository.FindOneFilm(data.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseFilms(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerFilms) DeleteFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.FilmsRepository.FindOneFilm(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.FilmsRepository.DeleteFilms(user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseFilms(data)}
	json.NewEncoder(w).Encode(response)
}
