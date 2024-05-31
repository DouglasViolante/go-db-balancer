package utils

import (
	"go-db-balancer/internal/pkg/domain"
)

func FindMaxPos(tables []domain.Table) int {
	maxPos := 0
	for i := 1; i < len(tables); i++ {
		if tables[i].ItemCount > tables[maxPos].ItemCount {
			maxPos = i
		}
	}
	return maxPos
}