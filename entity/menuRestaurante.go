package entity

type MenuRestaurante []ItemMenu

type ItemMenu struct {
	Nome      string  `json:"nome"`
	Preco     float64 `json:"preco"`
	Tipo      string  `json:"tipo"`
	Descricao string  `json:"descricao"`
	Imagem    []byte  `json:"img"`
}

type ItemMenuBson struct {
	Nome      string  `bson:"nome"`
	Preco     float64 `bson:"preco"`
	Tipo      string  `bson:"tipo"`
	Descricao string  `bson:"descricao"`
	Imagem    []byte  `bson:"img"`
}

type Mesa struct {
	Nome          string `bson:"nome"`
	NomeCliente   string `bson:"nomeCliente"`
	NumeroCliente string `bson:"numeroCliente"`
	Reserva       bool   `bson:"reserva"`
	NumeroPessoas int    `bson:"numeroPessoas"`
	DataReserva   string `bson:"dataReserva"`
}
