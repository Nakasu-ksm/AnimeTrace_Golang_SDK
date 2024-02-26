package main

import (
	"fmt"
	"go_sdk/animetrace"
	"os"
	"time"
)

type AIRecognitionError struct {
	errorTitle string
	timeNow    time.Time
}

func (e *AIRecognitionError) Error() string {
	e.timeNow = time.Now()
	return "エラー発生: " + e.errorTitle + "  At" + e.timeNow.String()
}

func catchError(err string) error {
	return &AIRecognitionError{errorTitle: err}
}

func main() {

	defer func() {
		if rec := recover(); rec != nil {
			strs, ok := rec.(string)
			if ok {
				fmt.Println(catchError(strs))
			}

		}
	}()
	worker := animetrace.API()

	worker.SetMultiple(true)
	worker.SetModel("anime_model_lovelive")
	worker.SetAI(true)
	imageBytes, err := os.ReadFile("demo.png")
	if err != nil {
		panic("画像の読み込みに失敗しました！")
	}
	worker.SetImage(imageBytes)
	worker.Recognition()
	err, jsonReturn := worker.ConvertToJson()
	if err != nil {
		panic("画像認識異常")
	}
	if worker.IsReturnMulti() {
		fmt.Println(jsonReturn.Data[0].Char[0].Name)
	} else {
		fmt.Println(jsonReturn.Data[0].Name)
	}

}
