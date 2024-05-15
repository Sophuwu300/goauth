package main

import "fmt"
import "sophuwu.site/goauth"

func main() {
	u, e := goauth.NewUser("sophie", "password")
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(u)
}
