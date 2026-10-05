package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Vivek09Chahal/student_api/internal/storage"
	"github.com/Vivek09Chahal/student_api/internal/types"
	"github.com/Vivek09Chahal/student_api/internal/utils/repsonse"
	"github.com/go-playground/validator"
)

func New(storage storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {


        var students types.Student

        err := json.NewDecoder(r.Body).Decode(&students)

        if errors.Is(err, io.EOF) {
            response.WriteJSON(w, http.StatusBadRequest, err.Error())
            return
        }
        
        slog.Info("creating a student")
        // w.Write([]byte("Welcom to student API" ))

        //request validate
        if err := validator.New().Struct(students); err != nil {
            validateErrs := err.(validator.ValidationErrors)
            response.WriteJSON(w, http.StatusBadRequest, response.ValidatonError(validateErrs))
            return 
        }

        lastID, err := storage.CreateStudent(students.Name, students.Email, students.Age)
        slog.Info("user created successfully", slog.String("UserID", fmt.Sprint(lastID)))
            
        if err != nil {
            response.WriteJSON(w, http.StatusInternalServerError, err)
            return
        }

        response.WriteJSON(w, http.StatusCreated, map[string]int64{"id":  lastID})
    }
}

func GetById(storage storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        id := r.PathValue("id")
        
        slog.Info("getting a student", slog.String("id", id))

        intID, err := strconv.ParseInt(id, 10, 64)
        if err != nil {
            response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
            return 
        }
        student, err := storage.GetStudentById(intID)

        if err != nil {
            slog.Error("err getting user", slog.String("id", id))
            response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
            return
        }

        response.WriteJSON(w, http.StatusOK, student)
    }
}

func GetList(storage storage.Storage) http.HandlerFunc  {
    return func(w http.ResponseWriter, r *http.Request) {
        slog.Info("getting all students")
        
        students, err := storage.GetStudents()

        if err != nil {
            response.WriteJSON(w, http.StatusInternalServerError, err)
            return
        }

        response.WriteJSON(w, http.StatusOK, students)
        
    }
}