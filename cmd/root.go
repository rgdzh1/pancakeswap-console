/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed To in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/walletConsole/pancakeswap-console/config"
	"log"
	"os"
	"path/filepath"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pancakeswap-console",
	Short: "PancakeSwap Console",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PancakeSwap Console version is v1.0")
		//fmt.Printf("%v",  config.CF)
	},
}

// Execute adds all child commands To the root command and sets flags appropriately.
// This is called by main.main(). It only needs To happen once To the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/tool-conf)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}

		sprintf := fmt.Sprintf("%s/%s", dir, "tool-conf.yaml")
		//fmt.Println(sprintf)
		_, err = os.Stat(sprintf)
		if err == nil {
			viper.AddConfigPath(dir)
			viper.SetConfigName("tool-conf")
		} else if os.IsNotExist(err) {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in home directory with name ".eth-tool" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName("tool-conf")
		} else {
			log.Fatal(fmt.Errorf("miss tool-conf file ($HOME/tool-conf or ./tool-conf)"))
		}

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
		if err := viper.Unmarshal(&config.CF); err != nil {
			panic(err)
		}
	}

	if len(config.CF.PrivateKey) == 0 {
		log.Fatal(fmt.Errorf("missing PrivateKey )"))
	}
	//log.Printf("%v",config.CF.BscToken)
	//client, err := ethclient.Dial(config.CF.RpcUrl)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//utils.Client = client

}

func initCommonChainFlags(command *cobra.Command) {
	command.Flags().Uint64VarP(&Gas, "gas", "g", 6, "Gas price(gwei)")
	//command.Flags().StringVarP(&Chain, "chain", "c", "", "链(eth,heco,bsc)")
	//command.MarkFlagRequired("chain")
}

func initCommonMnemonicFlags(command *cobra.Command) {
	//command.Flags().StringVarP(&ImportMnemonic, "mnemonic", "m", "", "助记词")
	command.Flags().IntVarP(&Start, "start", "s", 0, "start index")
	command.Flags().IntVarP(&Num, "num", "n", 1, "address number")

	//command.MarkFlagRequired("mnemonic")
}
