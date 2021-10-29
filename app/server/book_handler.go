package server

import (
	"dockerized-rest-api/model"
	"dockerized-rest-api/repository"
	"dockerized-rest-api/util/validator"
	"encoding/json"
	"fmt"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"

	"net/http"
	"strconv"
)

func (server *Server) HandleListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := repository.ListBooks(server.db)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}
	if books == nil {
		fmt.Fprint(w, "[]")
		return
	}
	dtos := books.ToDto()
	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (server *Server) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	form := &model.BookForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := server.validator.Struct(form); err != nil {
		server.logger.Warn().Err(err).Msg("")
		resp := validator.ToErrResponse(err)
		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrFormErrResponseFailure)
			return
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			server.logger.Warn().Err(err).Msg("")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	if err := server.validator.Struct(form); err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
		return
	}

	bookModel, err := form.ToModel()
	if err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}
	book, err := repository.CreateBook(server.db, bookModel)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataCreationFailure)
		return
	}
	server.logger.Info().Msgf("New book created: %d", book.ID)
	w.WriteHeader(http.StatusCreated)
}

func (server *Server) HandleReadBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		server.logger.Info().Msgf("can not parse ID: %v", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	book, err := repository.ReadBook(server.db, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}
	dto := book.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (server *Server) HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		server.logger.Info().Msgf("can not parse ID: %v", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	form := &model.BookForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := server.validator.Struct(form); err != nil {
		server.logger.Warn().Err(err).Msg("")
		resp := validator.ToErrResponse(err)
		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrFormErrResponseFailure)
			return
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			server.logger.Warn().Err(err).Msg("")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	if err := server.validator.Struct(form); err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
		return
	}

	bookModel, err := form.ToModel()
	if err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}
	bookModel.ID = uint(id)
	if err := repository.UpdateBook(server.db, bookModel); err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataUpdateFailure)
		return
	}
	server.logger.Info().Msgf("Book updated: %d", id)
	w.WriteHeader(http.StatusAccepted)
}

func (server *Server) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		server.logger.Info().Msgf("can not parse ID: %v", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err := repository.DeleteBook(server.db, uint(id)); err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}
	server.logger.Info().Msgf("Book deleted: %d", id)
	w.WriteHeader(http.StatusAccepted)
}
