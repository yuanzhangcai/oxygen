package tools

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	server   = "redis:6379"
	password = "12345678"
	prefix   = ""
)

func TestInitRedis(t *testing.T) {
	err := InitRedis("", password, prefix)
	assert.NotNil(t, err)

	err = InitRedis(server+"1", password, prefix)
	assert.NotNil(t, err)

	// 初次初始化
	err = InitRedis(server, password, prefix)
	assert.Nil(t, err)

	// 重复初始化
	err = InitRedis(server, password, prefix)
	assert.Nil(t, err)
}

func TestGetRedis(t *testing.T) {
	assert.Equal(t, client, GetRedis())
}

func TestCommand(t *testing.T) {
	_ = InitRedis(server, password, prefix)

	valStr := "SetObject"
	err := client.SetObject("test_set_string", valStr, 300*time.Second)
	assert.Nil(t, err)
	newValStr := ""
	err = client.GetObject("test_set_string", &newValStr)
	assert.Nil(t, err)
	assert.Equal(t, valStr, newValStr)

	valInt := 453
	err = client.SetObject("test_set_int", valInt, 300*time.Second)
	assert.Nil(t, err)
	newValInt := 0
	err = client.GetObject("test_set_int", &newValInt)
	assert.Nil(t, err)
	assert.Equal(t, valInt, newValInt)

	valMap := map[string]interface{}{"name": "zacyuan"}
	err = client.SetObject("test_set_map", valMap, 300*time.Second)
	assert.Nil(t, err)
	newValMap := make(map[string]interface{})
	err = client.GetObject("test_set_map", &newValMap)
	assert.Nil(t, err)
	assert.Equal(t, valMap["name"], newValMap["name"])

	type testStruct struct {
		Name string
		Age  int
	}
	valStc := testStruct{Name: "zacyuan", Age: 18}
	err = client.SetObject("test_set_struct", valStc, 300*time.Second)
	assert.Nil(t, err)
	newValStc := testStruct{}
	err = client.GetObject("test_set_struct", &newValStc)
	assert.Nil(t, err)
	assert.Equal(t, valStc.Name, newValStc.Name)
	assert.Equal(t, valStc.Age, newValStc.Age)
}
