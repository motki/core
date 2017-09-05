package command

import (
	"fmt"
	"strings"

	"github.com/motki/motkid/cli"
	"github.com/motki/motkid/cli/text"
)

// EVETypesCommand provides item type lookup and display.
type EVETypesCommand struct {
	env *cli.Prompter
}

func NewEVETypesCommand(prompter *cli.Prompter) EVETypesCommand {
	return EVETypesCommand{prompter}
}

func (c EVETypesCommand) Prefixes() []string {
	return []string{"types", "type"}
}

func (c EVETypesCommand) Description() string {
	return "Display type information from the EVE database."
}

func (c EVETypesCommand) Handle(subcmd string, args ...string) {
	switch {
	case len(subcmd) == 0:
		c.PrintHelp()

	case subcmd == "search" || subcmd == "show":
		item, ok := c.env.PromptItemTypeDetail("Specify Item Type", strings.Join(args, " "))
		if !ok {
			return
		}

		fmt.Println(text.CenterText(fmt.Sprintf("%s (ID: %d)", item.Name, item.ID), text.StandardTerminalWidthInChars))
		fmt.Println()
		fmt.Println(text.WrapText(strings.Replace(item.Description, "<br>", "\n", -1), text.StandardTerminalWidthInChars))
		fmt.Println()
		fmt.Printf("%s%s\n", text.PadTextRight("Type", 15), "Value")
		fmt.Printf("%s%s\n", text.PadTextRight("Volume", 15), item.Volume.StringFixed(2))
		fmt.Printf("%s%s\n", text.PadTextRight("Mass", 15), item.Mass.StringFixed(2))
		fmt.Printf("%s%s\n", text.PadTextRight("Capacity", 15), item.Capacity.StringFixed(2))
		fmt.Printf("%s%s\n", text.PadTextRight("Base Price", 15), item.BasePrice.StringFixed(2))
		fmt.Println()
		fmt.Printf("%s%d\n", text.PadTextRight("Parent Type ID", 15), item.ParentTypeID)
		fmt.Printf("%s%d\n", text.PadTextRight("Blueprint ID", 15), item.BlueprintID)
		fmt.Printf("%s%s (ID: %d)\n", text.PadTextRight("Group", 15), item.GroupName, item.GroupID)
		fmt.Printf("%s%s (ID: %d)\n", text.PadTextRight("Category", 15), item.CategoryName, item.CategoryID)

	default:
		fmt.Println("Unknown subcommand:", subcmd)
		c.PrintHelp()
	}
}

func (c EVETypesCommand) PrintHelp() {
	fmt.Println(text.WrapText(`Command "types" can be used to search for and display information
about items that exist in the EVE universe.`, text.StandardTerminalWidthInChars))
	fmt.Printf(`
Subcommands:
  %s Search for an item matching the given query.
  %s Display details for a given item type.
`, text.PadTextRight("search [query]", 15), text.PadTextRight("show [typeID]", 15))
	fmt.Println()
	fmt.Println(text.WrapText(`Note that the commands are identical, and either one accepts an optional integer or string argument. If an integer is given as an argument or in the prompt, the command will attempt to load the Item Type with the given Type ID. Likewise, if a string is given as an argument or in the prompt, the command will show results matching the input.`, text.StandardTerminalWidthInChars))
}
