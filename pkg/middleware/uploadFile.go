package middleware

import (
	"context"
	"encoding/json"
	dto "finaltask/dto/result"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadFile(point http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			file, header, err := r.FormFile("image")

			if err != nil {
				response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: "Form not fully filled!"}
				json.NewEncoder(w).Encode(response)
				return
			}
			defer file.Close()
			const MAX_UPLOAD_SIZE = 32 << 20 // 32MB

			r.ParseMultipartForm(MAX_UPLOAD_SIZE)
			if r.ContentLength > MAX_UPLOAD_SIZE {
				w.WriteHeader(http.StatusBadRequest)
				response := Result{Code: http.StatusBadRequest, Message: "Files sizes are too big!"}
				json.NewEncoder(w).Encode(response)
				return
			}

			if filepath.Ext(header.Filename) != ".jpg" && filepath.Ext(header.Filename) != ".jpeg" && filepath.Ext(header.Filename) != ".png" {
				w.WriteHeader(http.StatusBadRequest)
				response := dto.ErrorResult{
					Code:    400,
					Message: "The provided file format is not allowed. Please upload a JPG, JPEG or PNG image",
				}
				json.NewEncoder(w).Encode(response)
				return
			}

			tempFile, err := os.CreateTemp("uploads", "image-*.jpeg")
			if err != nil {
				fmt.Println(err)
				fmt.Println("path upload error")
				json.NewEncoder(w).Encode(err)
				return
			}
			defer tempFile.Close()

			var ctx = context.Background()
			cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUD_NAME"), os.Getenv("API_KEY"), os.Getenv("API_SECRET"))

			resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: "CinemaOnline"})
			if err != nil {
				fmt.Println(err.Error())
			}
			contx := context.WithValue(r.Context(), "cloudImage", resp.SecureURL)

			point.ServeHTTP(w, r.WithContext(contx))
		})
}

func MayUploadFile(point http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			file, _, _ := r.FormFile("image")

			if file == nil {
				filename := ""
				ctx := context.WithValue(r.Context(), "cloudImage", filename)
				point.ServeHTTP(w, r.WithContext(ctx))
				return
			} else {
				defer file.Close()
				const MAX_UPLOAD_SIZE = 32 << 20 // 32MB

				r.ParseMultipartForm(MAX_UPLOAD_SIZE)
				if r.ContentLength > MAX_UPLOAD_SIZE {
					w.WriteHeader(http.StatusBadRequest)
					response := Result{Code: http.StatusBadRequest, Message: "Files sizes are too big!"}
					json.NewEncoder(w).Encode(response)
					return
				}

				tempFile, err := os.CreateTemp("uploads", "image-*.jpeg")
				if err != nil {
					fmt.Println(err)
					fmt.Println("path upload error")
					json.NewEncoder(w).Encode(err)
					return
				}
				defer tempFile.Close()

				var ctx = context.Background()
				cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUD_NAME"), os.Getenv("API_KEY"), os.Getenv("API_SECRET"))

				resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: "CinemaOnline"})
				if err != nil {
					fmt.Println(err.Error())
				}

				contx := context.WithValue(r.Context(), "cloudImage", resp.SecureURL)

				point.ServeHTTP(w, r.WithContext(contx))
			}
		})
}
