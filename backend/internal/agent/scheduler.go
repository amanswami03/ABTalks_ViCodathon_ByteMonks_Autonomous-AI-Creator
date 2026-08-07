package agent

import (
	"log"
	"time"

	"ai-persona-agent/internal/llm"
	"ai-persona-agent/internal/store"
)

// StartScheduler kicks off the autonomous loop for one agent. This is
// what makes publishing happen "over time" without any further human
// input — call this once right after the agent is created in /init,
// and never call RunCycle manually again after that.
//
// TODO(copilot): tune the interval. 15-30 min is reasonable for a 48h
// evaluation window (gives ~100-190 cycles, plenty of chances to publish
// without spamming).
func StartScheduler(client *llm.Client, s *store.Store, agentID string, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Run once immediately so there's not a long silent gap after init.
		if err := RunCycle(client, s, agentID); err != nil {
			log.Printf("[agent %s] initial cycle error: %v", agentID, err)
		}

		for range ticker.C {
			if err := RunCycle(client, s, agentID); err != nil {
				log.Printf("[agent %s] cycle error: %v", agentID, err)
			}
		}
	}()
}
