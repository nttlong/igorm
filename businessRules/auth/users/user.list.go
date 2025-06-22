package auth

import (
	authModels "dbmodels/auth"
	"dbx"
	"errors"
)

func (u *User) List(filter struct {
	PageIndex int
	PageSize  int
	Sort      string
}) (interface{}, error) {

	tokenInfo, err := u.ValidateAccessToken(u.AccessToken)
	if err != nil {
		return nil, err
	}
	user, err := u.GetUser(tokenInfo.UserId)
	if err != nil {
		return nil, err
	}
	if user.IsSupperUser {
		if u.FeatureService.FeatureId == "" {
			return nil, errors.New("Feature is empty")
		}
		if u.FeatureService.IsDebug {
			u.FeatureService.RegisterFeature()
		}

	}

	if err != nil || tokenInfo == nil {

		return nil, err
	}
	qr := dbx.Pager[authModels.User](u.TenantDb, u.Context)

	qr.Select("UserId,Username,Email,CreatedAt,UpdatedAt,CreatedBy,UpdatedBy")
	if filter.PageSize == 0 {
		filter.PageSize = 50
	}
	if filter.Sort == "" {
		filter.Sort = "Id asc"
	}
	qr.Page(filter.PageIndex).Size(filter.PageSize).Sort(filter.Sort)
	// var wg sync.WaitGroup
	resultChan := make(chan dbx.QueryResult[authModels.User], 1) // Generic theo kiểu User
	defer close(resultChan)
	// ret, err := qr.Query()

	// return ret, err
	//wg.Add(1)
	go qr.QueryAsync(resultChan)

	//wg.Wait()

	result := <-resultChan
	return result.Data, result.Err

}
