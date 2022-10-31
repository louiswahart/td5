package vtypes

type Request struct {
	Choice []Alternative `json:"alts"`
}

type Response struct {
	Result int `json:"res"`
}

type Alternative int

type Profile [][]Alternative
