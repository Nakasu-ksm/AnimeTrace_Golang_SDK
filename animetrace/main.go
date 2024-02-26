package animetrace

import (
	"bytes"
	"encoding/json"
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
func (wk *WorkerType) SetMultiple(id int) {
	//if id != 0 {
	//	panic("自分でロジックを実装してください")
	//}
	wk.p.Is_multi = id
}

func (wk *WorkerType) SetModel(model string) {
	if _, ok := all_model_map[model]; !ok {
		panic("認識モデルは存在しない。参考資料 https://docs.animedb.cn/#/introduction を参照。")
	}
	wk.p.Model = model

}
func (wk *WorkerType) SetImage(imageBytes []byte) {
	defer wk.writer.Close()
	_ = wk.writer.WriteField("is_multi", strconv.Itoa(wk.p.Is_multi))
	_ = wk.writer.WriteField("model", wk.p.Model)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "1.jpg"))
	h.Set("Content-Type", "image/png")
	file_parameter, _ := wk.writer.CreatePart(h)
	file_parameter.Write(imageBytes)
	//wk.buffer.Write([]byte("\r\n" + "\r\n" + "--" + get_boundary + "--\r\n"))
}

func API() *WorkerType {
	worker := WorkerType{}
	worker.p = &Params{}
	worker.buffer = new(bytes.Buffer)
	worker.writer = multipart.NewWriter(worker.buffer)
	return &worker
}

type Params struct {
	Is_multi int
	Model    string
}

//
//type Worker interface {
//	SetImage(buffer *bytes.Buffer, imageBytes []byte) *multipart.Writer
//	SetMultiple(id int)
//	SetModel(model string)
//}

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

func (wk *WorkerType) ConvertToJson() Response {
	//fmt.Println(*wk.p)
	var err error
	var resp Response
	err = json.Unmarshal(*wk.result, &resp)
	if err != nil {
		panic("パースエラー")
	}
	return resp
}

type WorkerType struct {
	p      *Params
	writer *multipart.Writer
	buffer *bytes.Buffer
	result *[]byte
}

func (wk *WorkerType) IsReturnMulti() bool {
	if wk.p.Is_multi == 1 {
		return true
	}
	return false
}
