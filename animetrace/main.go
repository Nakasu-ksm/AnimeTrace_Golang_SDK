package animetrace

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
)

var all_model_map = map[string]int{
	"anime_model_lovelive": 0,
	"pre_stable":           0,
	"anime":                0,
	"game":                 0,
	"game_model_kirakira":  0,
}

func (wk *WorkerType) Recognition() {
	client := http.Client{}
	apiUrl := "https://aiapiv2.animedb.cn/ai/api/detect"

	req, err := http.NewRequest("POST", apiUrl, wk.buffer)
	req.Header.Set("Content-Type", wk.writer.FormDataContentType())
	if err != nil {
		panic("画像の読み込みに失敗しました！")
	}

	content, err := client.Do(req)
	defer content.Body.Close()
	all, err := io.ReadAll(content.Body)
	if err != nil {
		panic("画像の読み込みに失敗しました！")
	}
	fmt.Println("識別終了！")
	//fmt.Println(string(all))
	//fmt.Println(all)
	//fmt.Println(string(all))
	wk.result = &all
}
func (wk *WorkerType) SetMultiple(bool2 bool) {
	if wk.lock {
		panic("画像アップロード後の設定変更はできません。")
	}
	//if id != 0 {
	//	panic("自分でロジックを実装してください")
	//}
	if bool2 {
		wk.p.Is_multi = 1
	}
}

func (wk *WorkerType) SetModel(model string) {
	//fmt.Println(wk.p)
	if wk.lock {
		panic("画像アップロード後の設定変更はできません。")
	}
	if _, ok := all_model_map[model]; !ok {
		panic("認識モデルは存在しない。参考資料 https://docs.animedb.cn/#/introduction を参照。")
	}
	wk.p.Model = model

}
func (wk *WorkerType) SetAI(bool2 bool) {
	if wk.lock {
		panic("画像アップロード後の設定変更はできません。")
	}
	if bool2 {
		wk.p.ai = 1
	}
}
func (wk *WorkerType) SetImage(imageBytes []byte) {
	if wk.lock {
		panic("画像アップロード後の設定変更はできません。")
	}
	wk.lock = true
	defer wk.writer.Close()
	_ = wk.writer.WriteField("is_multi", strconv.Itoa(wk.p.Is_multi))
	_ = wk.writer.WriteField("model", wk.p.Model)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "1.jpg"))
	h.Set("Content-Type", "image/png")
	file_parameter, _ := wk.writer.CreatePart(h)
	file_parameter.Write(imageBytes)
	wk.buffer.Write([]byte("\r\n" + "\r\n" + "--" + wk.writer.Boundary() + "--\r\n"))
}

func API() Worker {
	var worker_return Worker
	worker := WorkerType{}
	worker.p = &Params{}
	worker.lock = false
	worker.buffer = new(bytes.Buffer)
	worker.writer = multipart.NewWriter(worker.buffer)
	worker_return = &worker
	return worker_return
}

type Worker interface {
	SetImage(imageBytes []byte)
	SetModel(model string)
	SetAI(bool2 bool)
	ConvertToJson() (error, Response)
	IsReturnMulti() bool
	Recognition()
	SetMultiple(bool2 bool)
	GetResultString() string
}

type Params struct {
	Is_multi int
	Model    string
	ai       int
}

type Response struct {
	Code    int              `json:"code"`
	Data    []AnimeCharacter `json:"data"`
	Ai      bool             `json:"ai"`
	NewCode int              `json:"new_code"`
}

type AnimeCharacter struct {
	Char    []MultipleCharacter `json:"char,omitempty"`
	Box     [5]float64          `json:"box"`
	Name    string              `json:"name,omitempty"`
	Cartoon string              `json:"cartoonname,omitempty"`
	Acc     float64             `json:"acc_percent"`
	BoxId   string              `json:"box_id"`
}

type MultipleCharacter struct {
	Name    string  `json:"name"`
	Cartoon string  `json:"cartoonname"`
	Acc     float64 `json:"acc"`
}

type ResultBytes []byte

func (wk *WorkerType) ConvertToJson() (error, Response) {
	if wk.result == nil {
		panic("このメソッドは、画像を認識した後にのみ呼び出すことができます。")
	}
	var err error
	var resp Response
	err = json.Unmarshal(*wk.result, &resp)
	if err != nil {
		panic("パースエラー")
	}
	if resp.Code != 0 && resp.Code != 17720 {
		return errors.New("error"), resp

	}
	return nil, resp
}

type WorkerType struct {
	p      *Params
	writer *multipart.Writer
	buffer *bytes.Buffer
	result *[]byte
	lock   bool
	res    *Response
}

func (wk *WorkerType) IsReturnMulti() bool {
	if wk.p.Is_multi == 1 {
		return true
	}
	return false
}

func (wk *WorkerType) GetResultString() string {
	return string(*wk.result)
}
