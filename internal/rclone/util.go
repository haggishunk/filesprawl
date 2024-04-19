package rclone

import "fmt"

// dumpListResponse is a toy to print a ListResponse to stdout
func dumpListResponse(lr ListResponse) error {
	for _, item := range lr.List {
		fmt.Printf("Item Id %s name %s\n%s %d %t\n", item.ID, item.Name, item.Path, item.Size, item.IsDir)
	}
	return nil
}
