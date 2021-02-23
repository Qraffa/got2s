package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"got2s/util"
	"io/ioutil"
	"os"
)

var (
	s              string
	packageNameSql string
	isJsonSql      bool
	isImportSql    bool
	isFuncSql      bool
)

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Parse create table statement to struct",
	Run:   parseSqlFunc,
}

func init() {
	sqlCmd.Flags().StringVarP(&s, "sql", "s", "", "set sql file")
	sqlCmd.MarkFlagRequired("sql")

	sqlCmd.Flags().StringVarP(&packageNameSql, "package", "p", "", "the name of struct package")
	sqlCmd.Flags().BoolVarP(&isJsonSql, "json", "j", false, "add json tag")
	sqlCmd.Flags().BoolVarP(&isImportSql, "import", "i", false, "auto import needed package")
	sqlCmd.Flags().BoolVarP(&isFuncSql, "func", "f", false, "add TableName function")

	rootCmd.AddCommand(sqlCmd)
}

func parseSqlFunc(cmd *cobra.Command, args []string) {
	data, err := ioutil.ReadFile(s)
	sql := string(data)
	if err != nil {
		fmt.Printf("read sql file %s failed. err: %s\n", s, err.Error())
		os.Exit(1)
	}
	var isPackage bool
	if packageNameSql != "" {
		isPackage = true
	}
	cfg := &util.Option{
		IsSQL:     true,
		Sql:       sql,
		IsJson:    isJsonSql,
		IsImport:  isImportSql,
		IsFunc:    isFuncSql,
		IsPackage: isPackage,
		Package:   packageNameSql,
	}
	structStr, err := util.Parse(cfg)
	if err != nil {
		fmt.Printf("parse sql statement to struct failed.\n\n%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(structStr)
	if err := clipboard.WriteAll(structStr); err != nil {
		fmt.Printf("copy to clipboard failed.\n\n%s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("copied OK!")
}
