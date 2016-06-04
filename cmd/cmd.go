package cmd

//all the commands have to register with their run func
import (
	"fmt"
	"github.com/redis/cache"
	"strings"
)

var (
	commands map[string]RunFunc
)

func init() {
	commands = make(map[string]RunFunc)
}

//Holds command name and the options
type Command struct {
	Name    string
	Options []string
}

//Constructs command obj from req
func NewCommand(req string) (Command, error) {
	fields := strings.Fields(req)

	if len(fields) == 0 {
		return Command{}, fmt.Errorf("Invalid command")
	}

	command := Command{}
	command.Name = strings.ToUpper(fields[0])
	command.Options = fields[1:]

	for i, option := range command.Options {
		command.Options[i] = strings.Replace(option, `"`, "", -1)
	}

	return command, nil
}

type RunFunc func(options []string, cache *cache.Cache) ([]byte, bool)

//Handler for command providers to register their service
func Register(cmdName string, runFunc RunFunc) error {
	if _, ok := commands[cmdName]; ok {
		return fmt.Errorf("command with name %s already registered", cmdName)
	}

	commands[cmdName] = runFunc
	return nil
}

//Identifies the runfunc for the command name and executes the func
func ExecuteCmd(cmd Command, cache *cache.Cache) ([]byte, bool) {
	runFunc, ok := commands[cmd.Name]
	if !ok {
		return GetErrMsg(fmt.Sprintf("ERR unknown command '%s'", cmd.Name)), false
	}

	return runFunc(cmd.Options, cache)
}
