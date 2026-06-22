package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Sameetpatro/go-crud/internal/storage"
	"github.com/Sameetpatro/go-crud/internal/types"
	"github.com/Sameetpatro/go-crud/internal/utils/response"
	"github.com/go-playground/validator/v10"
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

func GetByID(storage storage.Storage) http.HandlerFunc{
	return func(res http.ResponseWriter, req *http.Request){
		id := req.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))
		intid, err := strconv.ParseInt(id, 10, 64)
		if err != nil{
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, e := storage.GetStudentById(intid)
		if e != nil{
			slog.Error("Error getting User", slog.String("id", id))
			response.WriteJson(res, http.StatusInternalServerError, response.GeneralError(e))
			return 
		}
		response.WriteJson(res, http.StatusAccepted, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc{
	return func(res http.ResponseWriter, req *http.Request){
		slog.Info("Getting all students")
		students, errr:= storage.GetStudents()
		if errr != nil{
			response.WriteJson(res, http.StatusInternalServerError, errr)
			return
		}
		response.WriteJson(res, http.StatusAccepted, students)
	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc{
	return func(res http.ResponseWriter, req *http.Request){
		id := req.PathValue("id")
		slog.Info("Updating the Student with requested ID")
		intid, err := strconv.ParseInt(id, 10, 64)
		if err != nil{
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		var student types.Student
		er := json.NewDecoder(req.Body).Decode(&student)
		if errors.Is(er, io.EOF){
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(fmt.Errorf("body cant be empty")))
			return
		}
		if er != nil{
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid body")))
			return
		}
		// internal/http/handlers/students/student.go

		if e := validator.New().Struct(student); e != nil {
			var validationErrs validator.ValidationErrors
			if errors.As(e, &validationErrs) {
				response.WriteJson(res, http.StatusBadRequest, response.ValidationError(validationErrs))
			} else {
				response.WriteJson(res, http.StatusBadRequest, response.GeneralError(err))
			}
			return
		}
		updated, err := storage.UpdateTheStudent(intid, student.Name, student.Email, student.Age)
		if err != nil{
			response.WriteJson(res, http.StatusBadGateway, response.GeneralError(err))
			return
		}
		slog.Info("updation done bhai")
		response.WriteJson(res, http.StatusOK, updated)

	}
}

func DeleteStudent(storage storage.Storage) http.HandlerFunc{
	return func(res http.ResponseWriter, req *http.Request){
		id := req.PathValue("id")
		slog.Info("Deleting a student...")
		intid, err := strconv.ParseInt(id, 10, 64)
		if err != nil{
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		e := storage.DeleteTheStudent(intid)
		if e != nil{
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(e))
			return
		}
		slog.Info("student deleted successfully", slog.String("id", id))
		response.WriteJson(res, http.StatusOK, map[string]string{"status": "ok", "id": id})
	}
}
