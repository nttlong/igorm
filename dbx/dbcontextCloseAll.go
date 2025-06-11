package dbx

import "fmt"

func CloseAllDBXTenant() {
	cacheDBXTenant.Range(func(key, value interface{}) bool {
		fmt.Printf("Closing DBX connection:% s", key)
		value.(DBXTenant).Close()
		return true
	})
}
func CloseAllDBX() {
	dbxCache.Range(func(key, value interface{}) bool {
		fmt.Printf("Closing DBX connection:% s", key)
		value.(DBX).Close()
		return true
	})
}
func CloseAll() {
	fmt.Println("Closing all DBX connections")
	CloseAllDBX()
	CloseAllDBXTenant()
}
