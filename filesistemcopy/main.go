package main

import (
	//"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"

	"os"
	//"strings"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Bicom-Systems-Hackathon/goated-with-the-source.git/logic"
)

// var i *template.Template = template.Must(template.ParseFiles("./templates/index.html"))
var s *template.Template = template.Must(template.ParseFiles("./templates/search.html"))
var e *template.Template = template.Must(template.ParseFiles("./templates/edit.html"))
var c *template.Template = template.Must(template.ParseFiles("./templates/create.html"))
var baza = logic.GoMap{}

func main() {
	fmt.Println("hello")
	//dodajemo novi tag ili paragraf createTag, createItem
	http.HandleFunc("/edit/", editHandlerFunc)     //mijenjamo item
	http.HandleFunc("/delete/", deleteHandlerFunc) //brisemo item
	http.HandleFunc("/editResult/", editResultFunction)
	http.HandleFunc("/create", createHandlerFunc)
	http.HandleFunc("/read", readHandlerFunc)
	ime := "nesto"
	content := "content"
	lista := []string{"kupus", "voda"}
	item1 := logic.Item1{Name: ime, Index: len(baza.Item_list), Tag: lista}
	item1.Content = &content
	baza.Add(item1)
	item2 := logic.Item1{Name: "teta", Content: &ime, Index: len(baza.Item_list), Tag: lista}
	baza.Add(item2)
	http.ListenAndServe("localhost:8080", nil)
}
func editResultFunction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	content := r.FormValue("content")
	tags := r.FormValue("tags")
	index := r.URL.Query().Get("Index")
	fmt.Println("editResultFunction")
	fmt.Println("INDEKS JE: ", index)
	i, err := strconv.Atoi(index)

	if err != nil {
		panic(err)
	}
	tag := strings.Split(tags, " ")
	item := baza.Item_list[i] //ispravit

	item.Name = name
	item.Content = &content
	item.Tag = tag
	baza.Item_list[i] = item
	stringy := listItemToString(item)
	baza.Update(stringy, "templates/data.txt")
	s.Execute(w, nil)
}

func deleteHandlerFunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	index := r.FormValue("Index")
	i, err := strconv.Atoi(index)
	fmt.Println(index)
	if err != nil {
		panic(err)
	}
	baza.Delete(i)
	s.Execute(w, nil)
}
func editHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query().Get("Index"))
	index := r.URL.Query().Get("Index")
	fmt.Println("INDEKS JE: ", index)
	fmt.Println("editHanderFunc")
	i, err := strconv.Atoi(index)
	if err != nil {
		panic(err)
	}
	item := baza.Item_list[i]
	fmt.Println(item.Name)
	e.Execute(w, baza.Item_list[i])
}
func createHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		c.Execute(w, nil)
	}
	r.ParseForm()
	name := r.FormValue("name")
	content := r.FormValue("content")
	tags := r.FormValue("tags")
	tag := strings.Split(tags, " ")
	if name == "" || content == "" || tag[0] == "" {
		s.Execute(w, nil)
	} else {
		item := logic.Item1{Name: name, Content: &content, Index: len(baza.Item_list), Tag: tag}
		baza.Add(item)
	}
	s.Execute(w, nil)
}
func readHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		s.Execute(w, nil)
		return
	}
	r.ParseForm()
	tag := r.FormValue("tag")
	o := baza.ReadTag(tag)
	s.Execute(w, o)
}

/*func editHandlerFunc(w http.ResponseWriter, r *http.Request) {
	return
}
func deleteHandlerFunc(w http.ResponseWriter, r *http.Request) {
	return
}
func createHandlerFunc(w http.ResponseWriter, r *http.Request) {
	return
}
*/

func toFile(stringToWrite string, path string) {

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(stringToWrite); err != nil {
		panic(err)
	}

}
func fileToList(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	baza.Item_list = nil
	baza.Baza = make(map[string]logic.Map_value1)
	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	a := logic.Item1{}
	for {

		a = logic.Item1{}
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		a.Name = rec[0]
		a.Content = &rec[1]
		a.Index, err = strconv.Atoi(rec[2])

		for i := 3; i < len(rec); i++ {
			a.Tag = append(a.Tag, rec[i])

		}

		//logic.Item_list = append(logic.Item_list, a)
		baza.Add(a)
	}

}

func ListToString() string {
	x := ""
	for num := range baza.Item_list {
		x += baza.Item_list[num].Name
		x += ","
		x += *(baza.Item_list[num].Content)
		x += ","
		x += strconv.Itoa(baza.Item_list[num].Index)
		x += ","
		for num1, val1 := range baza.Item_list[num].Tag {
			x += baza.Item_list[num].Tag[num1]
			if num1 != len(val1) {
				x += ","
			}
		}
		x += "\n"

	}
	return x
}
func listItemToString(parsedItem logic.Item1) string {
	x := ""
	x += parsedItem.Name
	x += ","
	x += *parsedItem.Content
	x += ","
	x += strconv.Itoa(parsedItem.Index)
	x += ","
	for num, val := range parsedItem.Tag {
		x += parsedItem.Tag[num]
		if num != len(val)-1 {
			x += ","
		}
	}
	return x
}
