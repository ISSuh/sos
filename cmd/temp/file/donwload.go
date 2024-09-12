package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/download", downloadHandler)
	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// 파일을 엽니다.
	file, err := os.Open("/Users/issuh/workspace/git/issuh/sos/cmd/temp/file/sample.mp3")
	if err != nil {
		http.Error(w, "Unable to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 파일 정보를 가져옵니다.
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Unable to get file info", http.StatusInternalServerError)
		return
	}

	// 응답 헤더를 설정합니다.
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 파일을 반복문으로 읽고 클라이언트에게 전송합니다.
	buf := make([]byte, 32*1024) // 32KB 버퍼
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		if n == 0 {
			break
		}

		// buf의 내용을 클라이언트에게 전송합니다.
		if _, err := w.Write(buf); err != nil {
			http.Error(w, "Error writing to response", http.StatusInternalServerError)
			return
		}
	}
}
