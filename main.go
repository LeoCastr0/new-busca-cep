package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	http.HandleFunc("/", FindZipHandler)
	http.ListenAndServe(":8080", nil)
}

func FindZipHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	zipParam := r.URL.Query().Get("zip")
	if zipParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	zip, err := FindZip(zipParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// This lines below do the same thing as the encoder does
	// result, err := json.Marshal(zip)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(result)
	json.NewEncoder(w).Encode(zip)
}

func FindZip(zip string) (*ViaCEP, error) {
	res, err := http.Get("https://viacep.com.br/ws/" + zip + "/json/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
