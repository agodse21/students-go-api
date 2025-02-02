package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/agodse21/students-go-api/internal/storage"
	"github.com/agodse21/students-go-api/internal/types"
	"github.com/agodse21/students-go-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a new student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//Request validation

		if err := validator.New().Struct(student); err != nil {
			validateErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErr))
			return
		}

		lastId, err := storage.CreateStudent(
			student,
		)

		slog.Info("Student created", slog.String("id", fmt.Sprint(lastId)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "student created", "id": fmt.Sprint(lastId)})
		// w.Write([]byte("Welcome to Students API"))
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		slog.Info("Geting a student by id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			slog.Error("Error parsing id", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, er := storage.GetStudentById(intId)

		if er != nil {
			slog.Error("Error getting student by id", slog.String("id", id), slog.String("error", er.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(er))
			return
		}

		response.WriteJson(w, http.StatusOK, student)

	}
}

func GetAll(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Geting all students")
		students, err := storage.GetStudents()

		if err != nil {
			slog.Error("Error getting all students", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func Delete(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Deleting a student")

		id := r.PathValue("id")

		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			slog.Error("Error parsing id", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := storage.DeleteStudentById(intId); err != nil {
			slog.Error("Error deleting student", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "student deleted"})

	}
}

func Update(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Updating a student")

		id := r.PathValue("id")

		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			slog.Error("Error parsing id", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		var student types.Student

		err = json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			slog.Info("Request body is empty")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if err != nil {
			slog.Info("Error decoding request body", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validateErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErr))
			return
		}

		if err := storage.UpdateStudentById(intId, student); err != nil {
			slog.Error("Error updating student", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "student updated"})
	}
}
