package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePasswordHash(t *testing.T) {
	// 保存原始环境变量
	originalSalt := os.Getenv("PASSWORD_SALT")
	defer os.Setenv("PASSWORD_SALT", originalSalt)

	tests := []struct {
		name        string
		password    string
		salt        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "PASSWORD_SALT未配置",
			password:    "testpassword",
			salt:        "",
			expectError: true,
			errorMsg:    "请配置PASSWORD_SALT",
		},
		{
			name:        "正常生成哈希",
			password:    "testpassword123",
			salt:        "testsalt",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			if tt.salt != "" {
				os.Setenv("PASSWORD_SALT", tt.salt)
			} else {
				os.Unsetenv("PASSWORD_SALT")
			}

			hash, err := GeneratePasswordHash(tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误，但得到了 nil")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("期望错误信息 '%s'，但得到了 '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误，但得到了: %v", err)
				return
			}

			// 验证哈希格式：argon2 32字节 = 64个十六进制字符
			if len(hash) != 64 {
				t.Errorf("期望哈希长度为64，但得到了 %d", len(hash))
			}
			assert.Equal(t, hash, "537275f995fdd46eb2e5455b8a29adccb60c5637689d29646676a5f1bffb63f3")
		})
	}
}
