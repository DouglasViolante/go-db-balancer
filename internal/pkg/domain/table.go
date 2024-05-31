package domain

type Table struct {
	TableName  string  `json:"tableName"`
	ItemCount  int     `json:"itemCount"`
	Difference float64 `json:"difference"`
}