package main

import (
	"context"
	"day7/connection"
	"fmt"
	"html/template"
	"log"

	"math"
	"net/http"
	"strconv"

	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	connection.ConnectionProject()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.PathPrefix("/node_modules/").Handler(http.StripPrefix("/node_modules/", http.FileServer(http.Dir("./node_modules"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/form-project", formProject).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/detail-project/{id}", detailProject).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	route.HandleFunc("/update-project/{id}", updateProject).Methods("POST")

	fmt.Println("server on")
	http.ListenAndServe("localhost:5000", route)
}

type Blog struct {
	Id              int
	NameProject     string
	StarDate        time.Time
	Format_starDate string
	Edit_starDate   string
	EndDate         time.Time
	Format_endDate  string
	Edit_endDate    string
	Duration        string
	Message         string
	Tech            []string
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	var nameProject = r.PostForm.Get("projectName")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")
	var desc = r.PostForm.Get("Description")

	var tech []string
	for key, values := range r.Form {
		for _, value := range values {
			if key == "technologies" {
				tech = append(tech, value)
			}
		}
	}

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name,start_date,end_date,description,technologies) VALUES ($1,$2,$3,$4,$5)", nameProject, startDate, endDate, desc, tech)

	if err != nil {
		w.Write([]byte("Error baris 78 " + err.Error()))
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("view/index.html")

	if err != nil {
		w.Write([]byte("messege: " + err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies FROM tb_projects")

	var result []Blog

	for data.Next() {

		var each = Blog{}
		err := data.Scan(&each.Id, &each.NameProject, &each.StarDate, &each.EndDate, &each.Message, &each.Tech)

		if err != nil {
			fmt.Println("Error baris 126 " + err.Error())
			return
		}

		hs := each.EndDate.Sub(each.StarDate).Hours()
		day, _ := math.Modf(hs / 24)
		bulan := int64(day / 30)
		tahun := int64(day / 365)

		if tahun > 0 {
			each.Duration = strconv.FormatInt(tahun, 10) + " Year"
		} else if bulan > 0 {
			each.Duration = strconv.FormatInt(bulan, 10) + " Month"
		} else {
			each.Duration = fmt.Sprintf("%.0f", day) + " Day"
		}
		result = append(result, each)
	}

	response := map[string]interface{}{
		"Blogs": result,
	}

	tmpl.Execute(w, response)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("view/blog-detail.html")

	if err != nil {
		w.Write([]byte("messege: " + err.Error()))
		return
	}

	var ProjectDetail = Blog{}

	indexData, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies FROM tb_projects WHERE id=$1", indexData).Scan(&ProjectDetail.Id, &ProjectDetail.NameProject, &ProjectDetail.StarDate, &ProjectDetail.EndDate, &ProjectDetail.Message, &ProjectDetail.Tech)

	if err != nil {
		fmt.Println("Error baris 168 " + err.Error())
		return
	}
	hs := ProjectDetail.EndDate.Sub(ProjectDetail.StarDate).Hours()
	day, _ := math.Modf(hs / 24)
	bulan := int64(day / 30)
	tahun := int64(day / 365)

	if tahun > 0 {
		ProjectDetail.Duration = strconv.FormatInt(tahun, 10) + " Year"
	} else if bulan > 0 {
		ProjectDetail.Duration = strconv.FormatInt(bulan, 10) + " Month"
	} else {
		ProjectDetail.Duration = fmt.Sprintf("%.0f", day) + " Day"
	}

	ProjectDetail.Format_starDate = ProjectDetail.StarDate.Format("2 January 2006")
	ProjectDetail.Format_endDate = ProjectDetail.EndDate.Format("2 January 2006")

	response := map[string]interface{}{
		"Blogs": ProjectDetail,
	}

	tmpl.Execute(w, response)
}

func formProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("view/add.html")

	if err != nil {
		w.Write([]byte("messege: " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("view/contact.html")

	if err != nil {
		w.Write([]byte("messege: " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	indexDelete, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", indexDelete)

	if err != nil {
		w.Write([]byte("eror Baris 201 " + err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("view/edit-project.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var EditProject = Blog{}
	indexData, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies FROM tb_projects WHERE id=$1", indexData).Scan(&EditProject.Id, &EditProject.NameProject, &EditProject.StarDate, &EditProject.EndDate, &EditProject.Message, &EditProject.Tech)

	if err != nil {
		w.Write([]byte("error baris 223 " + err.Error()))
	}

	EditProject.Edit_starDate = EditProject.StarDate.Format("2006-01-02")
	EditProject.Edit_endDate = EditProject.EndDate.Format("2006-01-02")

	dataEdit := map[string]interface{}{
		"Edit": EditProject,
	}

	tmpl.Execute(w, dataEdit)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	indexData, _ := strconv.Atoi(mux.Vars(r)["id"])

	var nameProject = r.PostForm.Get("projectName")
	var startDate = r.PostForm.Get("startDate")
	var endDate = r.PostForm.Get("endDate")
	var desc = r.PostForm.Get("Description")

	var tech []string
	for key, values := range r.Form {
		for _, value := range values {
			if key == "technologies" {
				tech = append(tech, value)
			}
		}
	}

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name=$2,start_date=$3,end_date=$4,description=$5,technologies=$6 WHERE id=$1", indexData, nameProject, startDate, endDate, desc, tech)

	if err != nil {
		w.Write([]byte("error Baris 267: " + err.Error()))
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
