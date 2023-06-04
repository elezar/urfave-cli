package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSuggestFlag(t *testing.T) {
	// Given
	app := testFishCommand()

	for _, testCase := range []struct {
		provided, expected string
	}{
		{"", ""},
		{"a", "--another-flag"},
		{"hlp", "--help"},
		{"k", ""},
		{"s", "-s"},
	} {
		// When
		res := suggestFlag(app.Flags, testCase.provided, false)

		// Then
		expect(t, res, testCase.expected)
	}
}

func TestSuggestFlagHideHelp(t *testing.T) {
	// Given
	app := testFishCommand()

	// When
	res := suggestFlag(app.Flags, "hlp", true)

	// Then
	expect(t, res, "--fl")
}

func TestSuggestFlagFromError(t *testing.T) {
	// Given
	app := testFishCommand()

	for _, testCase := range []struct {
		command, provided, expected string
	}{
		{"", "hel", "--help"},
		{"", "soccer", "--socket"},
		{"config", "anot", "--another-flag"},
	} {
		// When
		res, _ := app.suggestFlagFromError(
			errors.New(providedButNotDefinedErrMsg+testCase.provided),
			testCase.command,
		)

		// Then
		expect(t, res, fmt.Sprintf(SuggestDidYouMeanTemplate+"\n\n", testCase.expected))
	}
}

func TestSuggestFlagFromErrorWrongError(t *testing.T) {
	// Given
	app := testFishCommand()

	// When
	_, err := app.suggestFlagFromError(errors.New("invalid"), "")

	// Then
	expect(t, true, err != nil)
}

func TestSuggestFlagFromErrorWrongCommand(t *testing.T) {
	// Given
	app := testFishCommand()

	// When
	_, err := app.suggestFlagFromError(
		errors.New(providedButNotDefinedErrMsg+"flag"),
		"invalid",
	)

	// Then
	expect(t, true, err != nil)
}

func TestSuggestFlagFromErrorNoSuggestion(t *testing.T) {
	// Given
	app := testFishCommand()

	// When
	_, err := app.suggestFlagFromError(
		errors.New(providedButNotDefinedErrMsg+""),
		"",
	)

	// Then
	expect(t, true, err != nil)
}

func TestSuggestCommand(t *testing.T) {
	// Given
	app := testFishCommand()

	for _, testCase := range []struct {
		provided, expected string
	}{
		{"", ""},
		{"conf", "config"},
		{"i", "i"},
		{"information", "info"},
		{"inf", "info"},
		{"con", "config"},
		{"not-existing", "info"},
	} {
		// When
		res := suggestCommand(app.Commands, testCase.provided)

		// Then
		expect(t, res, testCase.expected)
	}
}

func ExampleCommand_Suggest() {
	cmd := &Command{
		Name:                  "greet",
		ErrWriter:             os.Stdout,
		Suggest:               true,
		HideHelp:              false,
		HideHelpCommand:       true,
		CustomAppHelpTemplate: "(this space intentionally left blank)\n",
		Flags: []Flag{
			&StringFlag{Name: "name", Value: "squirrel", Usage: "a name to say"},
		},
		Action: func(cCtx *Context) error {
			fmt.Printf("Hello %v\n", cCtx.String("name"))
			return nil
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if cmd.Run(ctx, []string{"greet", "--nema", "chipmunk"}) == nil {
		fmt.Println("Expected error")
	}
	// Output:
	// Incorrect Usage: flag provided but not defined: -nema
	//
	// Did you mean "--name"?
	//
	// (this space intentionally left blank)
}

func ExampleCommand_Suggest_command() {
	cmd := &Command{
		Name:                  "greet",
		ErrWriter:             os.Stdout,
		Suggest:               true,
		HideHelpCommand:       true,
		CustomAppHelpTemplate: "(this space intentionally left blank)\n",
		Flags: []Flag{
			&StringFlag{Name: "name", Value: "squirrel", Usage: "a name to say"},
		},
		Action: func(cCtx *Context) error {
			fmt.Printf("Hello %v\n", cCtx.String("name"))
			return nil
		},
		Commands: []*Command{
			{
				Name:               "neighbors",
				HideHelp:           false,
				CustomHelpTemplate: "(this space intentionally left blank)\n",
				Flags: []Flag{
					&BoolFlag{Name: "smiling"},
				},
				Action: func(cCtx *Context) error {
					if cCtx.Bool("smiling") {
						fmt.Println("😀")
					}
					fmt.Println("Hello, neighbors")
					return nil
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if cmd.Run(ctx, []string{"greet", "neighbors", "--sliming"}) == nil {
		fmt.Println("Expected error")
	}
	// Output:
	// Incorrect Usage: flag provided but not defined: -sliming
	//
	// Did you mean "--smiling"?
	//
	// (this space intentionally left blank)
}
