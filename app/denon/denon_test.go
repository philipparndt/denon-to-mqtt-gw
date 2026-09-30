package denon

import (
	"reflect"
	"testing"
)

// The read loop reuses its buffer. A short read must not see the tail of an
// earlier, longer one: a stale PWON after PWSTANDBY kept the receiver "on"
// in MQTT after it went to standby (2026-09-28, 2026-09-29).
func TestShortReadAfterLongReadIgnoresStaleBytes(t *testing.T) {
	state = Message{}
	reply := make([]byte, 1024)
	var pending string

	read := func(data string) {
		n := copy(reply, data)
		var lines []string
		lines, pending = splitReplies(pending, reply[:n])
		for _, line := range lines {
			handleReply(line)
		}
	}

	read("PWON\rMVMAX 80\rMV58\rPWON\r")
	read("PWSTANDBY\r")

	if state.Power != "STANDBY" {
		t.Fatalf("power = %q, want STANDBY", state.Power)
	}
}

func TestSplitRepliesKeepsIncompleteReplyPending(t *testing.T) {
	lines, pending := splitReplies("", []byte("PWON\rMV5"))
	if !reflect.DeepEqual(lines, []string{"PWON"}) || pending != "MV5" {
		t.Fatalf("got %q, pending %q", lines, pending)
	}

	lines, pending = splitReplies(pending, []byte("85\r"))
	if !reflect.DeepEqual(lines, []string{"MV585"}) || pending != "" {
		t.Fatalf("got %q, pending %q", lines, pending)
	}
}

func TestVolumeReplies(t *testing.T) {
	state = Message{}
	handleReply("MVMAX 80")
	handleReply("MV585")

	if state.VolumeMax != 80 || state.Volume != 58.5 {
		t.Fatalf("volume = %v, max = %v", state.Volume, state.VolumeMax)
	}
}
