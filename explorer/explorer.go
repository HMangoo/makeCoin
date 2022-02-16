package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/HMangoo/makeCoin/blockchain"
)

const (
	port        string = ":4000"
	templateDir string = "explorer/templates/"
)

// templates : Template struct의 pointer
var templates *template.Template

// PageTitle : 만들어질 page의 title 
// Blocks : blockchain의 Block slice pointer
type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

// Handler
func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}
// ExecuteTemplate : 이름이 home인 template를 실행
// -> response로 template를 실행시킴 (ResponseWriter를 쓰고, data전달)
 
// Handler
func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // 아무 데이테 없이 template 씀
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm() // r.Form을 생성
		data := r.Form.Get("blockData") // Form에서 data를 가져옴 (같은 name을 써야함)
		blockchain.GetBlockchain().AddBlock(data) 
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start() {
	// load template / pattern(*)을 이용하여 template를 load, template.Must(helper function) : error check
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml")) // use standard libaray template
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml")) // use templates variable

	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
	
}