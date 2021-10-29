package server

import (
	"dockerized-rest-api/util/logger"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

type Server struct {
	logger    *logger.Logger
	db        *gorm.DB
	validator *validator.Validate
}

const (
	appErrDataAccessFailure      = "data access failure"
	appErrJsonCreationFailure    = "json creation failure"
	appErrDataCreationFailure    = "data creation failure"
	appErrFormDecodingFailure    = "form decoding failure"
	appErrDataUpdateFailure      = "data update failure"
	appErrFormErrResponseFailure = "form error response failure"
)

func New(logger *logger.Logger, db *gorm.DB, validator *validator.Validate) *Server {
	return &Server{logger: logger, db: db, validator: validator}
}

func (server *Server) Logger() *logger.Logger {
	return server.logger
}
