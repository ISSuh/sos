package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("서버가 8080 포트에서 실행 중입니다...")
	http.ListenAndServe(":22511", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "허용되지 않는 메서드입니다.", http.StatusMethodNotAllowed)
		return
	}

	// 요청 본문을 파싱합니다.
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "요청을 파싱할 수 없습니다.", http.StatusBadRequest)
		return
	}

	// 파일을 추출합니다.
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "파일을 추출할 수 없습니다.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 파일을 저장합니다.
	dst, err := os.Create(handler.Filename)
	if err != nil {
		http.Error(w, "파일을 저장할 수 없습니다.", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "파일을 저장하는 동안 오류가 발생했습니다.", http.StatusInternalServerError)
		return
	}

	// 텍스트 데이터를 추출합니다.
	text := r.FormValue("text")

	// 응답을 작성합니다.
	fmt.Fprintf(w, "파일 업로드 성공: %s\n", handler.Filename)
	fmt.Fprintf(w, "텍스트 데이터: %s\n", text)
}
