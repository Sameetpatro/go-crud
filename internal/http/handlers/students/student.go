package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Sameetpatro/go-crud/internal/types"
	"github.com/Sameetpatro/go-crud/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"github.com/Sameetpatro/go-crud/internal/storage"
)


func New(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request){
		
		var student types.Student

		slog.Info("Creating a student")

		err := json.NewDecoder(req.Body).Decode(&student)
		if errors.Is(err, io.EOF){
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(fmt.Errorf("body must not be empty")))
			return
		}
		if err != nil{
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid body")))
			return
		}

		//request validation
		if err := validator.New().Struct(student); err != nil{
			response.WriteJson(res, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}
		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil{
			response.WriteJson(res, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to create student")))
			return
		}
		slog.Info("user created successfully", slog.String("user_id", fmt.Sprint(lastId)))

		response.WriteJson(res, http.StatusCreated, map[string]string{"status": "Fine babe", "id": fmt.Sprint(lastId)})
	}
}