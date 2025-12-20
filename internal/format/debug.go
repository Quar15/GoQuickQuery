package format

import "fmt"

func PrintMap(m []map[string]any) {
	for i, data := range m {
		for k, v := range data {
			fmt.Printf("idx = %v | key = %v | val =  %v | type=%T\n", i, k, v, v)
		}
	}
}
