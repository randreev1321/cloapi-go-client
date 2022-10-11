package recipe

import (
	"encoding/json"
	"github.com/clo-ru/cloapi-go-client/clo"
	"io"
)

type RecipeListResponse struct {
	clo.Response
	Count   int              `json:"count"`
	Results []RecipeListItem `json:"results"`
}

func (r *RecipeListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type RecipeListItem struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	MinDisk        int      `json:"min_disk"`
	MinRam         int      `json:"min_ram"`
	MinVcpus       int      `json:"min_vcpus"`
	SuitableImages []string `json:"suitable_images"`
}
