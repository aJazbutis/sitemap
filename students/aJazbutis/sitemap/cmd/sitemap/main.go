package main

import (
	"flag"
	"fmt"
	"sitemap"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	urlFlag := flag.String("s", "", "website name, will return if not set")
	depth := flag.Int("d", 0, "how deep to map")
	fileFlag := flag.Bool("f", false, "output to a file")
	flag.Parse()
	if *depth < 0 || *urlFlag == "" {
		fmt.Println("Usage:")
		flag.VisitAll(func(f *flag.Flag){
			fmt.Printf("    -%s: %s (default %v)\n", f.Name, f.Usage, f.DefValue)
		})
		return
	}
	sitemap.MapSite(*urlFlag, *depth, *fileFlag)
}
