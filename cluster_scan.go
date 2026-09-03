package valkey

import (
	"strconv"
	"strings"

	"github.com/valkey-io/valkey-go/internal/cmds"
)

const clusterScanCommand = "CLUSTERSCAN"

func clusterScanSlot(cmd Completed) (uint16, bool) {
	if cmd.IsEmpty() {
		return cmds.InitSlot, false
	}
	return clusterScanArgsSlot(cmd.Commands())
}

func clusterScanArgsSlot(args []string) (uint16, bool) {
	if len(args) < 2 || !strings.EqualFold(args[0], clusterScanCommand) {
		return cmds.InitSlot, false
	}

	cursor := args[1]

	if cursor == "0" {
		return clusterScanOptionSlot(args[2:])
	}

	start := strings.IndexByte(cursor, '{')
	if start == -1 {
		return cmds.InitSlot, false
	}

	end := strings.IndexByte(cursor[start+1:], '}')
	if end == -1 {
		return cmds.InitSlot, false
	}

	end += start + 1

	if end == start+1 {
		return cmds.InitSlot, false
	}

	return cmds.Slot(cursor[start : end+1]), true
}

func clusterScanOptionSlot(args []string) (uint16, bool) {
	for i := 0; i < len(args); {
		switch strings.ToUpper(args[i]) {
		case "MATCH", "COUNT", "TYPE":
			i += 2
		case "SLOT":
			if i+1 < len(args) {
				slot, err := strconv.ParseUint(args[i+1], 10, 16)
				if err == nil && slot < 16384 {
					return uint16(slot), true
				}
			}
			return cmds.InitSlot, false
		default:
			i++
		}
	}

	return cmds.InitSlot, false
}
