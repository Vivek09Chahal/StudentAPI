package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Vivek09Chahal/student_api/internal/storage"
	"github.com/Vivek09Chahal/student_api/internal/types"
	response "github.com/Vivek09Chahal/student_api/internal/utils/repsonse"
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
        }

        response.WriteJSON(w, http.StatusCreated, map[string]int64{"id":  lastID})
    }
}