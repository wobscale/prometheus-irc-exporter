package stringprep

type TableElement struct {
	Lo  rune
	Hi  rune
	Map d // can be empty
}

type Table []TableElement

// Returns true if the rune is in table
func in_table(c rune, table Table) bool {
	for i := 0; i < len(table); i++ {
		if c == table[i].Lo || (table[i].Lo <= c && c <= table[i].Hi) {
			return true
		}
	}
	return false
}

// Returns a filtered rune sequence
func filter(input []rune, table Table) []rune {
	output := make([]rune, len(input))
	c := 0 // count

	for i := 0; i < len(input); i++ {
		if !in_table(input[i], table) {
			output[c] = input[i]
			c++
		}
	}

	return output[0:]
}

// mapLen returns our length of a map.
func mapLen(input d) int {
	for i := len(input) - 1; 0 <= i; i-- {
		if input[i] != 0 {
			return i + 1
		}
	}
	return 0
}

// Iterates over the input rune array and replaces runes with their maps
func map_table(input []rune, table Table) []rune {
	var output []rune

	for i := 0; i < len(input); i++ {
		// If rune is in table, replace it with its map
		if in_table(input[i], table) {
			for k := 0; k < len(table); k++ {
				if input[i] == table[k].Lo || (table[k].Lo <= input[i] && input[i] <= table[k].Hi) {
					if table[k].Map[0] != 0 || table[k].Map[1] != 0 || table[k].Map[2] != 0 || table[k].Map[3] != 0 {
						output = append(output, table[k].Map[0:mapLen(table[k].Map)]...)
					}
					break
				}
			}
		} else {
			output = append(output, input[i])
		}

	}
	return output
}

// Tables defines the various unicode tables.
var Tables = map[string]Table{
	"A1":  _A1,
	"B1":  _B1,
	"B2":  _B2,
	"B3":  _B3,
	"C11": _C11,
	"C12": _C12,
	"C21": _C21,
	"C22": _C22,
	"C3":  _C3,
	"C4":  _C4,
	"C5":  _C5,
	"C6":  _C6,
	"C7":  _C7,
	"C8":  _C8,
	"C9":  _C9,
	"D1":  _D1,
	"D2":  _D2,
}
