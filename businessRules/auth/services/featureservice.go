package services

import (
	authModels "dbmodels/auth"
	"dbx"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FeatureDataDetail struct {
	Module string `json:"module"`
	Action string `json:"action"`
}
type FeatureData struct {
	FeatureId  string                     `json:"featureId"`
	RefModules map[string]map[string]bool `json:"refModules"`
}
type FeatureService struct {
	CacheService
	TenantDb  *dbx.DBXTenant
	FeatureId string
	Module    string
	Action    string
	IsDebug   bool
}

func (f *FeatureService) GetUser(userId string) (user *authModels.User, err error) {
	ret, err := dbx.Query[authModels.User](f.TenantDb, f.Context).Where("UserId = ?", userId).First()
	if err != nil {
		return nil, err
	}
	return ret, nil
}
func (f *FeatureService) RegisterFeature() error {
	// crete folder store feature name featureConfig
	folderPath := "./data/features/" // thư mục muốn tạo

	err := os.MkdirAll(folderPath, os.ModePerm) // ModePerm = 0777
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}
	absPath, err := filepath.Abs(folderPath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return err
	}
	fmt.Println("Directory created at:", absPath)
	featureFilePath := filepath.Join(folderPath, f.FeatureId+".json")
	// check if feature file exist
	jsonData := &FeatureData{
		FeatureId:  f.FeatureId,
		RefModules: make(map[string]map[string]bool),
	}
	jsonData.RefModules[f.Module] = make(map[string]bool)
	jsonData.RefModules[f.Module][f.Action] = true
	if _, err := os.Stat(featureFilePath); os.IsNotExist(err) {
		jsonText, err := json.MarshalIndent(jsonData, "", "  ")

		if err != nil {
			return err
		}
		err = os.WriteFile(featureFilePath, jsonText, 0644)
		if err != nil {
			return err
		}
	} else {
		// update feature file
		fileData, err := os.ReadFile(featureFilePath)
		if err != nil {
			return err
		}
		jsonData := &FeatureData{}
		err = json.Unmarshal(fileData, jsonData)
		if err != nil {
			return err
		}
		if _, ok := jsonData.RefModules[f.Module]; !ok {
			jsonData.RefModules[f.Module] = make(map[string]bool)
		}
		if _, ok := jsonData.RefModules[f.Module][f.Action]; !ok {
			jsonData.RefModules[f.Module][f.Action] = true
		}
		jsonText, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return err
		}
		err = os.WriteFile(featureFilePath, jsonText, 0644)
		if err != nil {
			return err
		}

	}

	featureData := authModels.Features{
		Id:        dbx.NewUUID(),
		Name:      f.FeatureId,
		CreatedAt: time.Now(),
		CreatedBy: "root",
	}
	err = dbx.InsertWithContext(f.Context, f.TenantDb, &featureData)

	if err != nil {
		if dbxErr, ok := err.(*dbx.DBXError); ok {
			if dbxErr.Code == dbx.DBXErrorCodeDuplicate {

				return nil
			}
		}
		return err
	}
	return nil
}
