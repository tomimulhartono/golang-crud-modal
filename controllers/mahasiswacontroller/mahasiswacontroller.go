package mahasiswacontroller

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"gitlab.com/tomimulhartono/go-crud-modal/entities"
	"gitlab.com/tomimulhartono/go-crud-modal/models/mahasiswamodel"
)

var mahasiswaModel = mahasiswamodel.New()

func Index(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"data": template.HTML(GetData()),
	}

	temp, _ := template.ParseFiles("views/mahasiswa/index.html")
	temp.Execute(w, data)
}

func GetData() string {
	buffer := &bytes.Buffer{}

	temp, _ := template.New("data.html").Funcs(template.FuncMap{
		"increment": func(a int, b int) int {
			return a + b
		},
	}).ParseFiles("views/mahasiswa/data.html")

	var mahasiswa []entities.Mahasiswa
	err := mahasiswaModel.FindAll(&mahasiswa)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"mahasiswa": mahasiswa,
	}

	temp.ExecuteTemplate(buffer, "data.html", data)

	return buffer.String()
}

func GetForm(w http.ResponseWriter, r *http.Request) {

	queryString := r.URL.Query()
	id, err := strconv.ParseInt(queryString.Get("id"), 10, 64)

	var data map[string]interface{}
	var mahasiswa entities.Mahasiswa

	if err != nil {
		data = map[string]interface{}{
			"title":     "Add Data Mahasiswa",
			"mahasiswa": mahasiswa,
		}
	} else {

		err := mahasiswaModel.Find(id, &mahasiswa)
		if err != nil {
			panic(err)
		}

		data = map[string]interface{}{
			"title":     "Edit Data Mahasiswa",
			"mahasiswa": mahasiswa,
		}
	}

	temp, _ := template.ParseFiles("views/mahasiswa/form.html")
	temp.Execute(w, data)
}

func Store(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		r.ParseForm()
		var mahasiswa entities.Mahasiswa

		mahasiswa.FullName = r.Form.Get("txt_fullname")
		mahasiswa.Gender = r.Form.Get("txt_gender")
		mahasiswa.Birthplace = r.Form.Get("txt_birthplace")
		mahasiswa.Birthdate = r.Form.Get("txt_birthdate")
		mahasiswa.Address = r.Form.Get("txt_address")

		id, err := strconv.ParseInt(r.Form.Get("txt_id"), 10, 64)

		var data map[string]interface{}

		if err != nil {
			// insert data
			err := mahasiswaModel.Create(&mahasiswa)
			if err != nil {
				ResponseError(w, http.StatusInternalServerError, err.Error())
				return
			}
			data = map[string]interface{}{
				"message": "Data berhasil disimpan",
				"data":    template.HTML(GetData()),
			}
		} else {
			// update data
			mahasiswa.Id = id
			err := mahasiswaModel.Update(id, &mahasiswa)
			if err != nil {
				ResponseError(w, http.StatusInternalServerError, err.Error())
				return
			}
			data = map[string]interface{}{
				"message": "Data berhasil diubah",
				"data":    template.HTML(GetData()),
			}
		}

		ResponseJson(w, http.StatusOK, data)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)

	if err != nil {
		panic(err)
	}

	err = mahasiswaModel.Delete(id)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"message": "Data berhasil dihapus.",
		"data":    template.HTML(GetData()),
	}

	ResponseJson(w, http.StatusOK, data)

}

func ResponseError(w http.ResponseWriter, code int, message string) {
	ResponseJson(w, code, map[string]string{"error": message})
}

func ResponseJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
