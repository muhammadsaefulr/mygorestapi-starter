package request

type UserPoints struct {
	UserId     string `json:"-"`
	TypeUpdate string `json:"type_update" example:"add" enums:"add,subtract" format:"string" description:"Jenis update point: 'add' untuk menambah, 'subtract' untuk mengurangi"`
	Value      int    `json:"value_added" example:"1" minimum:"1" description:"Jumlah point yang ingin ditambah atau dikurangi"`
}

type QueryUserPoints struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}
