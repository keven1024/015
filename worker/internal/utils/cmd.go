package utils

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"go.uber.org/zap"
)

func RunCommand(bin string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		zap.L().Error("command failed",
			zap.String("bin", bin),
			zap.Strings("args", args),
			zap.ByteString("output", bytes),
			zap.Error(err),
		)
		return nil, fmt.Errorf("command failed: %w", err)
	}
	return bytes, nil
}
