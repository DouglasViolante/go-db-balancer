package domain

type Specs struct {
	ID    int     `json:"ID" dynamodbav:"ID"`
	Power float32 `json:"Power" dynamodbav:"Power"`
	Cost  float32 `json:"Cost" dynamodbav:"Cost"`
}