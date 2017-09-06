package cli

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/motki/motkid/cli/text"
	"github.com/motki/motkid/evedb"
	"github.com/motki/motkid/log"
	"github.com/peterh/liner"
	"github.com/shopspring/decimal"
)

type Prompter struct {
	logger log.Logger
	cli    *Server
	evedb  *evedb.EveDB
}

func NewPrompter(cli *Server, edb *evedb.EveDB, logger log.Logger) *Prompter {
	return &Prompter{
		logger: logger,
		cli:    cli,
		evedb:  edb,
	}
}

// PromptInt prompts the user for a valid integer input.
//
// If defVal is not nil, the prompt will be pre-populated with the default
// value.
//
// Additionally, any number of filter funcs can be passed in to perform
// custom validation and filtering on the user's input. Filter functions
// receive the value received from the prompt and return the new value and
// and indicator whether the value is valid.
func (p *Prompter) PromptInt(prompt string, defVal *int, filters ...func(int) (int, bool)) (int, bool) {
	var val int
	var valStr string
	var err error
	prompt = fmt.Sprintf("%s: ", prompt)
	for {
	begin:
		if defVal != nil {
			valStr, err = p.cli.PromptWithSuggestion(prompt, strconv.Itoa(*defVal), -1)
		} else {
			valStr, err = p.cli.Prompt(prompt)
		}
		if err != nil {
			if err == liner.ErrPromptAborted {
				return 0, false
			}
			if err == io.EOF {
				err = errors.New("unexpected EOF")
				fmt.Println()
			}
			p.logger.Debugf("unable to read input: %s", err.Error())
			goto begin
		}
		val, err = strconv.Atoi(valStr)
		if err != nil {
			fmt.Println("Invalid value, specify an integer.")
			goto begin
		}
		var ok bool
		for _, v := range filters {
			val, ok = v(val)
			if !ok {
				goto begin
			}
		}
		break
	}
	return val, true
}

// PromptString prompts the user for any string input.
//
// If defVal is not nil, the prompt will be pre-populated with the default
// value.
//
// Additionally, any number of filter funcs can be passed in to perform
// custom validation and filtering on the user's input. Filter functions
// receive the value received from the prompt and return the new value and
// and indicator whether the value is valid.
func (p *Prompter) PromptString(prompt string, defVal *string, filters ...func(string) (string, bool)) (string, bool) {
	var val string
	var err error
	prompt = fmt.Sprintf("%s: ", prompt)
	for {
	begin:
		if defVal != nil {
			val, err = p.cli.PromptWithSuggestion(prompt, *defVal, -1)
		} else {
			val, err = p.cli.Prompt(prompt)
		}
		if err != nil {
			if err == liner.ErrPromptAborted {
				return "", false
			}
			if err == io.EOF {
				err = errors.New("unexpected EOF")
				fmt.Println()
			}
			p.logger.Debugf("unable to read input: %s", err.Error())
			goto begin
		}
		var ok bool
		for _, v := range filters {
			val, ok = v(val)
			if !ok {
				goto begin
			}
		}
		break
	}
	return val, true
}

// PromptStringWithArgs prompts the user for any string input.
//
// This function differs from PromptString in that it will split the received
// input by spaces, returning the first value as the main value, and a slice
// of additional arguments.
//
// If defVal is not nil, the prompt will be pre-populated with the default
// value.
//
// Additionally, any number of filter funcs can be passed in to perform
// custom validation and filtering on the user's input. Filter functions
// receive the value received from the prompt and return the new value and
// and indicator whether the value is valid.
func (p *Prompter) PromptStringWithArgs(prompt string, defVal *string, filters ...func(string) (string, bool)) (string, []string, bool) {
	var val string
	var args []string
	var err error
	prompt = fmt.Sprintf("%s: ", prompt)
	for {
	begin:
		if defVal != nil {
			val, err = p.cli.PromptWithSuggestion(prompt, *defVal, -1)
		} else {
			val, err = p.cli.Prompt(prompt)
		}
		if err != nil {
			if err == liner.ErrPromptAborted {
				return "", nil, false
			}
			if err == io.EOF {
				err = errors.New("unexpected EOF")
				fmt.Println()
			}
			p.logger.Debugf("unable to read input: %s", err.Error())
			goto begin
		}
		parts := strings.Split(val, " ")
		val = parts[0]
		args = parts[1:]
		var ok bool
		for _, v := range filters {
			val, ok = v(val)
			if !ok {
				goto begin
			}
		}
		break
	}
	return val, args, true
}

// PromptDecimal prompts the user for a valid decimal input.
//
// If defVal is not nil, the prompt will be pre-populated with the default
// value.
//
// Additionally, any number of filter funcs can be passed in to perform
// custom validation and filtering on the user's input. Filter functions
// receive the value received from the prompt and return the new value and
// and indicator whether the value is valid.
func (p *Prompter) PromptDecimal(prompt string, defVal *decimal.Decimal, filters ...func(decimal.Decimal) (decimal.Decimal, bool)) (decimal.Decimal, bool) {
	var val decimal.Decimal
	var valStr string
	var defStr string
	if defVal != nil {
		defStr = defVal.StringFixed(2)
	}
	var err error
	prompt = fmt.Sprintf("%s: ", prompt)
	for {
	begin:
		if defVal != nil {
			valStr, err = p.cli.PromptWithSuggestion(prompt, defStr, -1)
		} else {
			valStr, err = p.cli.Prompt(prompt)
		}
		if err != nil {
			if err == liner.ErrPromptAborted {
				return decimal.Zero, false
			}
			if err == io.EOF {
				err = errors.New("unexpected EOF")
				fmt.Println()
			}
			p.logger.Debugf("unable to read input: %s", err.Error())
			goto begin
		}
		val, err = decimal.NewFromString(valStr)
		if err != nil {
			fmt.Println("Invalid input, specify a valid decimal value.")
			goto begin
		}
		var ok bool
		for _, v := range filters {
			val, ok = v(val)
			if !ok {
				goto begin
			}
		}
		break
	}
	return val, true
}

// PromptRegion prompts the user for a valid region input.
//
// If the user enters an integer, it is treated as the region's Region ID.
// Otherwise, the value is used to lookup regions.
//
// This function also accepts an initial input that should be used to
// as the first round of prompt input.
func (p *Prompter) PromptRegion(prompt string, initialInput string) (*evedb.Region, bool) {
	var val *evedb.Region
	var id int
	var err error
	var regions []*evedb.Region
	valStr := initialInput
	prompt = fmt.Sprintf("%s: ", prompt)
	regions, err = p.evedb.GetAllRegions()
	if err != nil {
		p.logger.Warnf("error loading regions: %s", err.Error())
		fmt.Println("Error loading regions, try again.")
		return nil, false
	}
	regionIndex := map[string]*evedb.Region{}
	for _, region := range regions {
		regionIndex[strings.ToUpper(region.Name)] = region
	}
	for {
		// This loop is ordered in such a way that it does the input validation
		// first. This allows us to specify an initial input and test that first,
		// before actually prompting the user for input.
		if valStr == "" {
			goto prompt
		}
		id, err = strconv.Atoi(valStr)
		if err != nil {
			matches := []*evedb.Region{}
			checkVal := strings.ToUpper(valStr)
			for capsName, region := range regionIndex {
				if strings.Contains(capsName, checkVal) {
					matches = append(matches, region)
				}
			}
			fmt.Printf("Top %d results for \"%s\":\n", len(matches), valStr)
			for _, r := range matches {
				fmt.Printf("%s  %s\n", text.PadTextLeft(fmt.Sprintf("%d", r.RegionID), 12), r.Name)
			}
			goto prompt
		}
		val, err = p.evedb.GetRegion(id)
		if err != nil {
			fmt.Printf("No region exists with ID %d.\n", id)
			goto prompt
		}
		// We have a valid value, break out of the loop.
		break

	prompt:
		valStr, err = p.cli.Prompt(prompt)
		if err != nil {
			if err == liner.ErrPromptAborted {
				return nil, false
			}
			if err == io.EOF {
				err = errors.New("unexpected EOF")
				fmt.Println()
			}
			p.logger.Debugf("unable to read input: %s", err.Error())
			goto prompt
		}
		// Loop back around and check valStr for valid input.
		continue
	}
	return val, true
}

// PromptItemTypeDetail prompts the user for a valid item type input.
//
// If the user enters an integer, it is treated as the item's Type ID.
// Otherwise, the value is used to lookup item types.
//
// This function also accepts an initial input that should be used to
// as the first round of prompt input.
func (p *Prompter) PromptItemTypeDetail(prompt string, initialInput string) (*evedb.ItemTypeDetail, bool) {
	var val *evedb.ItemTypeDetail
	var id int
	var err error
	valStr := initialInput
	prompt = fmt.Sprintf("%s: ", prompt)
	for {
		// This loop is ordered in such a way that it does the input validation
		// first. This allows us to specify an initial input and test that first,
		// before actually prompting the user for input.
		if valStr == "" {
			goto prompt
		}
		id, err = strconv.Atoi(valStr)
		if err != nil {
			its, err := p.evedb.QueryItemTypes(valStr)
			if err != nil || len(its) == 0 {
				if err != nil {
					p.logger.Debugf("error querying item types: %s", err.Error())
				}
				fmt.Printf("Nothing found for \"%s\".\n", valStr)
				goto prompt
			}
			fmt.Printf("Top %d results for \"%s\":\n", len(its), valStr)
			for _, it := range its {
				fmt.Printf("%s  %s\n", text.PadTextLeft(fmt.Sprintf("%d", it.ID), 8), it.Name)
			}
			goto prompt
		}
		val, err = p.evedb.GetItemTypeDetail(id)
		if err != nil {
			fmt.Printf("No item exists with ID %d.\n", id)
			p.logger.Warnf("error fetching item type detail: %s", err.Error())
			goto prompt
		}
		// We have a valid value, break out of the loop.
		break

	prompt:
		valStr, err = p.cli.Prompt(prompt)
		if err != nil {
			if err == liner.ErrPromptAborted {
				return nil, false
			}
			if err == io.EOF {
				err = errors.New("unexpected EOF")
				fmt.Println()
			}
			p.logger.Debugf("unable to read input: %s", err.Error())
			goto prompt
		}
		// Loop back around and check valStr for valid input.
		continue
	}
	return val, true
}
