package test

import (
	"bytes"
	"testing"

	"pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestInitEnvAndGetEnv(t *testing.T) {
	jsonData := `
{
	"test": {
		"value": "foobar",
		"empty": null
	}
}
`
	props := utils.EnvOption{
		ConfigData: bytes.NewBufferString(jsonData),
		ConfigType: []string{"json"},
	}
	utils.InitTestViper(props)

	// GetEnv应能拿到值
	val := utils.GetEnv("test.value")
	assert.Equal(t, "foobar", val)

	// GetEnv拿不到值时应该返回空字符串
	emptyVal := utils.GetEnv("test.empty")
	assert.Equal(t, "", emptyVal)

	// GetEnv拿不到值且有默认值时应该返回默认值
	notExistKey := "test.not_exist"
	valWithDefault := utils.GetEnvWithDefault(notExistKey, "defaultbar")
	assert.Equal(t, "defaultbar", valWithDefault)
}
