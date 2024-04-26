package entity

type MenuRestaurante []ItemMenu

type ItemMenu struct {
	Nome      string  `json:"nome"`
	Preco     float64 `json:"preco"`
	Tipo      string  `json:"tipo"`
	Descricao string  `json:"descricao"`
	Imagem    []byte  `json:"img"`
}
