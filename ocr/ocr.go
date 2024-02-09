package ocr

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"gocv.io/x/gocv"
)

func OCR(cell gocv.Mat) int {
	return callAPI(cell)
}

// APIを叩く関数
func callAPI(cell gocv.Mat) int {
	// 画像をbase64にエンコード
	bbuf, err := gocv.IMEncode(".png", cell)
	if err != nil {
		panic(err)
	}
	b64 := base64.StdEncoding.EncodeToString(bbuf.GetBytes())

	uri := "http://127.0.0.1:8888/ocr"
	jsonStr := fmt.Sprintf(`{"b64": "%s"}`, b64)
	reqBody := bytes.NewBuffer([]byte(jsonStr))

	resp, err := http.Post(uri, "application/json", reqBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := bytes.Buffer{}
	buf.ReadFrom(resp.Body)
	respBody := buf.String()
	pred, err := strconv.Atoi(respBody)
	if err != nil {
		panic(err)
	}
	return pred
}
