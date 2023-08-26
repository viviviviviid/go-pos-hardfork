package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/viviviviviid/go-coin/blockchain"
)

const (
	port        string = ":4000"
	templateDir string = "templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) { // (유저에게 보내고싶은 내용, pointer)
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData") // add.gohtml에서 버튼을 눌렀을때 얻어온 데이터
		// .Form은 함수가 아닌 map
		blockchain.GetBlockchain().AddBlock(data)
		fmt.Println(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect) // redirect
	}

}

func main() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml")) // pattern을 인자로 받음
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil)) // error 발생 시 log처리 log.Fatal()
}
