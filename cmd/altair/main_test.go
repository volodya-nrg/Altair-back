package main

import (
	"altair/api/request"
	"altair/pkg/helpers"
	"altair/pkg/service"
	"altair/server"
	"altair/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Info:
// https://golang.hotexamples.com/ru/examples/mime.multipart/-/NewWriter/golang-newwriter-function-examples.html
// https://github.com/gin-gonic/gin/blob/66d2c30c54ff8042f5ae13d9ebb26dfe556561fe/binding/binding_test.go#L530

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseUser, 0)
	testCases = append(testCases, &testCaseUser{
		Want: 200,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}
func TestGetUsersUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseUser, 0)
	serviceUsers := service.NewUserService()
	pUsers, err := serviceUsers.GetUsers()
	assert.NoError(t, err)

	// сделаем реверсивный список
	pReversedUsers := make([]*storage.User, 0)
	for i := range pUsers {
		pReversedUsers = append(pReversedUsers, pUsers[len(pUsers)-1-i])
	}

	// создаим нужные нам автоматичекие testCase-ы
	for i, user := range pReversedUsers {
		if i == 3 {
			break
		}

		a := &testCaseUser{MyStr: fmt.Sprint(user.UserId), Want: 200}
		b := &testCaseUser{MyStr: fmt.Sprint(user.UserId + 3), Want: 404}

		testCases = append(testCases, a, b)
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/api/v1/users/"+tc.MyStr, nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}
func TestPostUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseUser, 0)

	testCases = append(testCases, &testCaseUser{
		RequestPost: &request.PostUser{
			Email:           "test@" + helpers.RandStringRunes(5) + "." + helpers.RandStringRunes(3),
			Password:        "123456",
			PasswordConfirm: "123456",
			AgreeOffer:      true,
			AgreePolicy:     true,
		},
		Want: 201})
	testCases = append(testCases, &testCaseUser{
		RequestPost: &request.PostUser{
			Email:           "test@" + helpers.RandStringRunes(5) + "." + helpers.RandStringRunes(3),
			Password:        "123456",
			PasswordConfirm: "123456",
			AgreeOffer:      false,
			AgreePolicy:     true,
		},
		Want: 400})
	testCases = append(testCases, &testCaseUser{
		RequestPost: &request.PostUser{
			Email:           "test@" + helpers.RandStringRunes(5) + "." + helpers.RandStringRunes(3),
			Password:        "12345",
			PasswordConfirm: "123456",
			AgreeOffer:      true,
			AgreePolicy:     true,
		},
		Want: 400})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.RequestPost)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(b))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if !assert.Equal(t, tc.Want, w.Code) {
				t.Log("-----> Response: ", w.Body)
			}

			if w.Code == 201 {
				user := new(storage.User)
				assert.NoError(t, json.Unmarshal(w.Body.Bytes(), user))
				assert.NoError(t, server.Db.Debug().Delete(user).Error)
			}
		})
	}
}
func TestPutUsersUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceUsers := service.NewUserService()
	user := &storage.User{
		Email:    "test@" + helpers.RandStringRunes(5) + "." + helpers.RandStringRunes(3),
		Password: "123456",
	}
	assert.NoError(t, serviceUsers.Create(user))
	testCases := make([]*testCaseUser, 0)

	testCases = append(testCases, &testCaseUser{
		UserId:   fmt.Sprint(user.UserId),
		Email:    user.Email,
		UserName: helpers.RandStringRunes(5),
		Want:     200,
	})
	testCases = append(testCases, &testCaseUser{
		UserId:   fmt.Sprint(user.UserId + 1),
		Email:    "test@test.te",
		UserName: "test",
		Want:     404,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			body := new(bytes.Buffer)
			multiPartWriter := multipart.NewWriter(body)

			assert.NoError(t, multiPartWriter.WriteField("userId", tc.UserId))
			assert.NoError(t, multiPartWriter.WriteField("email", tc.Email))
			assert.NoError(t, multiPartWriter.WriteField("name", tc.UserName))

			// Закрываем запись данных
			assert.NoError(t, multiPartWriter.Close())

			// Создаем объект реквеста
			req, err := http.NewRequest(http.MethodPut, "/api/v1/users/"+tc.UserId, body)
			assert.NoError(t, err)

			req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if !assert.Equal(t, tc.Want, w.Code) {
				t.Log("-----> Response: ", w.Body)
			}
		})
	}

	assert.NoError(t, server.Db.Debug().Delete(user).Error)
}

func TestGetCats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseCat, 0)
	a := &testCaseCat{Want: 200}
	b := &testCaseCat{MyStr: "?asTree=true", Want: 200}
	testCases = append(testCases, a, b)

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/api/v1/cats"+tc.MyStr, nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}
func TestGetCatsCatId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseCat, 0)
	serviceCats := service.NewCatService()
	cats, err := serviceCats.GetCats()
	assert.NoError(t, err)
	offset := 2

	if len(cats) > offset {
		cats = cats[len(cats)-offset:]
	}

	// создаим нужные нам автоматичекие testCase-ы
	for _, cat := range cats {
		a := &testCaseCat{MyStr: fmt.Sprint(cat.CatId), Want: 200}
		b := &testCaseCat{MyStr: fmt.Sprint(cat.CatId + uint64(offset)), Want: 404}
		testCases = append(testCases, a, b)
	}

	// прогоним тесты
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/api/v1/cats/"+tc.MyStr, nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}
func TestPostCats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceCats := service.NewCatService()
	testCases := make([]*testCaseCat, 0)

	testCases = append(testCases, &testCaseCat{
		RequestPost: &request.PostCat{
			Name: helpers.RandStringRunes(5),
		},
		Want: 201,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.RequestPost)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/api/v1/cats", bytes.NewBuffer(b))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if !assert.Equal(t, tc.Want, w.Code) {
				t.Log("-----> Response: ", w.Body)
			}

			if w.Code == 201 {
				cat := new(storage.Cat)
				assert.NoError(t, json.Unmarshal(w.Body.Bytes(), cat))
				assert.NoError(t, serviceCats.Delete(cat.CatId))
			}
		})
	}
}
func TestPutCatsCatId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceCats := service.NewCatService()
	cat := &storage.Cat{
		Name: helpers.RandStringRunes(3),
	}
	assert.NoError(t, serviceCats.Create(cat))
	testCases := make([]*testCaseCat, 0)

	testCases = append(testCases, &testCaseCat{
		RequestPut: &request.PutCat{
			CatId: cat.CatId,
			Name:  helpers.RandStringRunes(5),
		},
		Want: 200,
	})
	testCases = append(testCases, &testCaseCat{
		RequestPut: &request.PutCat{
			CatId: cat.CatId + 1,
			Name:  cat.Name,
		},
		Want: 404,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.RequestPut)
			assert.NoError(t, err)
			url := "/api/v1/cats/" + fmt.Sprint(tc.RequestPut.CatId)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(b))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if !assert.Equal(t, tc.Want, w.Code) {
				t.Log("-----> Response: ", w.Body)
			}
		})
	}

	assert.NoError(t, serviceCats.Delete(cat.CatId))
}
func TestDeleteCatsCatId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceCats := service.NewCatService()
	testCases := make([]*testCaseCat, 0)

	testCases = append(testCases, &testCaseCat{
		Cat: &storage.Cat{
			Name: "test",
		},
		Want: 204,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			assert.NoError(t, serviceCats.Create(tc.Cat))
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/api/v1/cats/"+fmt.Sprint(tc.Cat.CatId), nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}

func TestGetAds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseAd, 0)

	testCases = append(testCases, &testCaseAd{Want: 200})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/api/v1/ads", nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}
func TestGetAdsAdId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseAd, 0)
	serviceAds := service.NewAdService()
	ads, err := serviceAds.GetAds()
	assert.NoError(t, err)
	offset := 2

	if len(ads) > offset {
		ads = ads[len(ads)-offset:]
	}

	// создаим нужные нам автоматичекие testCase-ы
	for _, ad := range ads {
		a := &testCaseAd{MyStr: fmt.Sprint(ad.AdId), Want: 200}
		b := &testCaseAd{MyStr: fmt.Sprint(ad.AdId + uint64(offset)), Want: 404}
		testCases = append(testCases, a, b)
	}

	// прогоним тесты
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/api/v1/ads/"+tc.MyStr, nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}
func TestPostAds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceAds := service.NewAdService()
	testCases := make([]*testCaseAd, 0)

	testCases = append(testCases, &testCaseAd{
		RequestPost: &request.PostAd{
			Title: helpers.RandStringRunes(5),
			CatId: 1,
		},
		Want: 201,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			body := new(bytes.Buffer)
			multiPartWriter := multipart.NewWriter(body)

			assert.NoError(t, multiPartWriter.WriteField("title", tc.RequestPost.Title))
			assert.NoError(t, multiPartWriter.WriteField("catId", fmt.Sprint(tc.RequestPost.CatId)))

			// Закрываем запись данных
			assert.NoError(t, multiPartWriter.Close())

			// Создаем объект реквеста
			req, err := http.NewRequest(http.MethodPost, "/api/v1/ads", body)
			assert.NoError(t, err)

			req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if !assert.Equal(t, tc.Want, w.Code) {
				t.Log("-----> Response: ", w.Body)
			}

			if w.Code == 201 {
				ad := new(storage.Ad)
				assert.NoError(t, json.Unmarshal(w.Body.Bytes(), ad))
				assert.NoError(t, serviceAds.Delete(ad.AdId))
			}
		})
	}
}
func TestPutAdsAdId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceAds := service.NewAdService()
	testCases := make([]*testCaseAd, 0)
	ad := &storage.Ad{
		Title: helpers.RandStringRunes(3),
		CatId: 1,
	}
	assert.NoError(t, serviceAds.Create(ad))

	testCases = append(testCases, &testCaseAd{
		RequestPut: &request.PutAd{
			AdId:  ad.AdId,
			Title: helpers.RandStringRunes(5),
			CatId: 2,
		},
		Want: 200,
	})
	testCases = append(testCases, &testCaseAd{
		RequestPut: &request.PutAd{
			AdId:  ad.AdId + 1,
			Title: ad.Title,
			CatId: 2,
		},
		Want: 404,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			body := new(bytes.Buffer)
			multiPartWriter := multipart.NewWriter(body)
			sAdId := fmt.Sprint(tc.RequestPut.AdId)
			assert.NoError(t, multiPartWriter.WriteField("adId", sAdId))
			assert.NoError(t, multiPartWriter.WriteField("title", tc.RequestPut.Title))
			assert.NoError(t, multiPartWriter.WriteField("catId", fmt.Sprint(tc.RequestPut.CatId)))

			// Закрываем запись данных
			assert.NoError(t, multiPartWriter.Close())

			// Создаем объект реквеста
			req, err := http.NewRequest(http.MethodPut, "/api/v1/ads/"+sAdId, body)
			assert.NoError(t, err)

			req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if !assert.Equal(t, tc.Want, w.Code) {
				t.Log("-----> Response: ", w.Body)
			}
		})
	}

	assert.NoError(t, serviceAds.Delete(ad.AdId))
}
func TestDeleteAdsAdId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceAds := service.NewAdService()
	testCases := make([]*testCaseAd, 0)

	testCases = append(testCases, &testCaseAd{
		Ad: &storage.Ad{
			Title: "test",
			CatId: 1,
		},
		Want: 204,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			assert.NoError(t, serviceAds.Create(tc.Ad))
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/api/v1/ads/"+fmt.Sprint(tc.Ad.AdId), nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}

func TestGetProperties(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseProperties, 0)
	a := &testCaseProperties{Want: 200}
	testCases = append(testCases, a)

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/api/v1/properties", nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}
func TestGetPropertiesPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseProperties, 0)
	serviceProperties := service.NewPropertyService()
	properties, err := serviceProperties.GetProperties(true)
	assert.NoError(t, err)
	offset := 2

	if len(properties) > offset {
		properties = properties[len(properties)-offset:]
	}

	// создаим нужные нам автоматичекие testCase-ы
	for _, v := range properties {
		a := &testCaseProperties{MyStr: fmt.Sprint(v.PropertyId), Want: 200}
		b := &testCaseProperties{MyStr: fmt.Sprint(v.PropertyId + uint64(offset)), Want: 404}
		c := &testCaseProperties{MyStr: "test", Want: 400}
		testCases = append(testCases, a, b, c)
	}

	// прогоним тесты
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/api/v1/properties/"+tc.MyStr, nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}
func TestPostProperties(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceProperties := service.NewPropertyService()
	serviceKindProperties := service.NewKindPropertyService()
	kindsProperties, err := serviceKindProperties.GetKindProperties()
	assert.NoError(t, err)

	myRequests := make([]*testCaseProperties, 0)
	a := &testCaseProperties{
		RequestPost: &request.PostProperty{
			Name:           "",
			KindPropertyId: 0,
		},
		Want: 400,
	}
	myRequests = append(myRequests, a)

	if len(kindsProperties) > 0 {
		tmp := &testCaseProperties{
			RequestPost: &request.PostProperty{
				Title:          helpers.RandStringRunes(5),
				Name:           helpers.RandStringRunes(5),
				KindPropertyId: kindsProperties[0].KindPropertyId,
			},
			Want: 201,
		}

		myRequests = append(myRequests, tmp)
	}

	for _, tc := range myRequests {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.RequestPost)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/api/v1/properties", bytes.NewBuffer(b))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if !assert.Equal(t, tc.Want, w.Code) {
				t.Log("-----> Response: ", w.Body)
			}

			if w.Code == 201 {
				p := new(storage.Property)
				assert.NoError(t, json.Unmarshal(w.Body.Bytes(), p))
				assert.NoError(t, serviceProperties.Delete(p.PropertyId))
			}
		})
	}
}
func TestPutPropertiesPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceProperties := service.NewPropertyService()
	serviceKindProperties := service.NewKindPropertyService()

	kindsProperties, err := serviceKindProperties.GetKindProperties()
	assert.NoError(t, err)

	if len(kindsProperties) < 1 {
		assert.Equal(t, 200, 200)
		return
	}

	pr := &storage.Property{
		Title:          helpers.RandStringRunes(5),
		Name:           helpers.RandStringRunes(5),
		KindPropertyId: kindsProperties[0].KindPropertyId,
	}

	assert.NoError(t, serviceProperties.Create(pr))
	defer func() {
		assert.NoError(t, serviceProperties.Delete(pr.PropertyId))
	}()

	testCases := make([]*testCaseProperties, 0)

	testCases = append(testCases, &testCaseProperties{
		RequestPut: &request.PutProperty{
			PropertyId:     pr.PropertyId,
			Title:          helpers.RandStringRunes(5),
			Name:           helpers.RandStringRunes(5),
			KindPropertyId: pr.KindPropertyId,
		},
		Want: 200,
	})
	testCases = append(testCases, &testCaseProperties{
		RequestPut: &request.PutProperty{
			PropertyId:     pr.PropertyId,
			Title:          "",
			Name:           helpers.RandStringRunes(5),
			KindPropertyId: pr.KindPropertyId,
		},
		Want: 400,
	})
	testCases = append(testCases, &testCaseProperties{
		RequestPut: &request.PutProperty{
			PropertyId:     pr.PropertyId + 1,
			Title:          helpers.RandStringRunes(5),
			Name:           pr.Name,
			KindPropertyId: pr.KindPropertyId,
		},
		Want: 404,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.RequestPut)
			assert.NoError(t, err)
			url := "/api/v1/properties/" + fmt.Sprint(tc.RequestPut.PropertyId)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(b))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if !assert.Equal(t, tc.Want, w.Code) {
				t.Log("-----> Response: ", w.Body)
			}
		})
	}
}
func TestDeletePropertiesPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceProperties := service.NewPropertyService()
	serviceKindProperties := service.NewKindPropertyService()
	kindsProperties, err := serviceKindProperties.GetKindProperties()
	assert.NoError(t, err)

	if len(kindsProperties) < 0 {
		assert.Equal(t, 200, 200)
		return
	}

	testCases := make([]*testCaseProperties, 0)

	testCases = append(testCases, &testCaseProperties{
		Pr: &storage.Property{
			Name:           helpers.RandStringRunes(5),
			KindPropertyId: kindsProperties[0].KindPropertyId,
		},
		Want: 204,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			assert.NoError(t, serviceProperties.Create(tc.Pr))
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/api/v1/properties/"+fmt.Sprint(tc.Pr.PropertyId), nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}

func TestGetKindProperties(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseKindProperties, 0)
	a := &testCaseKindProperties{Want: 200}
	testCases = append(testCases, a)

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/api/v1/kind_properties", nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}
func TestGetKindPropertiesKindPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	testCases := make([]*testCaseKindProperties, 0)
	serviceKindProperties := service.NewKindPropertyService()
	kindProperties, err := serviceKindProperties.GetKindProperties()
	assert.NoError(t, err)
	offset := 2

	if len(kindProperties) > offset {
		kindProperties = kindProperties[len(kindProperties)-offset:]
	}

	// создаим нужные нам автоматичекие testCase-ы
	for _, v := range kindProperties {
		a := &testCaseKindProperties{MyStr: fmt.Sprint(v.KindPropertyId), Want: 200}
		b := &testCaseKindProperties{MyStr: fmt.Sprint(v.KindPropertyId + uint64(offset)), Want: 404}
		testCases = append(testCases, a, b)
	}

	// прогоним тесты
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/api/v1/kind_properties/"+tc.MyStr, nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}
func TestPostKindProperties(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceKindProperties := service.NewKindPropertyService()
	testCases := make([]*testCaseKindProperties, 0)

	testCases = append(testCases, &testCaseKindProperties{
		RequestPost: &request.PostKindProperty{
			Name: helpers.RandStringRunes(5),
		},
		Want: 201,
	})
	testCases = append(testCases, &testCaseKindProperties{
		RequestPost: &request.PostKindProperty{
			Name: "",
		},
		Want: 400,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.RequestPost)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/api/v1/kind_properties", bytes.NewBuffer(b))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if !assert.Equal(t, tc.Want, w.Code) {
				t.Log("-----> Response: ", w.Body)
			}

			if w.Code == 201 {
				kp := new(storage.KindProperty)
				assert.NoError(t, json.Unmarshal(w.Body.Bytes(), kp))
				assert.NoError(t, serviceKindProperties.Delete(kp.KindPropertyId))
			}
		})
	}
}
func TestPutKindPropertiesKindPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceKindProperties := service.NewKindPropertyService()
	kp := &storage.KindProperty{
		Name: helpers.RandStringRunes(3),
	}
	assert.NoError(t, serviceKindProperties.Create(kp))
	defer func() { // именно так
		assert.NoError(t, serviceKindProperties.Delete(kp.KindPropertyId))
	}()

	testCases := make([]*testCaseKindProperties, 0)

	testCases = append(testCases, &testCaseKindProperties{
		RequestPut: &request.PutKindProperty{
			KindPropertyId: kp.KindPropertyId,
			Name:           helpers.RandStringRunes(5),
		},
		Want: 200,
	})
	testCases = append(testCases, &testCaseKindProperties{
		RequestPut: &request.PutKindProperty{
			KindPropertyId: kp.KindPropertyId + 1,
			Name:           kp.Name,
		},
		Want: 404,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.RequestPut)
			assert.NoError(t, err)
			url := "/api/v1/kind_properties/" + fmt.Sprint(tc.RequestPut.KindPropertyId)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(b))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if !assert.Equal(t, tc.Want, w.Code) {
				t.Log("-----> Response: ", w.Body)
			}
		})
	}
}
func TestDeleteKindPropertiesKindPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()
	serviceKindProperties := service.NewKindPropertyService()
	testCases := make([]*testCaseKindProperties, 0)

	testCases = append(testCases, &testCaseKindProperties{
		Kp: &storage.KindProperty{
			Name: "test" + helpers.RandStringRunes(5),
		},
		Want: 204,
	})

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			assert.NoError(t, serviceKindProperties.Create(tc.Kp))
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/api/v1/kind_properties/"+fmt.Sprint(tc.Kp.KindPropertyId), nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.Want, w.Code)
		})
	}
}

type testCaseUser struct {
	MyStr       string
	Want        int
	RequestPost *request.PostUser
	UserId      string
	Email       string
	UserName    string
}
type testCaseCat struct {
	MyStr       string
	Want        int
	RequestPost *request.PostCat
	RequestPut  *request.PutCat
	Cat         *storage.Cat
	CatId       string
	CatNewName  string
}
type testCaseAd struct {
	MyStr       string
	Want        int
	RequestPost *request.PostAd
	RequestPut  *request.PutAd
	Ad          *storage.Ad
}
type testCaseProperties struct {
	MyStr       string
	Want        int
	RequestPost *request.PostProperty
	RequestPut  *request.PutProperty
	Pr          *storage.Property
}
type testCaseKindProperties struct {
	MyStr       string
	Want        int
	RequestPost *request.PostKindProperty
	RequestPut  *request.PutKindProperty
	Kp          *storage.KindProperty
}
