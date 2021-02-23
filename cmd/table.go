package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"got2s/util"
	"os"
)

var (
	u           string
	db          string
	t           string
	packageName string
	isJson      bool
	isImport    bool
	isFunc      bool
)

var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "Convert mysql table to struct",
	Run:   parseTable,
}

func init() {
	tableCmd.Flags().StringVarP(&u, "url", "u", "root:123456@tcp(localhost:3306)", "the url of sql connection")
	tableCmd.Flags().StringVarP(&db, "db", "d", "", "the database of sql")
	tableCmd.Flags().StringVarP(&t, "table", "t", "", "the table of sql")
	tableCmd.MarkFlagRequired("db")
	tableCmd.MarkFlagRequired("table")

	tableCmd.Flags().StringVarP(&packageName, "package", "p", "", "the name of struct package")
	tableCmd.Flags().BoolVarP(&isJson, "json", "j", false, "add json tag")
	tableCmd.Flags().BoolVarP(&isImport, "import", "i", false, "auto import needed package")
	tableCmd.Flags().BoolVarP(&isFunc, "func", "f", false, "add TableName function")

	rootCmd.AddCommand(tableCmd)
}

func parseTable(cmd *cobra.Command, args []string) {
	var isPackage bool
	if packageName != "" {
		isPackage = true
	}
	cfg := &util.Option{
		DataBase:   db,
		Table:      t,
		Url:        u,
		IsSQL:      false,
		IsJson:     isJson,
		FormatType: 0,
		IsImport:   isImport,
		IsFunc:     isFunc,
		IsPackage:  isPackage,
		Package:    packageName,
	}

	structStr, err := util.Parse(cfg)
	if err != nil {
		fmt.Printf("parse table to struct failed.\n\n%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(structStr)
	if err := clipboard.WriteAll(structStr); err != nil {
		fmt.Printf("copy to clipboard failed.\n\n%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("copied OK!")
}
