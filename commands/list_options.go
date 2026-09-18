package commands

import (
	"dvc/models"
	"fmt"
	"strconv"
	"strings"
)

func parseOptions(args []string) (models.ListOptions, error) {
	options := models.QueryOptions{}
	format := models.OutputFormatTable
	var rawFormat string

	parseStringFlag := func(target *string, flag string, i *int) error {
		if *target != "" {
			return fmt.Errorf("duplicate %s", flag)
		}
		if *i+1 >= len(args) {
			return fmt.Errorf("missing value for %s", flag)
		}
		*i++
		*target = args[*i]
		return nil
	}

	parseIntFlag := func(target *int, flag string, i *int) error {
		var valStr string
		if err := parseStringFlag(&valStr, flag, i); err != nil {
			return err
		}
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %w", flag, err)
		}
		if val <= 0 {
			return fmt.Errorf("%s must be greater than 0", flag)
		}
		*target = val
		return nil
	}

	for i := 0; i < len(args); i++ {
		var err error
		switch args[i] {
		case "--select":
			err = parseStringFlag(&options.Select, args[i], &i)
		case "--filter":
			err = parseStringFlag(&options.Filter, args[i], &i)
		case "--order-by":
			err = parseStringFlag(&options.OrderBy, args[i], &i)
		case "--expand":
			err = parseStringFlag(&options.Expand, args[i], &i)
		case "--top":
			err = parseIntFlag(&options.Top, args[i], &i)
		case "--max-page-size":
			err = parseIntFlag(&options.MaxPageSize, args[i], &i)
		case "--output", "-o":
			if err = parseStringFlag(&rawFormat, args[i], &i); err != nil {
				break
			}
			switch strings.ToLower(rawFormat) {
			case string(models.OutputFormatJSON):
				format = models.OutputFormatJSON
			case string(models.OutputFormatTable):
				format = models.OutputFormatTable
			default:
				return models.ListOptions{}, fmt.Errorf("invalid output format %q (allowed: json, table)", rawFormat)
			}
		default:
			return models.ListOptions{}, fmt.Errorf("unknown option: %s", args[i])
		}

		if err != nil {
			return models.ListOptions{}, err
		}
	}

	return models.ListOptions{Query: options, Output: format}, nil
}
