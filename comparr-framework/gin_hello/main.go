package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	jsoniter "github.com/json-iterator/go"
)

type UserInput struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address struct {
		City string `json:"city"`
	} `json:"address"`
}

type TestApi struct{}

func (t *TestApi) DoSayHelloInVn(name string) string {
	return "xin chào, " + name
}

func (t *TestApi) DoSayHelloInEn(name string) string {
	return "hello " + name
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// custom binding
type jsonIterBinding struct{}

func (jsonIterBinding) Name() string {
	return "json-iter"
}

func (jsonIterBinding) Bind(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return errors.New("missing request body")
	}
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(obj)
}

// BindBody mới: nhận []byte
func (jsonIterBinding) BindBody(body []byte, obj interface{}) error {
	return json.Unmarshal(body, obj)
}

// override default JSON binding
func init() {
	binding.JSON = jsonIterBinding{} //<-- loi
	/*
		cannot use jsonIterBinding{} (value of struct type jsonIterBinding) as binding.BindingBody value in assignment: jsonIterBinding does not implement binding.BindingBody (missing method BindBody)compilerInvalidIfaceAssign
	*/
}
func CreateUser(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	// giả lập xử lý
	msg := fmt.Sprintf("User %s, %d tuổi, sống ở %s",
		input.Name, input.Age, input.Address.City)

	c.JSON(http.StatusOK, gin.H{"msg": msg})
}

func (t *TestApi) Hello(c *gin.Context) {
	name := c.Param("name")
	lang := c.Param("langCode")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if lang == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "langCode is required"})
		return
	}

	var msg string
	switch lang {
	case "vn":
		msg = t.DoSayHelloInVn(name)
	case "en":
		msg = t.DoSayHelloInEn(name)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is not supported", lang)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": msg})
}

func main() {
	r := gin.Default()
	testApi := &TestApi{}

	r.POST("/user", CreateUser)
	r.GET("/hello/:name/:langCode", testApi.Hello)

	r.Run("0.0.0.0:8081") // listen and serve on 0.0.0.0:8080
}
