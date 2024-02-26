package main

import (
	"fmt"
	"go_sdk/animetrace"
	"os"
)

func main() {
	worker := animetrace.API()

	worker.SetModel("anime_model_lovelive")
	worker.SetMultiple(1)

	imageBytes, err := os.ReadFile("demo.png")
	if err != nil {
		panic("画像の読み込みに失敗しました！")
	}
	worker.SetImage(imageBytes)
	worker.Recognition()
	jsonReturn := worker.ConvertToJson()
	//fmt.Println(jsonReturn)
	if worker.IsReturnMulti() {
		fmt.Println(jsonReturn.Data[0].Char[0].Name)
	} else {
		fmt.Println(jsonReturn.Data[0].Name)
	}

}
