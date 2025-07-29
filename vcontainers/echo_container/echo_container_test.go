package echo_container

import "testing"

const configPath = "./../../vapi/cmd/config.yaml"

func TestEchoContainer(t *testing.T) {
	c := NewEchoContainer(configPath)

	// (&core).Cache.Get().Set(t.Context(), "test", "test", 0)
	c.Core.Get().Cache.Get().Set(t.Context(), "test", "test", 0)
}
