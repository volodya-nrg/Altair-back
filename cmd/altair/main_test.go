package main

import (
	"altair/api/request"
	"altair/pkg/helpers"
	"altair/pkg/service"
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
/*
assert.ElementsMatch(t, [1, 3, 2, 3], [1, 3, 3, 2])
assert.Empty(t, obj) // object is empty. I.e. nil, "", false, 0 or either a slice or a channel with len == 0.
assert.Equal(t, 123, 123)
assert.EqualError(t, err,  expectedErrorString)
assert.EqualValues(t, uint32(123), int32(123))

actualObj, err := SomeFunction()
if assert.Error(t, err) {
   	assert.Equal(t, expectedError, err)
}

assert.Exactly(t, int32(123), int64(123)) // проверка по значению и типу
assert.False(t, myBool) // является ли ложным

// проверка одного больше чем другое
assert.Greater(t, 2, 1)
assert.Greater(t, float64(2), float64(1))
assert.Greater(t, "b", "a")

// больше либо равно
assert.GreaterOrEqual(t, 2, 1)
assert.GreaterOrEqual(t, 2, 2)
assert.GreaterOrEqual(t, "b", "a")
assert.GreaterOrEqual(t, "b", "b")

// http-помошник. Пустая строка если не удачно.
func HTTPBody(handler http.HandlerFunc, method, url string, values url.Values) string

// проверка относительно строки-примера
assert.HTTPBodyContains(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")

//...которое не содержит строку.
assert.HTTPBodyNotContains(t, myHandler, "GET", "www.google.com", nil, "I'm Feeling Lucky")

// ...возвращает код состояния ошибки
assert.HTTPError(t, myHandler, "POST", "/a/b/c", url.Values{"a": []string{"b", "c"}})

// handler returns a redirect status code.
assert.HTTPRedirect(t, myHandler, "GET", "/a/b/c", url.Values{"a": []string{"b", "c"}}

// returns a specified status code
assert.HTTPStatusCode(t, myHandler, "GET", "/notImplemented", nil, 501)

// handler returns a success status code.
assert.HTTPSuccess(t, myHandler, "POST", "http://www.google.com", nil)

// проверка двух json-ов
assert.JSONEq(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)

// проверка на длину
assert.Len(t, mySlice, 3)

// проверка на меньшенство
assert.Less(t, 1, 2)
assert.Less(t, float64(1), float64(2))
assert.Less(t, "a", "b")

// меньше либо равен
assert.LessOrEqual(t, 1, 2)
assert.LessOrEqual(t, 2, 2)
assert.LessOrEqual(t, "a", "b")
assert.LessOrEqual(t, "b", "b")

// проверка объекта на nil
assert.Nil(t, err)

// проверка на то что нет ошибки
actualObj, err := SomeFunction()
if assert.NoError(t, err) {
   	assert.Equal(t, expectedObj, actualObj)
}

// не содержит
assert.NotContains(t, "Hello World", "Earth")
assert.NotContains(t, ["Hello", "World"], "Earth")
assert.NotContains(t, {"Hello": "World"}, "Earth")

// не пустой объект
if assert.NotEmpty(t, obj) {
  	assert.Equal(t, "two", obj[1])
}

// объекты не равны
assert.NotEqual(t, obj1, obj2)

// не равен nil
assert.NotNil(t, err)

// нет подустановки
assert.NotSubset(t, [1, 3, 4], [1, 2], "But [1, 3, 4] does not contain [1, 2]")

// не ноль
func NotZero(t TestingT, i interface{}, msgAndArgs ...interface{}) bool

// объекты равны
func ObjectsAreEqual(expected, actual interface{}) bool

// объекты имеют одни значения
func ObjectsAreEqualValues(expected, actual interface{}) bool

// один внутри другого
assert.Subset(t, [1, 2, 3], [1, 2], "But [1, 2, 3] does contain [1, 2]")

// явл. ли true
assert.True(t, myBool)

// два времени находятся в пределах друг от друга, исходя от заданного времени
assert.WithinDuration(t, time.Now(), time.Now(), 10*time.Second)

// ноль
func Zero(t TestingT, i interface{}, msgAndArgs ...interface{}) bool

// можно создать уже готовую структуру для тестирования
type Assertions struct {
    // contains filtered or unexported fields
}
func New(t TestingT) *Assertions
...и далее все вышеупомянутые ф-ии принадлежат тоже ей

// ф-ия проверка на булев тип
type BoolAssertionFunc func(TestingT, bool, ...interface{}) bool

t := &testing.T{} // provided by test
isOkay := func(x int) bool {
    return x >= 42
}
tests := []struct {
    name      string
    arg       int
    assertion BoolAssertionFunc
}{
    {"-1 is bad", -1, False},
    {"42 is good", 42, True},
    {"41 is bad", 41, False},
    {"45 is cool", 45, True},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        tt.assertion(t, isOkay(tt.arg))
    })
}

// еще одна ф-ия проверка
type ComparisonAssertionFunc func(TestingT, interface{}, interface{}, ...interface{})
t := &testing.T{} // provided by test
adder := func(x, y int) int {
    return x + y
}
type args struct {
    x   int
    y   int
}
tests := []struct {
    name      string
    args      args
    expect    int
    assertion ComparisonAssertionFunc
}{
    {"2+2=4", args{2, 2}, 4, Equal},
    {"2+2!=5", args{2, 2}, 5, NotEqual},
    {"2+3==5", args{2, 3}, 5, Exactly},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        tt.assertion(t, tt.expect, adder(tt.args.x, tt.args.y))
    })
}

// ф-ия на ошибку
type ErrorAssertionFunc func(TestingT, error, ...interface{}) bool
t := &testing.T{} // provided by test
dumbParseNum := func(input string, v interface{}) error {
    return json.Unmarshal([]byte(input), v)
}
tests := []struct {
    name      string
    arg       string
    assertion ErrorAssertionFunc
}{
    {"1.2 is number", "1.2", NoError},
    {"1.2.3 not number", "1.2.3", Error},
    {"true is not number", "true", Error},
    {"3 is number", "3", NoError},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        var x float64
        tt.assertion(t, dumbParseNum(tt.arg, &x))
    })
}

// проверка на значение
type ValueAssertionFunc func(TestingT, interface{}, ...interface{}) bool
t := &testing.T{} // provided by test
dumbParse := func(input string) interface{} {
    var x interface{}
    json.Unmarshal([]byte(input), &x)
    return x
}
tests := []struct {
    name      string
    arg       string
    assertion ValueAssertionFunc
}{
    {"true is not nil", "true", NotNil},
    {"empty string is nil", "", Nil},
    {"zero is not nil", "0", NotNil},
    {"zero is zero", "0", Zero},
    {"false is zero", "false", Zero},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        tt.assertion(t, dumbParse(tt.arg))
    })
}
*/

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	if !a.NoError(err) {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	a.Equal(200, w.Code)
}
func TestGetUsersUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceUsers := service.NewUserService()
	var userId uint64
	type My struct {
		UserId string
		Want   int
	}

	usersDesc, err := serviceUsers.GetUsers()
	if !a.NoError(err) {
		return
	}
	if len(usersDesc) > 0 {
		userId = usersDesc[0].UserId
	}

	tests := []My{
		{fmt.Sprint(userId + 1), 404},
		{"test", 500},
	}

	if len(usersDesc) > 0 {
		tests = append(tests, My{fmt.Sprint(userId), 200})
	}

	for _, tt := range tests {
		t.Run("Get one user", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/users/"+tt.UserId, nil)
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
				Email:           "test@" + helpers.RandStringRunes(5) + "." + helpers.RandStringRunes(3),
				Password:        "123456",
				PasswordConfirm: "123456",
				AgreeOffer:      true,
				AgreePolicy:     true,
			},
			201,
		},
		{
			request.PostUser{
				Email:           "test@" + helpers.RandStringRunes(5) + "." + helpers.RandStringRunes(3),
				Password:        "123456",
				PasswordConfirm: "123456",
				AgreeOffer:      false,
				AgreePolicy:     true,
			}, 400,
		},
		{
			request.PostUser{
				Email:           "test@" + helpers.RandStringRunes(5) + "." + helpers.RandStringRunes(3),
				Password:        "12345",
				PasswordConfirm: "123456",
				AgreeOffer:      true,
				AgreePolicy:     true,
			}, 400,
		},
	}

	for _, tt := range tests {
		t.Run("Create one user", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Post)
			if !a.NoError(err) {
				return
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(b))
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code, w.Body)

			if w.Code == 201 {
				user := new(storage.User)
				if a.NoError(json.Unmarshal(w.Body.Bytes(), user)) {
					a.NoError(serviceUsers.Delete(user.UserId, nil))
				}
			}
		})
	}
}
func TestPutUsersUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceUsers := service.NewUserService()
	user := &storage.User{
		Email:    "test@" + helpers.RandStringRunes(5) + "." + helpers.RandStringRunes(3),
		Password: "123456",
	}

	err := a.NoError(serviceUsers.Create(user, nil))
	defer func() {
		a.NoError(serviceUsers.Delete(user.UserId, nil))
	}()
	if !err {
		return
	}

	tests := []struct {
		UserId   string
		Email    string
		UserName string
		Want     int
	}{
		{fmt.Sprint(user.UserId), user.Email, helpers.RandStringRunes(5), 200},
		{fmt.Sprint(user.UserId + 1), "test@test.te", "test", 404},
	}

	for _, tt := range tests {
		t.Run("Edit user", func(t *testing.T) {
			body := new(bytes.Buffer)

			multiPartWriter := multipart.NewWriter(body)
			_ = multiPartWriter.WriteField("userId", tt.UserId)
			_ = multiPartWriter.WriteField("email", tt.Email)
			_ = multiPartWriter.WriteField("name", tt.UserName)
			_ = multiPartWriter.Close()

			// Создаем объект реквеста
			req, err := http.NewRequest(http.MethodPut, "/api/v1/users/"+tt.UserId, body)
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
func TestDeleteUsersUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceUsers := service.NewUserService()
	user := &storage.User{
		Email:    "test@" + helpers.RandStringRunes(5) + "." + helpers.RandStringRunes(3),
		Password: "123456",
	}

	if !a.NoError(serviceUsers.Create(user, nil)) {
		return
	}

	tests := []struct {
		UserId string
		Want   int
	}{
		{fmt.Sprint(user.UserId), 204},
		{fmt.Sprint(user.UserId + 1), 204},
		{"test", 500},
	}

	for _, tt := range tests {
		t.Run("Delete user", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/users/"+tt.UserId, nil)
			if !a.NoError(err) {
				return
			}

			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}

func TestGetKindProperties(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/api/v1/kind_properties", nil)
	if !a.NoError(err) {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	a.Equal(200, w.Code)
}
func TestGetKindPropertiesKindPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceKindProperties := service.NewKindPropertyService()
	var elId uint64 = 0
	type My struct {
		ElId string
		Want int
	}

	kindProperties, err := serviceKindProperties.GetKindProperties("kind_property_id desc")
	if !a.NoError(err) {
		return
	}

	if len(kindProperties) > 0 {
		elId = kindProperties[0].KindPropertyId
	}

	tests := []My{
		{fmt.Sprint(elId + uint64(1)), 404},
	}

	if len(kindProperties) > 0 {
		tests = append(tests, My{fmt.Sprint(elId), 200})
	}

	for _, tt := range tests {
		t.Run("Get KindPropertyId one", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/kind_properties/"+tt.ElId, nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestPostKindProperties(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceKindProperties := service.NewKindPropertyService()

	tests := []struct {
		Post request.PostKindProperty
		Want int
	}{
		{request.PostKindProperty{Name: helpers.RandStringRunes(5)}, 201},
		{request.PostKindProperty{Name: ""}, 400},
	}

	for _, tt := range tests {
		t.Run("Create one KindProperty", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Post)
			if !a.NoError(err) {
				return
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/kind_properties", bytes.NewBuffer(b))
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code, w.Body)

			if w.Code == 201 {
				kp := new(storage.KindProperty)

				if a.NoError(json.Unmarshal(w.Body.Bytes(), kp)) {
					a.NoError(serviceKindProperties.Delete(kp.KindPropertyId, nil))
				}
			}
		})
	}
}
func TestPutKindPropertiesKindPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceKindProperties := service.NewKindPropertyService()
	kp := &storage.KindProperty{
		Name: helpers.RandStringRunes(3),
	}

	err := a.NoError(serviceKindProperties.Create(kp, nil))
	defer func() {
		a.NoError(serviceKindProperties.Delete(kp.KindPropertyId, nil))
	}()
	if !err {
		return
	}

	tests := []struct {
		Put  request.PutKindProperty
		Want int
	}{
		{
			request.PutKindProperty{
				KindPropertyId: kp.KindPropertyId,
				Name:           helpers.RandStringRunes(5),
			},
			200,
		},
		{
			request.PutKindProperty{
				KindPropertyId: kp.KindPropertyId + 1,
				Name:           kp.Name,
			},
			404,
		},
	}

	for _, tt := range tests {
		t.Run("Edit one KindProperty", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Put)
			if !a.NoError(err) {
				return
			}

			url := "/api/v1/kind_properties/" + fmt.Sprint(tt.Put.KindPropertyId)
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
func TestDeleteKindPropertiesKindPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceKindProperties := service.NewKindPropertyService()

	tests := []struct {
		Kp   *storage.KindProperty
		Want int
	}{
		{&storage.KindProperty{Name: "test" + helpers.RandStringRunes(5)}, 204},
	}

	for _, tt := range tests {
		t.Run("Delete one KindProperty", func(t *testing.T) {
			if !a.NoError(serviceKindProperties.Create(tt.Kp, nil)) {
				return
			}

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/kind_properties/"+fmt.Sprint(tt.Kp.KindPropertyId), nil)
			if !a.NoError(err) {
				return
			}

			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}

func TestGetProperties(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/api/v1/properties", nil)
	if !a.NoError(err) {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	a.Equal(200, w.Code)
}
func TestGetPropertiesPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceProperties := service.NewPropertyService()
	var elId uint64 = 0
	type My struct {
		PropertyId string
		Want       int
	}

	properties, err := serviceProperties.GetProperties("property_id desc")
	if !a.NoError(err) {
		return
	}

	if len(properties) > 0 {
		elId = properties[0].PropertyId
	}

	tests := []My{
		{fmt.Sprint(elId + uint64(1)), 404},
		{"test", 400},
	}

	if len(properties) > 0 {
		tests = append(tests, My{fmt.Sprint(elId), 200})
	}

	for _, tt := range tests {
		t.Run("GET PropertyId", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/properties/"+tt.PropertyId, nil)
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			a.Equal(tt.Want, w.Code)
		})
	}
}
func TestPostProperties(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()

	tests := []struct {
		Post request.PostProperty
		Want int
	}{
		{request.PostProperty{}, 400},
		{
			request.PostProperty{
				Title:          helpers.RandStringRunes(5),
				Name:           helpers.RandStringRunes(5),
				KindPropertyId: 1,
			},
			201,
		},
	}

	for _, tt := range tests {
		t.Run("POST Properties", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Post)
			if !a.NoError(err) {
				return
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/properties", bytes.NewBuffer(b))
			if !a.NoError(err) {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			a.Equal(tt.Want, w.Code, w.Body)

			if w.Code == 201 {
				p := new(storage.Property)
				if a.NoError(json.Unmarshal(w.Body.Bytes(), p)) {
					a.NoError(serviceProperties.Delete(p.PropertyId, serviceValuesProperties, nil))
				}
			}
		})
	}
}
func TestPutPropertiesPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceProperties := service.NewPropertyService()
	serviceValuesProperties := service.NewValuesPropertyService()
	pr := &storage.Property{
		Title:          helpers.RandStringRunes(5),
		Name:           helpers.RandStringRunes(5),
		KindPropertyId: 1,
	}

	noError := a.NoError(serviceProperties.Create(pr, nil))
	defer func() {
		assert.NoError(t, serviceProperties.Delete(pr.PropertyId, serviceValuesProperties, nil))
	}()
	if !noError {
		return
	}

	tests := []struct {
		Put  request.PutProperty
		Want int
	}{
		{
			request.PutProperty{
				PropertyId:     pr.PropertyId,
				Title:          helpers.RandStringRunes(5),
				Name:           helpers.RandStringRunes(5),
				KindPropertyId: 2,
			},
			200,
		},
		{
			request.PutProperty{
				PropertyId:     pr.PropertyId,
				Title:          "",
				Name:           helpers.RandStringRunes(5),
				KindPropertyId: 3,
			},
			400,
		},
		{
			request.PutProperty{
				PropertyId:     pr.PropertyId + 1,
				Title:          helpers.RandStringRunes(5),
				Name:           pr.Name,
				KindPropertyId: 4,
			},
			404,
		},
	}

	for _, tt := range tests {
		t.Run("Put Property", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Put)
			if !a.NoError(err) {
				return
			}

			url := "/api/v1/properties/" + fmt.Sprint(tt.Put.PropertyId)
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
func TestDeletePropertiesPropertyId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceProperties := service.NewPropertyService()

	tests := []struct {
		Pr   *storage.Property
		Want int
	}{
		{
			&storage.Property{
				Name:           helpers.RandStringRunes(5),
				KindPropertyId: 1,
			},
			204,
		},
	}

	for _, tt := range tests {
		t.Run("DELETE PropertyId", func(t *testing.T) {
			noError := a.NoError(serviceProperties.Create(tt.Pr, nil))
			if !noError {
				return
			}

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/properties/"+fmt.Sprint(tt.Pr.PropertyId), nil)
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
func TestGetCatsCatId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceCats := service.NewCatService()
	var elId uint64
	type My struct {
		CatId string
		Query string
		Want  int
	}

	cats, err := serviceCats.GetCats()
	if !a.NoError(err) {
		return
	}

	// возьмем самый последний
	for _, v := range cats {
		if elId < v.CatId {
			elId = v.CatId
		}
	}

	tests := []My{
		{fmt.Sprint(elId + 1), "", 404},
		{fmt.Sprint(elId + 2), "withPropsOnlyFiltered=true", 404},
		{fmt.Sprint(elId + 3), "withPropsOnlyFiltered=1", 404},
	}

	if len(cats) > 0 {
		tests = append(tests, My{fmt.Sprint(elId), "", 200})
		tests = append(tests, My{fmt.Sprint(elId), "withPropsOnlyFiltered=true", 200})
		tests = append(tests, My{fmt.Sprint(elId), "withPropsOnlyFiltered=1", 200})
	}

	for _, tt := range tests {
		t.Run("Get one cat", func(t *testing.T) {
			w := httptest.NewRecorder()
			query := ""

			if tt.Query != "" {
				query += "?" + query
			}

			req, err := http.NewRequest(http.MethodGet, "/api/v1/cats/"+tt.CatId+query, nil)
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
		{request.PostCat{Name: helpers.RandStringRunes(5)}, 201},
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
					a.NoError(serviceCats.Delete(cat.CatId, nil))
				}
			}
		})
	}
}
func TestPutCatsCatId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceCats := service.NewCatService()
	cat := &storage.Cat{
		Name: helpers.RandStringRunes(3),
	}

	noError := a.NoError(serviceCats.Create(cat, nil))
	defer func() {
		a.NoError(serviceCats.Delete(cat.CatId, nil))
	}()
	if !noError {
		return
	}

	tests := []struct {
		Put  request.PutCat
		Want int
	}{
		{request.PutCat{CatId: cat.CatId, Name: helpers.RandStringRunes(5)}, 200},
		{request.PutCat{CatId: cat.CatId + 1, Name: cat.Name}, 404},
	}

	for _, tt := range tests {
		t.Run("PUT cat", func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.Put)
			if !a.NoError(err) {
				return
			}

			url := "/api/v1/cats/" + fmt.Sprint(tt.Put.CatId)
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
func TestDeleteCatsCatId(t *testing.T) {
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

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/cats/"+fmt.Sprint(tt.Cat.CatId), nil)
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
	var catId uint64

	cats, err := serviceCats.GetCats()
	if !a.NoError(err) {
		return
	}

	// возьмем самый последний
	for _, v := range cats {
		if catId < v.CatId {
			catId = v.CatId
		}
	}

	tests := []struct {
		CatId string
		Want  int
	}{
		{"catId=" + fmt.Sprint(catId), 200},
		{"", 200},
		{"catId=test", 500},
	}

	for _, tt := range tests {
		t.Run("GET ads", func(t *testing.T) {
			w := httptest.NewRecorder()
			query := ""

			if tt.CatId != "" {
				query += "?" + tt.CatId
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
func TestGetAdsAdId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceAds := service.NewAdService()
	var adId uint64
	type My struct {
		AdId string
		Want int
	}

	ads, err := serviceAds.GetAds(true)
	if !a.NoError(err) {
		return
	}

	if len(ads) > 0 {
		adId = ads[0].AdId
	}

	tests := []My{
		{fmt.Sprint(adId + 1), 404},
	}

	if len(ads) > 0 {
		tests = append(tests, My{fmt.Sprint(adId), 200})
	}

	for _, tt := range tests {
		t.Run("Get ad one", func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/ads/"+tt.AdId, nil)
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
	serviceImages := service.NewImageService()
	serviceAdDetail := service.NewAdDetailService()

	tests := []struct {
		Post request.PostAd
		Want int
	}{
		{request.PostAd{
			Title:       helpers.RandStringRunes(10),
			CatId:       1,
			Description: helpers.RandStringRunes(5)}, 201},
		{request.PostAd{}, 400},
	}

	for _, tt := range tests {
		t.Run("POST ad", func(t *testing.T) {
			body := new(bytes.Buffer)

			multiPartWriter := multipart.NewWriter(body)
			_ = multiPartWriter.WriteField("title", fmt.Sprint(tt.Post.Title))
			_ = multiPartWriter.WriteField("catId", fmt.Sprint(tt.Post.CatId))
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
					a.NoError(serviceAds.Delete(ad.AdId, nil, serviceImages, serviceAdDetail))
				}
			}
		})
	}
}
func TestPutAdsAdId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := assert.New(t)
	r := setupRouter()
	serviceAds := service.NewAdService()
	serviceImages := service.NewImageService()
	serviceAdDetail := service.NewAdDetailService()
	ad := &storage.Ad{
		Title:       helpers.RandStringRunes(10),
		CatId:       1,
		Description: helpers.RandStringRunes(3),
	}

	noError := a.NoError(serviceAds.Create(ad, nil))
	defer func() {
		a.NoError(serviceAds.Delete(ad.AdId, nil, serviceImages, serviceAdDetail))
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
				AdId: ad.AdId,
				PostAd: request.PostAd{
					Title:       helpers.RandStringRunes(10),
					CatId:       2,
					Description: helpers.RandStringRunes(5),
				},
			}, 200,
		},
		{
			request.PutAd{
				AdId: ad.AdId + 1,
				PostAd: request.PostAd{
					Title:       helpers.RandStringRunes(10),
					CatId:       2,
					Description: helpers.RandStringRunes(5),
				},
			}, 404,
		},
		{request.PutAd{}, 400},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			body := new(bytes.Buffer)
			sAdId := fmt.Sprint(tt.Put.AdId)

			multiPartWriter := multipart.NewWriter(body)
			_ = multiPartWriter.WriteField("adId", sAdId)
			_ = multiPartWriter.WriteField("title", tt.Put.Title)
			_ = multiPartWriter.WriteField("catId", fmt.Sprint(tt.Put.CatId))
			_ = multiPartWriter.WriteField("description", tt.Put.Description)
			_ = multiPartWriter.Close()

			req, err := http.NewRequest(http.MethodPut, "/api/v1/ads/"+sAdId, body)
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
func TestDeleteAdsAdId(t *testing.T) {
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
				Title:       helpers.RandStringRunes(10),
				CatId:       1,
				Description: helpers.RandStringRunes(5),
			}, 204,
		},
	}

	for _, tt := range tests {
		t.Run("Delete ad", func(t *testing.T) {
			if !a.NoError(serviceAds.Create(tt.Ad, nil)) {
				return
			}

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, "/api/v1/ads/"+fmt.Sprint(tt.Ad.AdId), nil)
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
		{helpers.RandStringRunes(10), 200},
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
