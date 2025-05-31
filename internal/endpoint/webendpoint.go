package webendpoint

import (
	"SearchServices/internal/prediction"
	"fmt"
	"io"
	"net/http"
)

//const form = `<html>
//    <head>
//    <title></title>
//    </head>
//    <body>
//        <form action="/main" method="post">
//            <label>Название услуги</label><input type="text" name="answer">
//            <input type="submit" value="watch">
//        </form>
//    </body>
//</html>`

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Counter равен")
	} else {
		fmt.Fprintln(w, "Поддерживается только метод GET")
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешаем все источники
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodPost {
		read := r.FormValue("answer")
		io.WriteString(w, prediction.FinalResponce(read))
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
}
