package main

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func NewComputer(program string) (*Computer, error) {
	c := &Computer{
		Registers: map[byte]int{'A': 0, 'B': 0, 'C': 0},
		Inputs:    []byte{},
		PC:        0,
	}

	for i, line := range strings.Split(program, "\n") {
		line = strings.TrimSpace(line)
		debug("reading line %d: %s", i+1, line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "Register ") {
			line = strings.TrimPrefix(line, "Register ")
			regN, regV, ok := strings.Cut(line, ": ")
			if !ok {
				return nil, fmt.Errorf("invalid line %d: %s", i+1, line)
			}
			regI, err := strconv.Atoi(regV)
			if err != nil {
				return nil, fmt.Errorf("invalid number in line %d: %s", i+1, line)
			}
			c.Registers[regN[0]] = regI
			continue
		}

		if strings.HasPrefix(line, "Program: ") {
			line = strings.TrimPrefix(line, "Program: ")
			for _, i := range strings.Split(line, ",") {
				v, err := strconv.Atoi(i)
				if err != nil {
					return nil, fmt.Errorf("invalid number in program: %s", line)
				}
				c.Inputs = append(c.Inputs, byte(v&0b111))
			}
			continue
		}

		return nil, fmt.Errorf("unkown line %d: %s", i+1, line)
	}
	return c, nil
}

type Instruction struct {
	Eval func(operand byte, c *Computer)
	Name string
}

func Div(name string, targetReg byte) Instruction {
	return Instruction{
		Name: name,
		Eval: func(operand byte, c *Computer) {
			numerator := c.Registers['A']
			denominator := c.Combo(operand)
			c.Registers[targetReg] = numerator >> denominator
			debug("  %s %d %[2]b", string(targetReg), c.Registers[targetReg])
		},
	}
}

var instructions = map[byte]Instruction{
	0: Div("adv", 'A'),
	1: {
		Name: "bxl",
		Eval: func(operand byte, c *Computer) {
			c.Registers['B'] = c.Registers['B'] ^ int(operand)
			debug("  B %d %[1]b", c.Registers['B'])
		},
	},
	2: {
		Name: "bst",
		Eval: func(operand byte, c *Computer) {
			c.Registers['B'] = c.Combo(operand) & 0b111
			debug("  B %d %[1]b", c.Registers['B'])
		},
	},
	3: {
		Name: "jnz",
		Eval: func(operand byte, c *Computer) {
			if c.Registers['A'] == 0 {
				return
			}
			c.PC = int(operand) & 0b111
			debug("  PC %d", c.PC)
		},
	},
	4: {
		Name: "bxc",
		Eval: func(operand byte, c *Computer) {
			c.Registers['B'] = c.Registers['B'] ^ c.Registers['C']
			debug("  B %d %[1]b", c.Registers['B'])
		},
	},
	5: {
		Name: "out",
		Eval: func(operand byte, c *Computer) {
			res := c.Combo(operand) & 0b111
			c.Output = append(c.Output, byte(res))
			debug("  output %d  %[1]b", res)
		},
	},
	6: Div("bdv", 'B'),
	7: Div("cdv", 'C'),
}

type Computer struct {
	Registers map[byte]int
	Inputs    []byte
	PC        int
	Output    []byte
}

func (c *Computer) String() string {
	b := bytes.Buffer{}
	regNames := []byte{}
	for k := range c.Registers {
		regNames = append(regNames, k)
	}
	slices.Sort(regNames)
	for _, k := range regNames {
		fmt.Fprintf(&b, "%s=%d ", string(k), c.Registers[k])
	}
	fmt.Fprintf(&b, "PC=%d", c.PC)
	return b.String()
}

func (c *Computer) Run() []byte {
	for {
		pc := c.PC
		if pc < 0 || pc >= len(c.Inputs)-1 {
			debug("halting: %d", pc)
			break
		}

		opcode, operand := c.Inputs[pc], c.Inputs[pc+1]
		inst := instructions[opcode]
		debug("-----\n  %s %d on %s", inst.Name, operand, c)
		inst.Eval(operand, c)
		if c.PC == pc {
			c.PC += 2
		}
	}
	return c.Output
}

func (c *Computer) RunExpect(expected []byte) []byte {
	for i := 1; ; i++ {
		pc := c.PC
		if pc < 0 || pc >= len(c.Inputs)-1 {
			debug("halting: %d", pc)
			break
		}

		opcode, operand := c.Inputs[pc], c.Inputs[pc+1]
		inst := instructions[opcode]
		debug("%d -----\n  %s %d on %s", i, inst.Name, operand, c)
		inst.Eval(operand, c)
		if c.PC == pc {
			c.PC += 2
		}
		if opcode == 5 && !slices.Equal(c.Output, expected[0:len(c.Output)]) {
			// fmt.Println("abording at len", len(c.Output))
			return nil
		}
	}
	return c.Output
}

func joinRes(i []byte) string {
	b := bytes.Buffer{}
	for i, v := range i {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('0' + v)
	}
	return b.String()
}

func (c *Computer) Combo(v byte) int {
	// Combo operands 0 through 3 represent literal values 0 through 3.
	// Combo operand 4 represents the value of register A.
	// Combo operand 5 represents the value of register B.
	// Combo operand 6 represents the value of register C.
	// Combo operand 7 is reserved and will not appear in valid programs.
	switch v {
	case 0:
		debug("  combo 0")
		return 0
	case 1:
		debug("  combo 1")
		return 1
	case 2:
		debug("  combo 2")
		return 2
	case 3:
		debug("  combo 3")
		return 3
	case 4:
		debug("  combo A %d", c.Registers['A'])
		return c.Registers['A']
	case 5:
		debug("  combo B %d", c.Registers['B'])
		return c.Registers['B']
	case 6:
		debug("  combo C %d", c.Registers['C'])
		return c.Registers['C']
	case 7:
		panic("reserved")
	default:
		panic("invalid combo operand")
	}
}

func runA(file string) (int, error) {
	in, err := os.ReadFile(file)
	if err != nil {
		return -1, err
	}

	c, err := NewComputer(string(in))

	for _, r := range c.Registers {
		fmt.Printf("%3b ", r)
	}
	fmt.Println()

	if err != nil {
		return -1, err
	}
	debug("%+v", c)

	res := c.Run()
	fmt.Println("output", joinRes(res))
	for _, r := range res {
		fmt.Printf("%3b ", r)
	}
	fmt.Println()

	return -1, nil
}

func runB(file string) (int, error) {
	in, err := os.ReadFile(file)
	if err != nil {
		return -1, err
	}

	c, err := NewComputer(string(in))
	if err != nil {
		return -1, err
	}
	c.Registers['A'] = 0
	c.Registers['B'] = 0
	c.Registers['C'] = 0

	debug("%+v", c)

	stop := false
	cnt := &atomic.Int64{}
	cnt.Add(86063176041)
	wg := &sync.WaitGroup{}
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	quit := make(chan struct{})

	go func() {
		start := time.Now()
		lastT := start
		lastV := cnt.Load()

		for {
			select {
			case t := <-ticker.C:
				v := cnt.Load()
				tdiff := t.Sub(lastT)
				fmt.Printf("progress %d | %.0f/sec | running %s\n", v, float64(v-lastV)/tdiff.Seconds(), time.Since(start))
				lastV = v
				lastT = t
			case <-quit:
				return
			}
		}
	}()

	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for !stop {
				v := cnt.Add(1) - 1
				tc := Computer{
					Inputs:    slices.Clone(c.Inputs),
					Registers: maps.Clone(c.Registers),
				}
				tc.Registers['A'] = int(v)
				res := tc.RunExpect(c.Inputs)
				if !slices.Equal(res, c.Inputs) {
					continue
				}
				stop = true
				fmt.Println("EQUAL", v)
				return
			}
		}()
	}

	wg.Wait()
	quit <- struct{}{}

	return -1, nil
}

const debugEnabled = true

func debug(msg string, args ...any) {
	if !debugEnabled {
		return
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, msg)
		return
	}
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func main() {

	if _, err := runA("input-test.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
		os.Exit(1)
	}

	// if _, err := runB("input.txt"); err != nil {
	// 	fmt.Fprintf(os.Stderr, "puzzle errored: %v\n", err)
	// 	os.Exit(1)
	// }
}
