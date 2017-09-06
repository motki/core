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

		colWidth := 15
		fmt.Println(text.CenterText(fmt.Sprintf("%s (ID: %d)", item.Name, item.ID), text.StandardTerminalWidthInChars))
		fmt.Println()
		fmt.Println(text.WrapText(strings.Replace(item.Description, "<br>", "\n", -1), text.StandardTerminalWidthInChars))
		fmt.Println()
		fmt.Printf("%s%s\n", text.PadTextRight("Type", colWidth), "Value")
		fmt.Printf("%s%s\n", text.PadTextRight("Volume", colWidth), item.Volume.StringFixed(2))
		fmt.Printf("%s%s\n", text.PadTextRight("Mass", colWidth), item.Mass.StringFixed(2))
		fmt.Printf("%s%s\n", text.PadTextRight("Capacity", colWidth), item.Capacity.StringFixed(2))
		fmt.Printf("%s%s\n", text.PadTextRight("Base Price", colWidth), item.BasePrice.StringFixed(2))
		fmt.Println()
		fmt.Printf("%s%d\n", text.PadTextRight("Parent Type ID", colWidth), item.ParentTypeID)
		fmt.Printf("%s%v\n", text.PadTextRight("Derivative Type IDs", colWidth), item.DerivativeTypeIDs)
		fmt.Printf("%s%d\n", text.PadTextRight("Blueprint ID", colWidth), item.BlueprintID)
		fmt.Printf("%s%s (ID: %d)\n", text.PadTextRight("Group", colWidth), item.GroupName, item.GroupID)
		fmt.Printf("%s%s (ID: %d)\n", text.PadTextRight("Category", colWidth), item.CategoryName, item.CategoryID)

	default:
		fmt.Println("Unknown subcommand:", subcmd)
		c.PrintHelp()
	}
}

func (c EVETypesCommand) PrintHelp() {
	fmt.Println()
	fmt.Println(text.WrapText(fmt.Sprintf(`Command "%s" can be used to search for and display information
about items that exist in the EVE universe.`, text.Boldf("types")), text.StandardTerminalWidthInChars))
	fmt.Println()
	fmt.Printf(`Subcommands:
  %s Search for an item matching the given query.
  %s Display details for a given item type.`,
		text.Boldf(text.PadTextRight("search [query|typeID]", 25)),
		text.Boldf(text.PadTextRight("show [query|typeID]", 25)))
	fmt.Println()
	fmt.Println()
	fmt.Println(text.WrapText(`Note that the commands are identical, and either one accepts an optional integer or string argument. If an integer is given as an argument or in the prompt, the command will attempt to load the Item Type with the given Type ID. If a string is given as an argument or in the prompt, the command will show results matching the input.`, text.StandardTerminalWidthInChars))
	fmt.Println()
}
