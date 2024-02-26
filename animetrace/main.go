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

func Recognition(buffer *bytes.Buffer, boundary string) ResultBytes {
	client := http.Client{}
	apiUrl := "https://aiapiv2.animedb.cn/ai/api/detect"

	req, err := http.NewRequest("POST", apiUrl, buffer)
	req.Header.Set("Content-Type", boundary)
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

	return all
}
func (p *Params) SetMultiple(id int) {
	if id != 0 {
		panic("自分でロジックを実装してください")
	}
	p.Is_multi = id
}

func (p *Params) SetModel(model string) {
	if _, ok := all_model_map[model]; !ok {
		panic("認識モデルは存在しない。参考資料 https://docs.animedb.cn/#/introduction を参照。")
	}
	p.Model = model

}
func (p Params) SetConfig(buffer *bytes.Buffer, imageBytes []byte) *multipart.Writer {

	writer := multipart.NewWriter(buffer)
	defer writer.Close()
	_ = writer.WriteField("is_multi", strconv.Itoa(p.Is_multi))
	_ = writer.WriteField("model", p.Model)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "1.jpg"))
	h.Set("Content-Type", "image/png")
	file_parameter, _ := writer.CreatePart(h)
	file_parameter.Write(imageBytes)
	//buffer.Write([]byte("\r\n" + "\r\n" + "--" + get_boundary + "--\r\n"))

	return writer
}

func API() *Params {
	return &Params{}
}

type Params struct {
	Is_multi int
	Model    string
}

type Worker interface {
	SetConfig(buffer *bytes.Buffer, imageBytes []byte) *multipart.Writer
	SetMultiple(id int)
	SetModel(model string)
}

type Response struct {
	Code    int              `json:"code"`
	Data    []AnimeCharacter `json:"data"`
	Ai      bool             `json:"ai"`
	NewCode int              `json:"new_code"`
}

type AnimeCharacter struct {
	Box     [5]float64 `json:"box"`
	Name    string     `json:"name"`
	Cartoon string     `json:"cartoonname"`
	Acc     float64    `json:"acc_percent"`
	BoxId   string     `json:"box_id"`
}

type ResultBytes []byte

func (json_string ResultBytes) ConvertToJson() Response {
	var resp Response
	err := json.Unmarshal(json_string, &resp)
	if err != nil {
		panic("パースエラー")
	}
	return resp
}
