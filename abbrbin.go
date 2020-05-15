package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

const aliasPath = ".config/abbr/"
const aliasFile = "aliases"
const aliasTmpFile = "aliastmp"

func main() {
	numberOfArgs := len(os.Args) - 1
	if numberOfArgs == 0 {
		showHelp()
		return
	}

	if os.Args[1] == "add" || os.Args[1] == "a" {
		if numberOfArgs == 2 {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			addAlias(os.Args[2], scanner.Text())
		} else if numberOfArgs == 3 {
			addAlias(os.Args[2], os.Args[3])
		} else {
			fmt.Print("Error: Unexpected arguments.\n\n")
			showHelp()
			return
		}
	} else if os.Args[1] == "list" || os.Args[1] == "l" {
		if numberOfArgs != 1 {
			fmt.Print("Error: Unexpected arguments.\n\n")
			showHelp()
			return
		}
		listAliases()
	} else if os.Args[1] == "move" || os.Args[1] == "m" {
		if numberOfArgs != 3 {
			fmt.Print("Error: Unexpected arguments.\n\n")
			showHelp()
			return
		}
		moveAlias(os.Args[2], os.Args[3])
	} else if os.Args[1] == "put" || os.Args[1] == "p" {
		if numberOfArgs == 2 {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			putAlias(os.Args[2], scanner.Text())
			// todo: last command, if stdin empty
		} else if numberOfArgs == 3 {
			putAlias(os.Args[2], os.Args[3])
		} else {
			fmt.Print("Error: Unexpected arguments.\n\n")
			showHelp()
			return
		}
	} else if os.Args[1] == "remove" || os.Args[1] == "r" {
		if numberOfArgs != 2 {
			fmt.Print("Error: Unexpected arguments.\n\n")
			showHelp()
			return
		}
		removeAlias(os.Args[2])
	} else if os.Args[1] == "show" || os.Args[1] == "s" {
		if numberOfArgs != 2 {
			fmt.Print("Error: Unexpected arguments.\n\n")
			showHelp()
			return
		}
		showAlias(os.Args[2])
	} else {
		showHelp()
	}
}

func showHelp() {
	message := `Usage: abbr [(a)dd|(l)ist|(m)ove|(p)ut|(r)emove|(s)how] args...
    
 add <name> <command>   Adds a new alias <name> for the command <command>
 add <name>             Adds a new alias <name> for one of either:
                         - Command given as string that is piped in
                         - Last executed command according to history
 list                   Prints a list of all managed aliases
 move <name> <new name> Changes the alias name from <name> to <new name>
 put <name> <command>   Like "add" but may overwrite existing alias
 put <name>             Like "add" but may overwrite existing alias
 remove <name>          Removes the alias of name <name>
 show <name>            Shows the definition of alias <name>
`
	fmt.Println(message)
}

func addAlias(alias string, command string) {
	if !checkAliasNameValid(alias) {
		return
	}

	if _, ok := getAliasList()[alias]; ok {
		fmt.Println("Alias already exists.")
	} else {
		addAliasExecute(alias, command)
	}
}

func moveAlias(aliasOld string, aliasNew string) {
	if command, ok := getAliasList()[aliasOld]; ok {
		fmt.Println(command)
		removeAliasExecute(aliasOld)
		addAlias(aliasNew, command)
	} else {
		fmt.Println("Alias does not exist.")
	}
}

func listAliases() {
	if list := getAliasList(); len(list) > 0 {
		fmt.Println("Alias            Command")
		fmt.Println("------------------------")

		keys := make([]string, len(list))
		for k := range list {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, alias := range keys {
			if alias != "" {
				if len(alias) <= 16 {
					fmt.Printf("%-16s %s\n", alias, list[alias])
				} else {
					fmt.Println(alias)
					fmt.Println("                ", list[alias])
				}
			}
		}
	}
}

func putAlias(alias string, command string) {
	if !checkAliasNameValid(alias) {
		return
	}
	removeAliasExecute(alias)
	addAliasExecute(alias, command)
}

func removeAlias(alias string) {
	if !checkAliasNameValid(alias) {
		return
	}
	// todo: error if alias does not exist?
	removeAliasExecute(alias)
}

func showAlias(alias string) {
	if command, ok := getAliasList()[alias]; ok {
		fmt.Println(command)
	}
}

func addAliasExecute(alias string, command string) {
	writeEntryToFile(getAbsoluteFilePath(), false, alias, command)
	writeEntryToFile(getAbsoluteTmpFilePath(), false, alias, command)
}

func removeAliasExecute(alias string) {
	file := getAbsoluteFilePath()

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(content), "\n")
	result := make([]string, 0)

	for _, line := range lines {
		if !strings.HasPrefix(line, "alias "+alias+"=") && strings.HasPrefix(line, "alias ") {
			result = append(result, line)
		}
	}

	err = ioutil.WriteFile(file, []byte(strings.Join(result, "\n")), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	writeEntryToFile(getAbsoluteTmpFilePath(), true, alias, "")
}

func writeEntryToFile(filePath string, unaliasFlag bool, alias string, command string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	var content string
	if unaliasFlag {
		content = "\nunalias " + alias + "\n"
	} else {
		content = "\nalias " + alias + "=\"" + escapeString(command) + "\"\n"
	}
	if _, err := file.Write([]byte(content)); err != nil {
		log.Fatalln(err)
	}
}

func getAbsoluteFilePath() string {
	return os.Getenv("HOME") + "/" + aliasPath + aliasFile
}

func getAbsoluteTmpFilePath() string {
	return os.Getenv("HOME") + "/" + aliasPath + aliasTmpFile
}

func checkAliasNameValid(alias string) bool {
	if len(alias) == 0 {
		fmt.Println("Alias name cannot be empty.")
		return false
	}
	if alias == "alias" || alias == "unalias" {
		fmt.Println("Alias name cannot be \"alias\" or \"unalias\".")
		return false
	}
	for c := range alias {
		if c == ' ' || c == '\n' || c == '\t' {
			fmt.Println("Alias name must not contain whitespaces.")
			return false
		}
	}

	matched, err := regexp.MatchString("^[\\w~@#%\\^\\.,][\\w-~@#%\\^\\.,]*$", alias) // should not start with '-'
	if err != nil {
		log.Fatalln(err)
	}
	if !matched {
		fmt.Println("Alias name contains illegal character.")
	}
	return matched
}

func getAliasList() map[string]string {
	list := make(map[string]string)

	file, err := os.Open(getAbsoluteFilePath())
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if line[0:6] != "alias " {
			continue
		}

		var aliasName string
		equalSignPosition := 0

		for pos, char := range line {
			if pos < 5 || (char == ' ' && aliasName == "") {
				continue
			}
			if char == '=' {
				equalSignPosition = pos
				break
			}
			if char != ' ' {
				aliasName += string(char)
			}
		}

		if aliasContentString := line[equalSignPosition+1:]; aliasContentString[0] == '"' && aliasContentString[len(aliasContentString)-1] == '"' {
			if aliasContent, err := unescapeString(aliasContentString[1 : len(aliasContentString)-1]); err == nil {
				list[aliasName] = aliasContent
			}
		}
	}

	return list
}

func escapeString(str string) string {
	var result string
	for _, c := range str {
		if c == '\\' {
			result += `\\`
		} else if c == '"' {
			result += `\"`
		} else {
			result += string(c)
		}
	}
	return result
}

func unescapeString(str string) (string, error) {
	var result string
	lastWasBackslash := false
	for _, c := range str {
		if c == '\\' && !lastWasBackslash {
			lastWasBackslash = true
			continue
		}

		if lastWasBackslash {
			if c == '\\' {
				result += `\`
			} else if c == '"' {
				result += `"`
			} else {
				return "", errors.New("String format error")
			}
		} else {
			if c != '"' {
				result += string(c)
			} else {
				return "", errors.New("String format error")
			}
		}
		lastWasBackslash = false
	}
	return result, nil
}
