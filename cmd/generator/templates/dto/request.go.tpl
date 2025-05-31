package request

type Create{{.PascalName}} struct {
	// contoh field, sesuaikan dengan kebutuhan
}

type Update{{.PascalName}} struct {
	// contoh field, sesuaikan dengan kebutuhan
}

type UpdatePassOrVerify struct {
	// contoh field, sesuaikan dengan kebutuhan
}

type Query{{.PascalName}} struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}
