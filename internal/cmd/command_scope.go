package cmd

type commandScope string

const (
	commandScopeContextAware commandScope = "context-aware"
	commandScopeLocalOnly    commandScope = "local-only"
	commandScopeStatic       commandScope = "static"
)

func scopeForCommand(name string) commandScope {
	switch name {
	case "start", "enqueue", "status", "history", "stop", "retry", "restart", "dequeue":
		return commandScopeContextAware
	case "version", "schema", "xlsx":
		return commandScopeStatic
	case "context":
		return commandScopeStatic
	default:
		return commandScopeLocalOnly
	}
}
