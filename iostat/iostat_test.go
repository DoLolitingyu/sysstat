package iostat

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func Test_GetData(t *testing.T) {
	go func() {
		stat, _ := GetData()
		fmt.Println(stat)
	}()
	go func() {
		ctx := context.Background()
		cmd := `iostat -x 1 2 | grep -v dm | sed '0,/^Device/d' | tr -s " " | awk 'NF'`
		output, err := exec.CommandContext(ctx, "bash", "-c", cmd).CombinedOutput()
		if err != nil {
			fmt.Println(err)
			return
		}
		idx := 0
		strArr := strings.Split(string(output), "\n")
		for _, item := range strArr {
			idx++
			if strings.Contains(item, "Device:") {
				break
			}
		}

		newStrArr := strArr[idx:]
		fmt.Println(newStrArr)
	}()
	time.Sleep(2 * time.Second)
}

