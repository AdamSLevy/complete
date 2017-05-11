package complete

import "github.com/posener/complete/match"

// Command represents a command line
// It holds the data that enables auto completion of command line
// Command can also be a sub command.
type Command struct {
	// Sub is map of sub commands of the current command
	// The key refer to the sub command name, and the value is it's
	// Command descriptive struct.
	Sub Commands

	// Flags is a map of flags that the command accepts.
	// The key is the flag name, and the value is it's predictions.
	Flags Flags

	// Args are extra arguments that the command accepts, those who are
	// given without any flag before.
	Args Predictor
}

// Commands is the type of Sub member, it maps a command name to a command struct
type Commands map[string]Command

// Flags is the type Flags of the Flags member, it maps a flag name to the flag predictions.
type Flags map[string]Predictor

// Predict returns all possible predictions for args according to the command struct
func (c *Command) Predict(a Args) (predictions []string) {
	predictions, _ = c.predict(a)
	return
}

func (c *Command) predict(a Args) (options []string, only bool) {

	// if wordCompleted has something that needs to follow it,
	// it is the most relevant completion
	if predictor, ok := c.Flags[a.LastCompleted]; ok && predictor != nil {
		Log("Predicting according to flag %s", a.Last)
		return predictor.Predict(a), true
	}

	sub, options, only := c.searchSub(a)
	if only {
		return
	}

	// if no sub command was found, return a list of the sub commands
	if sub == "" {
		options = append(options, c.subCommands(a.Last)...)
	}

	// add global available complete Predict
	for flag := range c.Flags {
		if m := match.Prefix(flag); m.Match(a.Last) {
			options = append(options, m.String())
		}
	}

	// add additional expected argument of the command
	if c.Args != nil {
		options = append(options, c.Args.Predict(a)...)
	}

	return
}

// searchSub searches recursively within sub commands if the sub command appear
// in the on of the arguments.
func (c *Command) searchSub(a Args) (sub string, all []string, only bool) {
	for i, arg := range a.Completed {
		if cmd, ok := c.Sub[arg]; ok {
			sub = arg
			all, only = cmd.predict(a.from(i))
			return
		}
	}
	return
}

// subCommands returns a list of matching sub commands
func (c *Command) subCommands(last string) (prediction []string) {
	for sub := range c.Sub {
		if m := match.Prefix(sub); m.Match(last) {
			prediction = append(prediction, m.String())
		}
	}
	return
}
