package main

import (
	"encoding/json"
	"fmt"

	"github.com/starius/api2/typegen"
	"github.com/starius/api2/typegen/example/types"
)

func main() {
	s := typegen.New()
	s.Add(types.T{})
	s.Add(types.User{})

	str := s.RenderToSwagger()
	res, _ := json.MarshalIndent(str, "", " ")
	fmt.Print(string(res))
	// if err != nil {
	// 	panic(err)
	// }
}
