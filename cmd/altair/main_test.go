package main

import (
	"altair/api/request"
	"altair/api/response"
	"altair/pkg/manager"
	"altair/pkg/service"
	"altair/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
)

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()

	tests := []struct {
		Want int
	}{
		{200},
	}

	for _, tt := range tests {
		t.Run("Get all user", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestGetUsersUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceUsers := service.NewUserService()
	var userID uint64
	type My struct {
		UserID string
		Want   int
	}

	usersDesc, err := serviceUsers.GetUsers("created_at desc")
	if !a.NoError(err) {
		return
	}
	if len(usersDesc) > 0 {
		userID = usersDesc[0].UserID
	}

	tests := []My{
		{fmt.Sprint(userID + 1), 404},
		{"test", 500},
	}

	if len(usersDesc) > 0 {
		tests = append(tests, My{fmt.Sprint(userID), 200})
	}

	for _, tt := range tests {
		t.Run("Get one user", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/users/"+tt.UserID, nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestPostUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceUsers := service.NewUserService()

	tests := []struct {
		Post request.PostUser
		Want int
	}{
		{
			request.PostUser{
				Email:            "test@" + manager.RandStringRunes(5) + "." + manager.RandStringRunes(3),
				Password:         "123456",
				PasswordConfirm:  "123456",
				Name:             manager.RandStringRunes(5),
				IsEmailConfirmed: true,
			}, 201,
		},
		{
			request.PostUser{
				Email:            "@" + manager.RandStringRunes(5) + "." + manager.RandStringRunes(3),
				Password:         "123456",
				PasswordConfirm:  "123456",
				Name:             manager.RandStringRunes(5),
				IsEmailConfirmed: true,
			}, 500,
		},
		{
			request.PostUser{
				Email:            "test@" + manager.RandStringRunes(5) + "." + manager.RandStringRunes(3),
				Password:         "12345",
				PasswordConfirm:  "123456",
				Name:             manager.RandStringRunes(5),
				IsEmailConfirmed: true,
			}, 400,
		},
	}

	for _, tt := range tests {
		t.Run("Create one user", func(t *testing.T) {
			body := new(bytes.Buffer)

			multiPartWriter := multipart.NewWriter(body)
			_ = multiPartWriter.WriteField("email", tt.Post.Email)
			_ = multiPartWriter.WriteField("password", tt.Post.Password)
			_ = multiPartWriter.WriteField("passwordConfirm", tt.Post.PasswordConfirm)
			_ = multiPartWriter.WriteField("name", tt.Post.Name)
			_ = multiPartWriter.WriteField("isEmailConfirmed", fmt.Sprint(tt.Post.IsEmailConfirmed))
			_ = multiPartWriter.Close()

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users", body)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code, w.Body)

			if w.Code == 201 {
				user := new(storage.User)
				if a.NoError(json.Unmarshal(w.Body.Bytes(), user)) {
					a.NoError(serviceUsers.Delete(user.UserID, nil))
				}
			}
		})
	}
}
func TestPutUsersUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceUsers := service.NewUserService()
	user := &storage.User{
		Email:    "test@" + manager.RandStringRunes(5) + "." + manager.RandStringRunes(3),
		Password: "123456",
	}

	err := a.NoError(serviceUsers.Create(user, nil))
	defer func() {
		a.NoError(serviceUsers.Delete(user.UserID, nil))
	}()
	if !err {
		return
	}

	tests := []struct {
		UserID   string
		Email    string
		UserName string
		Want     int
	}{
		{fmt.Sprint(user.UserID), user.Email, manager.RandStringRunes(5), 200},
		{fmt.Sprint(user.UserID + 1), "test@test.te", "test", 404},
	}

	for _, tt := range tests {
		t.Run("Edit user", func(t *testing.T) {
			body := new(bytes.Buffer)

			multiPartWriter := multipart.NewWriter(body)
			_ = multiPartWriter.WriteField("userId", tt.UserID)
			_ = multiPartWriter.WriteField("email", tt.Email)
			_ = multiPartWriter.WriteField("name", tt.UserName)
			_ = multiPartWriter.Close()

			// Создаем объект реквеста
			req, err := http.NewRequest(http.MethodPut, "/api/v1/users/"+tt.UserID, body)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code, w.Body)
		})
	}
}
func TestDeleteUsersUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceUsers := service.NewUserService()
	user := &storage.User{
		Email:    "test@" + manager.RandStringRunes(5) + "." + manager.RandStringRunes(3),
		Password: "123456",
	}

	if !a.NoError(serviceUsers.Create(user, nil)) {
		return
	}

	tests := []struct {
		UserID string
		Want   int
	}{
		{fmt.Sprint(user.UserID), 204},
		{fmt.Sprint(user.UserID + 1), 204},
		{"test", 500},
	}

	for _, tt := range tests {
		t.Run("Delete user", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/users/"+tt.UserID, nil)
			if !a.NoError(err) {
				return
			}

			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}

func TestGetKindProps(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()

	tests := []struct {
		Want int
	}{
		{200},
	}

	for _, tt := range tests {
		t.Run("Get all kind props", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/kind_props", nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestGetKindPropsKindPropID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceKindProps := service.NewKindPropService()
	var elID uint64 = 0
	type My struct {
		ElID string
		Want int
	}

	kindProps, err := serviceKindProps.GetKindProps("kind_prop_id desc")
	if !a.NoError(err) {
		return
	}

	if len(kindProps) > 0 {
		elID = kindProps[0].KindPropID
	}

	tests := []My{
		{fmt.Sprint(elID + uint64(1)), 404},
	}

	if len(kindProps) > 0 {
		tests = append(tests, My{fmt.Sprint(elID), 200})
	}

	for _, tt := range tests {
		t.Run("Get KindPropID one", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/kind_props/"+tt.ElID, nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestPostKindProps(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceKindProps := service.NewKindPropService()

	tests := []struct {
		Post request.PostKindProp
		Want int
	}{
		{request.PostKindProp{Name: manager.RandStringRunes(5)}, 201},
		{request.PostKindProp{Name: ""}, 400},
	}

	for _, tt := range tests {
		t.Run("Create one KindProp", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Post)
			if !a.NoError(err) {
				return
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/kind_props", bytes.NewBuffer(b))
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code, w.Body)

			if w.Code == 201 {
				kp := new(storage.KindProp)

				if a.NoError(json.Unmarshal(w.Body.Bytes(), kp)) {
					a.NoError(serviceKindProps.Delete(kp.KindPropID, nil))
				}
			}
		})
	}
}
func TestPutKindPropsKindPropID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceKindProps := service.NewKindPropService()
	kp := &storage.KindProp{
		Name: manager.RandStringRunes(3),
	}

	err := a.NoError(serviceKindProps.Create(kp, nil))
	defer func() {
		a.NoError(serviceKindProps.Delete(kp.KindPropID, nil))
	}()
	if !err {
		return
	}

	tests := []struct {
		Put  request.PutKindProp
		Want int
	}{
		{
			request.PutKindProp{
				KindPropID: kp.KindPropID,
				Name:       manager.RandStringRunes(5),
			},
			200,
		},
		{
			request.PutKindProp{
				KindPropID: kp.KindPropID + 1,
				Name:       kp.Name,
			},
			404,
		},
	}

	for _, tt := range tests {
		t.Run("Edit one KindProp", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Put)
			if !a.NoError(err) {
				return
			}

			url := "/api/v1/kind_props/" + fmt.Sprint(tt.Put.KindPropID)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(b))
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code, w.Body)
		})
	}
}
func TestDeleteKindPropsKindPropID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceKindProps := service.NewKindPropService()

	tests := []struct {
		Kp   *storage.KindProp
		Want int
	}{
		{&storage.KindProp{Name: "test" + manager.RandStringRunes(5)}, 204},
	}

	for _, tt := range tests {
		t.Run("Delete one KindProp", func(t *testing.T) {
			if !a.NoError(serviceKindProps.Create(tt.Kp, nil)) {
				return
			}

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/kind_props/"+fmt.Sprint(tt.Kp.KindPropID), nil)
			if !a.NoError(err) {
				return
			}

			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}

func TestGetProps(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()

	tests := []struct {
		Want int
	}{
		{200},
	}

	for _, tt := range tests {
		t.Run("Get all props", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/props", nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestGetPropsPropID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceProps := service.NewPropService()
	var elID uint64 = 0
	type My struct {
		PropID string
		Want   int
	}

	props, err := serviceProps.GetProps("prop_id desc")
	if !a.NoError(err) {
		return
	}

	if len(props) > 0 {
		elID = props[0].PropID
	}

	tests := []My{
		{fmt.Sprint(elID + uint64(1)), 404},
		{"test", 400},
	}

	if len(props) > 0 {
		tests = append(tests, My{fmt.Sprint(elID), 200})
	}

	for _, tt := range tests {
		t.Run("GET PropID", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/props/"+tt.PropID, nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestPostProps(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceProps := service.NewPropService()

	tests := []struct {
		Post request.PostProp
		Want int
	}{
		{request.PostProp{}, 400},
		{
			request.PostProp{
				Title:      manager.RandStringRunes(5),
				Name:       manager.RandStringRunes(5),
				KindPropID: "1",
			},
			201,
		},
	}

	for _, tt := range tests {
		t.Run("POST Props", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Post)
			if !a.NoError(err) {
				return
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/props", bytes.NewBuffer(b))
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code, w.Body)

			if w.Code == 201 {
				p := new(storage.Prop)
				if a.NoError(json.Unmarshal(w.Body.Bytes(), p)) {
					a.NoError(serviceProps.Delete(p.PropID, nil))
				}
			}
		})
	}
}
func TestPutPropsPropID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceProps := service.NewPropService()
	pr := &storage.Prop{
		Title:      manager.RandStringRunes(5),
		Name:       manager.RandStringRunes(5),
		KindPropID: 1,
	}

	noError := a.NoError(serviceProps.Create(pr, nil))
	defer func() {
		assert.NoError(t, serviceProps.Delete(pr.PropID, nil))
	}()
	if !noError {
		return
	}

	tests := []struct {
		Put  request.PutProp
		Want int
	}{
		{
			request.PutProp{
				PropID:     pr.PropID,
				Title:      manager.RandStringRunes(5),
				Name:       manager.RandStringRunes(5),
				KindPropID: "2",
			},
			200,
		},
		{
			request.PutProp{
				PropID:     pr.PropID,
				Title:      "",
				Name:       manager.RandStringRunes(5),
				KindPropID: "3",
			},
			400,
		},
		{
			request.PutProp{
				PropID:     pr.PropID + 1,
				Title:      manager.RandStringRunes(5),
				Name:       pr.Name,
				KindPropID: "4",
			},
			404,
		},
	}

	for _, tt := range tests {
		t.Run("Put Prop", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Put)
			if !a.NoError(err) {
				return
			}

			url := "/api/v1/props/" + fmt.Sprint(tt.Put.PropID)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(b))
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code, w.Body)
		})
	}
}
func TestDeletePropsPropID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceProps := service.NewPropService()

	tests := []struct {
		Pr   *storage.Prop
		Want int
	}{
		{
			&storage.Prop{
				Name:       manager.RandStringRunes(5),
				KindPropID: 1,
			},
			204,
		},
	}

	for _, tt := range tests {
		t.Run("DELETE PropID", func(t *testing.T) {
			noError := a.NoError(serviceProps.Create(tt.Pr, nil))
			if !noError {
				return
			}

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/props/"+fmt.Sprint(tt.Pr.PropID), nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}

func TestGetCats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()

	tests := []struct {
		Query string
		Want  int
	}{
		{"", 200},
		{"asTree=true", 200},
		{"asTree=1", 200},
	}

	for _, tt := range tests {
		t.Run("GET cats", func(t *testing.T) {
			w := httptest.NewRecorder()
			query := ""

			if tt.Query != "" {
				query += "?" + tt.Query
			}

			req, err := http.NewRequest(http.MethodGet, "/api/v1/cats"+query, nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestGetCatsCatID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceCats := service.NewCatService()
	var elID uint64
	type My struct {
		CatID string
		Query string
		Want  int
	}

	cats, err := serviceCats.GetCats(0)
	if !a.NoError(err) {
		return
	}

	// возьмем самый последний
	for _, v := range cats {
		if elID < v.CatID {
			elID = v.CatID
		}
	}

	tests := []My{
		{fmt.Sprint(elID + 1), "", 404},
		{fmt.Sprint(elID + 2), "withPropsOnlyFiltered=true", 404},
		{fmt.Sprint(elID + 3), "withPropsOnlyFiltered=1", 404},
	}

	if len(cats) > 0 {
		tests = append(
			tests,
			My{fmt.Sprint(elID), "", 200},
			My{fmt.Sprint(elID), "withPropsOnlyFiltered=true", 200},
			My{fmt.Sprint(elID), "withPropsOnlyFiltered=1", 200},
		)
	}

	for _, tt := range tests {
		t.Run("Get one cat", func(t *testing.T) {
			w := httptest.NewRecorder()
			query := ""

			if tt.Query != "" {
				query += "?" + query
			}

			req, err := http.NewRequest(http.MethodGet, "/api/v1/cats/"+tt.CatID+query, nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestPostCats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceCats := service.NewCatService()
	type My struct {
		Post request.PostCat
		Want int
	}

	tests := []My{
		{request.PostCat{Name: manager.RandStringRunes(5)}, 201},
		{request.PostCat{}, 400},
	}

	for _, tt := range tests {
		t.Run("POST cat", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Post)
			if !a.NoError(err) {
				return
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/cats", bytes.NewBuffer(b))
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code, w.Body)

			if w.Code == 201 {
				cat := new(storage.Cat)
				if a.NoError(json.Unmarshal(w.Body.Bytes(), cat)) {
					a.NoError(serviceCats.Delete(cat.CatID, nil))
				}
			}
		})
	}
}
func TestPutCatsCatID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceCats := service.NewCatService()
	cat := &storage.Cat{
		Name: manager.RandStringRunes(3),
	}

	noError := a.NoError(serviceCats.Create(cat, nil))
	defer func() {
		a.NoError(serviceCats.Delete(cat.CatID, nil))
	}()
	if !noError {
		return
	}

	tests := []struct {
		Put  request.PutCat
		Want int
	}{
		{request.PutCat{CatID: cat.CatID, Name: manager.RandStringRunes(5)}, 200},
		{request.PutCat{CatID: cat.CatID + 1, Name: cat.Name}, 404},
	}

	for _, tt := range tests {
		t.Run("PUT cat", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Put)
			if !a.NoError(err) {
				return
			}

			url := "/api/v1/cats/" + fmt.Sprint(tt.Put.CatID)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(b))
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code, w.Body)
		})
	}
}
func TestDeleteCatsCatID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceCats := service.NewCatService()

	tests := []struct {
		Cat  *storage.Cat
		Want int
	}{
		{&storage.Cat{Name: "test"}, 204},
	}

	for _, tt := range tests {
		t.Run("DELETE cat", func(t *testing.T) {
			if !a.NoError(serviceCats.Create(tt.Cat, nil)) {
				return
			}

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/cats/"+fmt.Sprint(tt.Cat.CatID), nil)
			if !a.NoError(err) {
				return
			}

			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}

func TestGetAds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceCats := service.NewCatService()
	var catID uint64

	cats, err := serviceCats.GetCats(0)
	if !a.NoError(err) {
		return
	}

	// возьмем самый последний
	for _, v := range cats {
		if catID < v.CatID {
			catID = v.CatID
		}
	}

	tests := []struct {
		CatID string
		Want  int
	}{
		{"catID=" + fmt.Sprint(catID), 200},
		{"", 200},
		{"catID=test", 500},
	}

	for _, tt := range tests {
		t.Run("GET ads", func(t *testing.T) {
			w := httptest.NewRecorder()
			query := ""

			if tt.CatID != "" {
				query += "?" + tt.CatID
			}

			req, err := http.NewRequest(http.MethodGet, "/api/v1/ads"+query, nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestGetAdsAdID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceAds := service.NewAdService()
	var adID uint64
	type My struct {
		AdID string
		Want int
	}

	ads, err := serviceAds.GetAds("created_at desc")
	if !a.NoError(err) {
		return
	}

	if len(ads) > 0 {
		adID = ads[0].AdID
	}

	tests := []My{
		{fmt.Sprint(adID + 1), 404},
	}

	if len(ads) > 0 {
		tests = append(tests, My{fmt.Sprint(adID), 200})
	}

	for _, tt := range tests {
		t.Run("Get ad one", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/ads/"+tt.AdID, nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestPostAds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceAds := service.NewAdService()

	tests := []struct {
		Post request.PostAd
		Want int
	}{
		{request.PostAd{
			Title:       manager.RandStringRunes(10),
			CatID:       1,
			Description: manager.RandStringRunes(5)}, 400},
		{request.PostAd{}, 400},
	}

	for _, tt := range tests {
		t.Run("POST ad", func(t *testing.T) {
			body := new(bytes.Buffer)

			multiPartWriter := multipart.NewWriter(body)
			_ = multiPartWriter.WriteField("title", fmt.Sprint(tt.Post.Title))
			_ = multiPartWriter.WriteField("catID", fmt.Sprint(tt.Post.CatID))
			_ = multiPartWriter.WriteField("description", tt.Post.Description)
			_ = multiPartWriter.Close()

			req, err := http.NewRequest(http.MethodPost, "/api/v1/ads", body)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code, w.Body)

			if w.Code == 201 {
				ad := new(storage.Ad)
				if a.NoError(json.Unmarshal(w.Body.Bytes(), ad)) {
					a.NoError(serviceAds.Delete(ad.AdID, nil))
				}
			}
		})
	}
}
func TestPutAdsAdID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceAds := service.NewAdService()
	ad := &storage.Ad{
		Title:       manager.RandStringRunes(10),
		CatID:       1,
		Description: manager.RandStringRunes(3),
	}

	noError := a.NoError(serviceAds.Create(ad, nil))
	defer func() {
		a.NoError(serviceAds.Delete(ad.AdID, nil))
	}()
	if !noError {
		return
	}

	tests := []struct {
		Put  request.PutAd
		Want int
	}{
		{
			request.PutAd{
				AdID: ad.AdID,
				PostAd: request.PostAd{
					Title:       manager.RandStringRunes(10),
					CatID:       2,
					Description: manager.RandStringRunes(5),
				},
			}, 400,
		},
		{
			request.PutAd{
				AdID: ad.AdID + 1,
				PostAd: request.PostAd{
					Title:       manager.RandStringRunes(10),
					CatID:       2,
					Description: manager.RandStringRunes(5),
				},
			}, 404,
		},
		{request.PutAd{}, 400},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			body := new(bytes.Buffer)
			sAdID := fmt.Sprint(tt.Put.AdID)

			multiPartWriter := multipart.NewWriter(body)
			_ = multiPartWriter.WriteField("adID", sAdID)
			_ = multiPartWriter.WriteField("title", tt.Put.Title)
			_ = multiPartWriter.WriteField("catID", fmt.Sprint(tt.Put.CatID))
			_ = multiPartWriter.WriteField("description", tt.Put.Description)
			_ = multiPartWriter.Close()

			req, err := http.NewRequest(http.MethodPut, "/api/v1/ads/"+sAdID, body)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code, w.Body)
		})
	}
}
func TestDeleteAdsAdID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceAds := service.NewAdService()

	tests := []struct {
		Ad   *storage.Ad
		Want int
	}{
		{
			&storage.Ad{
				Title:       manager.RandStringRunes(10),
				CatID:       1,
				Description: manager.RandStringRunes(5),
			}, 204,
		},
	}

	for _, tt := range tests {
		t.Run("Delete ad", func(t *testing.T) {
			if !a.NoError(serviceAds.Create(tt.Ad, nil)) {
				return
			}

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/ads/"+fmt.Sprint(tt.Ad.AdID), nil)
			if !a.NoError(err) {
				return
			}

			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}

func TestSearchAds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()

	tests := []struct {
		Query string
		Want  int
	}{
		{manager.RandStringRunes(10), 200},
		{"", 400},
		{"t", 400},
	}

	for _, tc := range tests {
		t.Run("GET Search", func(t *testing.T) {
			w := httptest.NewRecorder()
			var query string

			if tc.Query != "" {
				query += "q=" + tc.Query
			}
			if query != "" {
				query = "?" + query
			}

			req, err := http.NewRequest(http.MethodGet, "/api/v1/search/ads"+query, nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tc.Want, w.Code)
		})
	}
}

/* func TestAllCatsOnWrite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceCat := service.NewCatService()
	maxUserConnectionsToMySQL := 5

	cats, err := serviceCat.GetCats(0)
	if !a.NoError(err) {
		return
	}

	//testCat(479, t, r)
	catTree := serviceCat.GetCatsAsTree(cats)
	myCh := make(chan struct{}, maxUserConnectionsToMySQL)
	var wg sync.WaitGroup
	walkToCatTree(catTree.Childes, t, r, &wg, myCh)
	wg.Wait()
} */

func walkToCatTree(list []*response.CatTree, t *testing.T, r *gin.Engine, wg *sync.WaitGroup, ch chan struct{}) {
	for _, leaf := range list {
		// если это ветка
		if len(leaf.Childes) > 0 {
			walkToCatTree(leaf.Childes, t, r, wg, ch)
			continue
		}
		// тут мы находимся в "листе"
		wg.Add(1)
		go func(catID uint64) {
			defer wg.Done()
			ch <- struct{}{} // благодаря каналам создаим нормальную очередь
			testCat(catID, t, r)
			<-ch
		}(leaf.CatID)
	}
}
func testCat(catID uint64, t *testing.T, r *gin.Engine) {
	a := assert.New(t)
	serviceCat := service.NewCatService()

	catFull, err := serviceCat.GetCatFullByID(catID, false, 0)
	if !a.NoError(err) {
		return
	}

	// тут надо создать карту с нужными данными
	receiver := make(map[string]string)
	receiver["title"] = manager.RandStringRunes(10)
	receiver["catID"] = fmt.Sprint(catFull.CatID)
	receiver["description"] = manager.RandStringRunes(10)
	receiver["price"] = "0"
	receiver["youtube"] = manager.RandStringRunes(10)

	// заполним карту доп. св-вами
	for _, v1 := range catFull.PropsFull {
		val := manager.RandStringRunes(5)

		switch v1.KindPropName {
		case "checkbox", "radio", "select":
			for _, v2 := range v1.Values {
				val = fmt.Sprint(v2.ValueID)
				break
			}
		case "photo":
			val = v1.Comment
		case "input_number":
			val = "0"
		}

		receiver[v1.Name] = val
	}

	tests := []struct {
		Map  map[string]string
		Want int
	}{
		{
			receiver,
			201},
	}

	for _, tt := range tests {
		t.Run("POST ad (all)", func(t *testing.T) {
			body := new(bytes.Buffer)

			form := multipart.NewWriter(body)
			for k, v := range receiver {
				// если это файлы
				if k == "files" {
					maxFiles, err := strconv.Atoi(v)
					if !a.NoError(err) || maxFiles < 0 {
						continue
					}

					// для подстраховки установим лимит
					limitImages := 10
					if maxFiles > limitImages {
						maxFiles = limitImages
					}

					// обратимся к файлам
					dirTestImg := "./web/assets/img/test/"
					files, err := ioutil.ReadDir(dirTestImg)
					if !a.NoError(err) {
						continue
					}

					totalAcceptFiles := 0
					for _, file := range files {
						ext := filepath.Ext(file.Name())

						if file.IsDir() || ext != ".jpg" {
							continue
						}

						// берем только подходящее и нужное кол-во
						if totalAcceptFiles >= maxFiles {
							break
						}
						totalAcceptFiles++

						if err := attachFileInMultipart(form, dirTestImg+file.Name()); !a.NoError(err) {
							continue
						}
					}

					continue
				}
				_ = form.WriteField(k, v)
			}
			_ = form.Close()

			req, err := http.NewRequest(http.MethodPost, "/api/v1/ads", body)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", form.FormDataContentType())
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code, w.Body)
		})
	}
}
func attachFileInMultipart(mp *multipart.Writer, pathFile string) error {
	fileBase := filepath.Base(pathFile)
	fileExt := filepath.Ext(pathFile)

	if !manager.FolderOrFileExists(pathFile) {
		return fmt.Errorf("%s%s", "not found file: ", pathFile)
	}

	file, err := os.Open(pathFile)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		return err
	}

	h := textproto.MIMEHeader{}
	contentDisposition := fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "files", fileBase)
	h.Set("Content-Disposition", contentDisposition)
	h.Set("Content-Type", "application/octet-stream")

	if ct := mime.TypeByExtension(fileExt); ct != "" {
		h.Set("Content-Type", ct)
	}

	part, err := mp.CreatePart(h)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	return nil
}
