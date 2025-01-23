package term

import (
	"fmt"

	"github.com/iBug/ibugtool/pkg/util"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "term",
	Short: "Commands for managing terminal features",
	Args:  cobra.NoArgs,
	RunE:  util.ShowHelp,
}

var colorCmd = &cobra.Command{
	Use:   "color [-n|--length size]",
	Short: "Display a colored stripe to test terminal color support",
	Args:  cobra.NoArgs,
	RunE:  colorRunE,
}

var colorSize int

/* Original AWK code:
#!/usr/bin/awk -f

BEGIN {
    for (colnum = 0; colnum < 77; colnum++) {
        r = 255 - (colnum * 255 / 76);
        g = colnum * 510 / 76;
        b = 255 - r;
        if (g > 255)
            g = 510 - g;
        printf "\033[48;2;%d;%d;%dm", r, g, b;
        printf "\033[38;2;%d;%d;%dm", 255 - r, 255 - g, 255 - b;
        printf "%s\033[0m", substr("/\\", colnum % 2 + 1, 1);
    }
    printf "\n";
    exit;
}
*/

func colorRunE(cmd *cobra.Command, args []string) error {
	for i := 0; i < colorSize; i++ {
		r := 255 - (i * 255 / (colorSize - 1))
		g := i * 2 * 255 / (colorSize - 1)
		b := 255 - r
		if g > 255 {
			g = 510 - g
		}
		fmt.Fprintf(cmd.OutOrStdout(), "\x1B[48;2;%d;%d;%dm", r, g, b)
		fmt.Fprintf(cmd.OutOrStdout(), "\x1B[38;2;%d;%d;%dm", 255-r, 255-g, 255-b)
		fmt.Fprintf(cmd.OutOrStdout(), "%c", "/\\"[i%2])
	}
	fmt.Fprintf(cmd.OutOrStdout(), "\x1B[0m\n")
	return nil
}

func init() {
	flags := colorCmd.Flags()
	flags.IntVarP(&colorSize, "length", "n", 77, "Length of color stripe")
	Cmd.AddCommand(colorCmd)
}
