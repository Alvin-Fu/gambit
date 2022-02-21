package main

import (
	"fmt"

	"github.com/maaslalani/gambit/border"
)

func main() {
	tests := []struct {
		borderFunc func() string
		want       string
	}{

		{
			border.Top,
			"   ┌───┬───┬───┬───┬───┬───┬───┬───┐\n",
		},
		{
			border.Middle,
			"   ├───┼───┼───┼───┼───┼───┼───┼───┤\n",
		},
		{
			border.Bottom,
			"   └───┴───┴───┴───┴───┴───┴───┴───┘\n",
		},
	}

	for _, test := range tests {
		got := test.borderFunc()
		if got != test.want {
			fmt.Printf("want %s, got %s\n", test.want, got)
		}
		fmt.Println(got)
	}
	bo := border.Build(tests[0].want, tests[1].want, tests[2].want)
	fmt.Println(bo)
	fmt.Println(border.BottomLabels(false))
}
