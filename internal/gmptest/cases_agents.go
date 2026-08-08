package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// ModifyAgents.Authorized/MinInterval/HeartbeatInterval are *int: unset,
	// a nonzero value, and an explicit zero (gvmd documents
	// HeartbeatInterval: 0 as a valid example) all round-trip distinctly.
	RegisterRequestCase("modify_agents", "unset", &commands.ModifyAgents{
		Agents: []commands.Agent{{ID: "fb6451bf-ec5a-45a8-8bab-5cf4b862e51b"}},
	})

	RegisterRequestCase("modify_agents", "fields_set", &commands.ModifyAgents{
		Agents:            []commands.Agent{{ID: "fb6451bf-ec5a-45a8-8bab-5cf4b862e51b"}},
		Authorized:        commands.Int(1),
		MinInterval:       commands.Int(1000),
		HeartbeatInterval: commands.Int(300),
	})

	RegisterRequestCase("modify_agents", "fields_zero", &commands.ModifyAgents{
		Agents:            []commands.Agent{{ID: "fb6451bf-ec5a-45a8-8bab-5cf4b862e51b"}},
		Authorized:        commands.Int(0),
		MinInterval:       commands.Int(0),
		HeartbeatInterval: commands.Int(0),
	})
}
