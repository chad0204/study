package effective

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

type Movie struct {
	Title  string
	Year   int  `json:"released"`        // 转json Year -> released
	Color  bool `json:"color,omitempty"` // 转json Color -> color, omitempty就是忽略零值, 不显示该字段
	Actors []string
}

var movies = []Movie{
	{Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
}

func TestMarshal(t *testing.T) {
	data, err := json.Marshal(movies)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	fmt.Printf("%s\n", data)

	data1, err1 := json.MarshalIndent(movies, "", "    ") //添加缩进
	if err1 != nil {
		log.Fatalf("JSON marshaling failed: %s", err1)
	}
	fmt.Printf("%s\n", data1)
}

func TestUnMarshal(t *testing.T) {
	data, err := json.Marshal(movies)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	var titles []struct{ Title string }
	err = json.Unmarshal(data, &titles)
	if err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}

	fmt.Println(titles)

}
