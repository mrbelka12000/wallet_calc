package v1

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func (c *Controller) UploadStatement(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		c.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Retrieve the file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Now let’s save it locally
	dst, err := createFile(handler.Filename)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	defer func() {
		if err := os.Remove(handler.Filename); err != nil {
			c.log.With("error", err).Error("delete file")
		}
	}()

	// Copy the uploaded file to the destination file
	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
	}

	out, err := c.transactionUsecase.ParseStatement(r.Context(), handler.Filename, uuid.MustParse("a6fb4757-62dd-48e5-acbd-a91e0fe9d3ff"))
	if err != nil {
		c.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(out); err != nil {
		c.errorResponse(w, err, http.StatusBadRequest)
		return
	}
}

func createFile(filename string) (*os.File, error) {
	// Build the file path and create it
	dst, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return dst, nil
}
