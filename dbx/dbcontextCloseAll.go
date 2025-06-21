package dbx

import "fmt"

func CloseAllDBXTenant() {
	cacheDBXTenant.Range(func(key, value interface{}) bool {
		fmt.Printf("Closing DBX connection:% s", key)

		err := value.(DBXTenant).DB.Close()
		if err != nil {
			panic(err)
		}
		return true
	})
}
func CloseAllDBX() {
	dbxCache.Range(func(key, value interface{}) bool {
		fmt.Printf("Closing DBX connection:% s", key)
		err := value.(DBX).DB.Close()
		if err != nil {
			panic(err)
		}
		return true
	})
}
func CloseAll() {
	fmt.Println("Closing all DBX connections")
	CloseAllDBX()
	CloseAllDBXTenant()
}
