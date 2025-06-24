package dbx

import "fmt"

func CloseAllDBXTenant() {
	cacheDBXTenant.Range(func(key, value interface{}) bool {
		if value == nil {
			return false
		}
		fmt.Printf("Closing DBX connection:% s", key)
		if value.(DBXTenant).DB != nil {
			db := value.(DBX).DB
			if db.Stats().InUse > 0 {
				err := value.(DBXTenant).DB.Close()
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			return true
		}
		return false
	})
}
func CloseAllDBX() {

	dbxCache.Range(func(key, value interface{}) bool {
		fmt.Printf("Closing DBX connection:% s", key)
		if value == nil {
			return false
		}
		if value.(DBX).DB != nil {
			db := value.(DBX).DB
			if db.Stats().InUse > 0 {
				err := value.(DBX).DB.Close()
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			return true
		}
		return false
	})
}
func CloseAll() {
	fmt.Println("Closing all DBX connections")
	CloseAllDBX()
	CloseAllDBXTenant()
}
