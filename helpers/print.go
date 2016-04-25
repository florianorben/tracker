package helpers

import (
	"fmt"
	"github.com/spf13/viper"
)

func PrintBold(s string) string {
	if viper.GetBool("colors.enabled") == false {
		return s
	}

	return fmt.Sprintf("\033[1m%s\033[0m", s)
}

func PrintPurple(s string) string {
	if viper.GetBool("colors.enabled") == false {
		return s
	}

	return fmt.Sprintf("\033[35m%s\033[0m", s)
}

func PrintGreen(s string) string {
	if viper.GetBool("colors.enabled") == false {
		return s
	}

	return fmt.Sprintf("\033[32m%s\033[0m", s)
}

func PrintBlue(s string) string {
	if viper.GetBool("colors.enabled") == false {
		return s
	}

	return fmt.Sprintf("\033[34m%s\033[0m", s)
}

func PrintRed(s string) string {
	if viper.GetBool("colors.enabled") == false {
		return s
	}

	return fmt.Sprintf("\033[31m%s\033[0m", s)
}

func PrintTeal(s string) string {
	if viper.GetBool("colors.enabled") == false {
		return s
	}

	return fmt.Sprintf("\033[36m%s\033[0m", s)
}
