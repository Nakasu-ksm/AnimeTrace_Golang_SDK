package main

import (
	"bytes"
	"fmt"
	"go_sdk/animetrace"
	"os"
)

func main() {
	var worker animetrace.Worker
	worker = animetrace.API()
	worker.SetMultiple(0)
	worker.SetModel("anime_model_lovelive")
	imageBytes, err := os.ReadFile("demo.png")
	if err != nil {
		panic("画像の読み込みに失敗しました！")
	}
	buffer := new(bytes.Buffer)
	writer := worker.SetConfig(buffer, imageBytes)
	result := animetrace.Recognition(buffer, writer.FormDataContentType())
	response := result.ConvertToJson()
	fmt.Println(response.Data[0].Name)
}
