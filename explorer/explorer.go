package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/viviviviviid/go-coin/blockchain"
)

var templates *template.Template

const (
	templateDir string = "explorer/templates/"
)

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) { // (유저에게 보내고싶은 내용, pointer)
	data := homeData{"Home", nil}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		blockchain.Blockchain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect) // redirect
		// statusPer- 는 Redirect시 300번대 이상의 status를 파라미터로 넣어야하는데
		// 뭘 넣을지 고민하는 사람들을 대신해서, 자동으로 대신 넣어주는 메서드함수
	}

}

func Start(port int) {
	handler := http.NewServeMux()
	// Must를 찍어보면, 그냥 단순하게 에러가 있는지 확인해주는 helper
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))     // template들을 load
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml")) // template들을 load한 내용을 가지고 또다시 load
	// **/*.gohtml을 하지못해서 위 두줄이 된 것
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler)) // error 발생 시 log처리 log.Fatal()
	// 우리가 만든 Mux로 defaultMux를 대체함
}
