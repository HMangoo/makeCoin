package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/HMangoo/makeCoin/blockchain"
	"github.com/HMangoo/makeCoin/utils"
	"github.com/gorilla/mux"
)



var port string



// URLDescription이 어떻게 출력되는지 결정할 수 있음
type url string
// Marshal : interface를 json으로 encoding한 것을 return
// MarshalText : json string으로써 어떻게 보일지 결정하는 method
func (u url) MarshalText() ([]byte, error) { // 정확한 이름과 시그니처를 이용하여 interface를 구현할 수 있음
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL url	`json:"url"`	// JSON에서 보여지는 방식, JSON에서는 소문자를 사용
	Method string	`json:"method"`
	Description string `json:"description"`
	Payload string	`json:"payload,omitempty"`
}

// blocks : 1. create struct
type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}


func documentation (rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL: url("/"), // homepage
			Method: "GET",	// GET : GET을 통해 해당 리소스를 조회. 리소스를 조회하고 해당 document에 대한 자세한 정보를 가져온다 
			Description: "See Documentation",
		},
		{
			URL: url("/blocks"),
			Method: "GET",
			Description: "See All Blocks",
		},
		{
			URL: url("/blocks"),
			Method: "POST",
			Description: "Add A Block",
			Payload: "data:string",
		},
		{
			URL: url("/blocks/{height}"),
			Method: "GET",
			Description: "See A Block",
		},
	}
	// Response Header의 Content-Type을 수정 -> 브라우저에게 JSON이라는 것을 알림
	
	// struct를 JSON으로 바꿔서 원하는 포맷으로 보여줄 수 있음 // NewEncoder returns a new encoder that writes to w
	json.NewEncoder(rw).Encode(data) 
}
// Marshal : 메모리 형식으로 저장된 객체를, 저장/송신 할 수 있도록 변환, Go interface into JSON   

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json") 
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks()) // Endcode가 Marshal의 일을 해주고, 결과를 ResponseWriter에 작성
	case "POST":
		// {"data": "my block data"}
		// POST request에서 유저가 보내는 위와같은 블록 정보를 받아와서 우리가 작업할 수 있도록 golang의 struct로 변환
		
		/*
			1. create struct
			2. decode : request body -> struct
		*/
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		// body request를 decode하고 결과물을 addBlockBody에 넣음
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated) // HTTP response와 함께 status code를 보냄
	}

}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["height"])
	utils.HandleErr(err)
	block, err := blockchain.GetBlockchain().GetBlock(id)
	endcoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		endcoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		endcoder.Encode(block)
	}
	
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	// JSON이 middleware로 시작되면, 이 함수가 가장 먼저 호출
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter() // https://github.com/gorilla/mux
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// .Methods("GET") : GET만 사용 (다른 method로 부터 막아줌)