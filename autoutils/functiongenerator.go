/*
Copyright (C) 2019 Leo Tenenbaum

This file is part of AutoArt.

AutoArt is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

AutoArt is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with AutoArt.  If not, see <https://www.gnu.org/licenses/>.
*/

package autoutils

import (
    "math"
    "math/rand"
    "fmt"
)

// Operators
const (
    CONST = iota
    ADD
    SUB
    MUL
    DIV
	MIN
	MAX
    SQRT
    SIN
    COS
    TAN
	LOG
	EXP
    OPERATOR_COUNT
)
const FIRST_BINARY = ADD
const FIRST_UNARY = SQRT
const BINARY_COUNT = FIRST_UNARY - 1 // -1 for CONST
const UNARY_COUNT = OPERATOR_COUNT - FIRST_UNARY
const FIRST_VAR = OPERATOR_COUNT

type Operator struct {
    op int // Operator number. If operator is a variable, v, it is equal to FIRST_VAR + v
    constant float64 // Constant (if op = CONST)
}

type Function struct {
    nvars int
    operators []Operator
}

// Generate a random function f with the given length (i.e. len(f.operators))
// and the given number of variables
func (f *Function) Generate(nvars int, length int) {
    f.nvars = nvars
    f.operators = make([]Operator, length)
    nsOnStack := 0
    i := 0
    for nsOnStack + i < length {
        var operator Operator
        var optype int
        if nsOnStack == 0 {
            // Pick a random variable
            optype = 0
        } else if nsOnStack == 1 {
            // Pick a constant/variable/unary operator
            optype = rand.Intn(3)
        } else {
            // Pick a constant/variable/unary/binary operator
            optype = rand.Intn(4)
        }
        switch optype {
        case 0:
            // variable
            operator.op = FIRST_VAR + rand.Intn(nvars)
            nsOnStack++
        case 1:
            // Constant
            operator.op = CONST
            operator.constant = rand.Float64()
            nsOnStack++
        case 2:
            // unary
            operator.op = rand.Intn(UNARY_COUNT) + FIRST_UNARY
        case 3:
            // binary
            operator.op = rand.Intn(BINARY_COUNT) + FIRST_BINARY
            nsOnStack--
        }
        f.operators[i] = operator
        i++
    }

    if nsOnStack + i == length {
        // Add a unary operator
        f.operators[i].op = rand.Intn(UNARY_COUNT) + FIRST_UNARY
        i++
    }

    // Keep adding binary operators until nsOnStack == 1
    for nsOnStack > 1 {
        f.operators[i].op = rand.Intn(BINARY_COUNT) + FIRST_BINARY
        nsOnStack--
        i++
    }

}

func (f *Function) Evaluate(vars []float64) float64 {
    var stack []float64
    for _, op := range f.operators {
        l := len(stack)
        switch (op.op) {
        case CONST:
            stack = append(stack, op.constant)
        case ADD:
            stack[l-2] += stack[l-1]
            stack = stack[:l-1]
        case SUB:
            stack[l-2] -= stack[l-1]
            stack = stack[:l-1]
        case MUL:
            stack[l-2] *= stack[l-1]
            stack = stack[:l-1]
        case DIV:
            if stack[l-1] == 0 { // Check for division by 0
                stack[l-1] = 0.01
            }
            stack[l-2] /= stack[l-1]
            stack = stack[:l-1]
		case MIN:
			stack[l-2] = math.Min(stack[l-2], stack[l-1])
			stack = stack[:l-1]
		case MAX:
			stack[l-2] = math.Max(stack[l-2], stack[l-1])
			stack = stack[:l-1]
        case SQRT:
            stack[l-1] = math.Sqrt(math.Abs(stack[l-1]))
        case SIN:
            stack[l-1] = math.Sin(stack[l-1])
        case COS:
            stack[l-1] = math.Cos(stack[l-1])
        case TAN:
            stack[l-1] = math.Tan(stack[l-1])
		case LOG:
			stack[l-1] = math.Log(math.Abs(stack[l-1]))
		case EXP:
			stack[l-1] = math.Exp(stack[l-1])
        default:
            stack = append(stack, vars[op.op - FIRST_VAR])
        }
    }
    return stack[0]
}

func (f *Function) String() string {
    var str string
    for _, op := range f.operators {
        switch (op.op) {
        case CONST:
            str += fmt.Sprintf("%v",op.constant)
        case ADD:
            str += "+"
        case SUB:
            str += "-"
        case MUL:
            str += "*"
        case DIV:
            str += "/"
        case SQRT:
            str += "sqrt"
        case SIN:
            str += "sin"
        case COS:
            str += "cos"
        case TAN:
            str += "tan"
        default:
            str += fmt.Sprintf("v%v", op.op - FIRST_VAR)
        }
        str += " "
    }
    return str
}

const mutationRate = 0.01

func (f *Function) Mutate() {
	for i, op := range f.operators {
		if op.op == CONST && rand.Float64() < mutationRate {
			f.operators[i].constant += rand.NormFloat64() / 5 // Nudge constant
		}
	}
}

func (f *Function) Breed(f1* Function, f2 *Function) {
	// f(x) = (f1(x) + f2(x)) / 2
	f.operators = make([]Operator, len(f1.operators) + len(f2.operators) + 3)
	for i, o := range f1.operators {
		f.operators[i] = o
	}
	for i, o := range f2.operators {
		f.operators[i + len(f1.operators)] = o
	}
	i := len(f1.operators) + len(f2.operators)
	f.operators[i].op = ADD
	f.operators[i+1].op = CONST
	f.operators[i+1].constant = 2
	f.operators[i+2].op = DIV
}

func (f *Function) CopyFrom(other *Function) {
	f.nvars = other.nvars
	f.operators = make([]Operator, len(other.operators))
	for i, op := range other.operators {
		f.operators[i] = op
	}
}
