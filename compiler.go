package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var inputFile string
var fileString string
type symbol_table map[string]int

var keywords = map[string]string {
	"False": "OH GOD PLEASE NO",
	"True": "CHILLAX",

	"If": "WHERE ARE THE TURTLES",
	"Else": "HOW THE TURNTABLES",
	"EndIf": "WHY ARE THE WAY YOU ARE",

	"While": "I AM BEYONCE ALWAYS",
	"EndWhile": "SHES NOT YO HO NO MO",

	"PlusOperator": "I THINK THAT PRETTY MUCH SUMS IT UP",
	"MinusOperator": "IM NOT A MORON",
	"MultiplicationOperator": "SUE ME",
	"DivisionOperator": "ILL KILL YOU",
	"Modulo": "I AM COLLAR BLIND",

	"EqualTo": "I AM NOT SUPERSTITIOUS BUT IM A LTTLE STITIOUS",
	"GreaterThan": "MO MONEY MO PROBLEMS",
	"Or": "I KNOW EXACTLY WHAT TO DO BUT IN A MUCH MORE REAL SENSE I HAVE NO IDEA WHAT TO DO",
	"And": "FOOL ME ONCE STRIKE ONE FOOL ME TWICE STRIKE THREE",
	// "Not": "NOPE DONT LIKE THAT",

	"DeclareMethod": "SHOULD HAVE BURNT THE PLACE WHEN I HAD THE CHANCE",
	"NonVoidMethod": "I WANT PEOPLE TO BE AFRAID OF HOW MUCH THEY LOVE ME",
	"MethodArguments": "MAKE FRIENDS FIRST SALES SECOND LOVE THIRD IN NO PARTICULAR ORDER",
	"Return": "I LOVE INSIDE JOKES ID LOVE TO BE PART OF ONE SOMEDAY",
	"EndMethodDeclaration": "I DONT EVEN CONSIDER MYSELF PART OF SOCIETY",

	"CallMethod": "TELL HIM TO CALL ME ASAP AS POSSIBLE",
	"AssignVariableFromMethodCall": "WORLDS BEST BOSS",

	"DeclareInt": "I DIDNT SAY IT I DECLARED IT",
	"SetInitialValue": "GAME SET MATCH POINT GAME OVER END OF GAME",

	"BeginMain": "I DECLARE BANKRUPTCY",
	"EndMain": "IM DEAD INSIDE",

	"Print": "THATS WHAT SHE SAID",

	"Read": "IM RUNNING AWAY FROM MY RESPONSIBLITIES AND IT FEELS GOOD",

	"AssignVariable": "EXPLAIN THIS TO ME LIKE IM FIVE",
	"SetValue": "ITS NEVER TOO EARLY FOR ICE CREAM",
	"EndAssignVariable": "I UNDERSTAND NOTHING",

	"ParseError": "YOU MISS 100 PERCENT OF THE SHOTS YOU DONT TAKE - WAYNE GRETZKY - MICHEAL SCOTT",
}

type Token struct {
	tok_type string
	value string
	line_no int
}

func lexSource (program string) []Token {
	lines := strings.Split(program, "\n")
	token_list := make([]Token, 0)

	for line_idx, line := range lines {
		line = strings.TrimSpace(line)
		for k,v := range keywords {
			if strings.HasPrefix(line, v) {
				token_list = append(token_list, Token{k, v, line_idx+1})
			}
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, token_list[len(token_list)-1].value))
		if len(rest) != 0 {
			token_list = append(token_list, Token{"params", rest, line_idx+1})
		}
	}
	token_list = append(token_list, Token{"END_PROGRAM", "END PROGRAM", -1})
	return token_list
}

type AST struct {
	ast_tok Token
	children []AST
}

func ParseUnitExp(tok_list []Token) (AST, []Token) {
	return AST{tok_list[0], make([]AST, 0)}, tok_list[1:]
}

func ParseExp(tok_list []Token) (AST, []Token) {
	var left_ast AST
	left_ast, tok_list = ParseUnitExp(tok_list)

	for tok_list[0].tok_type == "DivisionOperator" || tok_list[0].tok_type == "MultiplicationOperator" ||
	tok_list[0].tok_type == "MinusOperator" || tok_list[0].tok_type == "PlusOperator" || tok_list[0].tok_type == "Modulo" ||
	tok_list[0].tok_type == "EqualTo" || tok_list[0].tok_type == "GreaterThan" || tok_list[0].tok_type == "Or" ||
	tok_list[0].tok_type == "And" {
		op_token := tok_list[0]
		var right_ast AST
		right_ast, tok_list = ParseExp(tok_list[1:])
		left_ast = AST{op_token, []AST{left_ast, right_ast}}
	}
	return left_ast, tok_list
}

func ParseStatements(tok_list []Token) ([]AST, []Token) {
	statement_list := make([]AST, 0)
	for tok_list[0].tok_type != "EndMain" && tok_list[0].tok_type != "EndMethodDeclaration" && tok_list[0].tok_type != "Else" && tok_list[0].tok_type != "EndIf" && tok_list[0].tok_type != "EndWhile" {
		if tok_list[0].tok_type == "Print" {
			print_token := AST{tok_list[0], make([]AST, 0)}

			if tok_list[1].tok_type != "params" {
				fmt.Println("ParseError: Expected a print value or string at line " + strconv.Itoa(tok_list[1].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			print_token.children = append(print_token.children, AST{tok_list[1], make([]AST, 0)})

			statement_list = append(statement_list, print_token)
			tok_list = tok_list[2:]
		} else if tok_list[0].tok_type == "DeclareInt" {
			decInt_token := AST{tok_list[0], make([]AST, 0)}

			if tok_list[1].tok_type != "params" {
				fmt.Println("ParseError: Expected a variable name at line " + strconv.Itoa(tok_list[1].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			varInt_token := AST{tok_list[1], make([]AST, 0)}

			if tok_list[2].tok_type != "SetInitialValue" {
				fmt.Println("ParseError: Wrong token, expected SetInitialValue but got " + tok_list[2].tok_type + " at line " + strconv.Itoa(tok_list[2].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			initVal_token := AST{tok_list[2], make([]AST, 0)}

			if tok_list[3].tok_type != "params" {
				fmt.Println("ParseError: Expected an integer value at line " + strconv.Itoa(tok_list[3].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			intVal_token := AST{tok_list[3], make([]AST, 0)}

			initVal_token.children = append(initVal_token.children, intVal_token)
			varInt_token.children = append(varInt_token.children, initVal_token)
			decInt_token.children = append(decInt_token.children, varInt_token)

			statement_list = append(statement_list, decInt_token)
			tok_list = tok_list[4:]
		} else if tok_list[0].tok_type == "AssignVariable" {
			assVar_token := AST{tok_list[0], make([]AST, 0)}

			if tok_list[1].tok_type != "params" {
				fmt.Println("ParseError: Expected a variable name at line " + strconv.Itoa(tok_list[1].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			varName_token := AST{tok_list[1], make([]AST, 0)}

			if tok_list[2].tok_type != "SetValue" {
				fmt.Println("ParseError: Wrong token, expected SetValue but got " + tok_list[2].tok_type + " at line " + strconv.Itoa(tok_list[2].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			setVal_token := AST{tok_list[2], make([]AST, 0)}

			var ast_list AST
			ast_list, tok_list = ParseExp(tok_list[3:])

			if tok_list[0].tok_type != "EndAssignVariable" {
				fmt.Println("ParseError: Wrong token, expected EndAssignVariable but got " + tok_list[0].tok_type + " at line " + strconv.Itoa(tok_list[0].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			// endAss_token := AST{tok_list[0], make([]AST, 0)}

			setVal_token.children = append(setVal_token.children, ast_list)
			varName_token.children = append(varName_token.children, setVal_token)
			assVar_token.children = append(assVar_token.children, varName_token)

			statement_list = append(statement_list, assVar_token)
			tok_list = tok_list[1:]
		} else if tok_list[0].tok_type == "If" {
			if_token := AST{tok_list[0], make([]AST, 0)}

			if tok_list[1].tok_type != "params" {
				fmt.Println("ParseError: Expected a variable or boolean constant name at line " + strconv.Itoa(tok_list[1].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			condVar_token := AST{tok_list[1], make([]AST, 0)}
			if_token.children = append(if_token.children, condVar_token)

			true_token := AST{Token{"True", "", -1}, make([]AST, 0)}

			var ast_list []AST
			ast_list, tok_list = ParseStatements(tok_list[2:])
			true_token.children = append(true_token.children, ast_list...)
			if_token.children = append(if_token.children, true_token)

			if tok_list[0].tok_type == "EndIf" {
				// endIf_token := AST{tok_list[0], make([]AST, 0)}

				statement_list = append(statement_list, if_token)
				tok_list = tok_list[1:]
			} else if tok_list[0].tok_type == "Else" {
				else_token := AST{tok_list[0], make([]AST, 0)}
				if_token.children = append(if_token.children, else_token)

				false_token := AST{Token{"False", "", -1}, make([]AST, 0)}
				var ast_list []AST
				ast_list, tok_list = ParseStatements(tok_list[1:])
				false_token.children = append(false_token.children, ast_list...)
				if_token.children = append(if_token.children, false_token)

				if tok_list[0].tok_type != "EndIf" {
					fmt.Println("ParseError: Wrong token, expected EndIf but got " + tok_list[0].tok_type + " at line " + strconv.Itoa(tok_list[0].line_no))
					fmt.Println(keywords["ParseError"])
					os.Exit(1)
				}
				// endIf_token := AST{tok_list[0], make([]AST, 0)}

				statement_list = append(statement_list, if_token)
				tok_list = tok_list[1:]
			} else {
				fmt.Println("ParseError: Wrong token, expected Else or EndIf but got " + tok_list[0].tok_type + " at line " + strconv.Itoa(tok_list[0].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			
		} else if tok_list[0].tok_type == "While" {
			while_token := AST{tok_list[0], make([]AST, 0)}

			if tok_list[1].tok_type != "params" {
				fmt.Println("ParseError: Expected a variable or boolean constant name at line " + strconv.Itoa(tok_list[1].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			condVar_token := AST{tok_list[1], make([]AST, 0)}
			while_token.children = append(while_token.children, condVar_token)

			true_token := AST{Token{"True", "", -1}, make([]AST, 0)}
			var ast_list []AST

			ast_list, tok_list = ParseStatements(tok_list[2:])

			true_token.children = append(true_token.children, ast_list...)
			while_token.children = append(while_token.children, true_token)

			if tok_list[0].tok_type != "EndWhile" {
				fmt.Println("ParseError: Wrong token, expected EndWhile but got " + tok_list[0].tok_type + " at line " + strconv.Itoa(tok_list[0].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			// endWhile_token := AST{tok_list[0], make([]AST, 0)}

			statement_list = append(statement_list, while_token)
			tok_list = tok_list[1:]
		} else if tok_list[0].tok_type == "AssignVariableFromMethodCall" {
			fromMC_token := AST{tok_list[0], make([]AST, 0)}

			if tok_list[1].tok_type != "params" {
				fmt.Println("ParseError: Expected a variable or boolean constant name at line " + strconv.Itoa(tok_list[1].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			varMC_token := AST{tok_list[1], make([]AST, 0)}

			if tok_list[2].tok_type != "CallMethod" {
				fmt.Println("ParseError: Wrong token, expected CallMethod but got " + tok_list[2].tok_type + " at line " + strconv.Itoa(tok_list[2].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			methCall_token := AST{tok_list[2], make([]AST, 0)}


			var method_token AST
			if tok_list[3].tok_type == "Read" {
				method_token = AST{tok_list[3], make([]AST, 0)}
			} else if tok_list[3].tok_type == "params" {
				method_token = AST{tok_list[3], make([]AST, 0)}
				for _, p := range strings.Split(tok_list[3].value, " ") {
					method_token.children = append(method_token.children, AST{Token{"params", p, tok_list[3].line_no}, make([]AST, 0)})
				}
			} else {
				fmt.Println("ParseError: Expected an function call at line " + strconv.Itoa(tok_list[3].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}			

			methCall_token.children = append(methCall_token.children, method_token)
			varMC_token.children = append(varMC_token.children, methCall_token)
			fromMC_token.children = append(fromMC_token.children, varMC_token)

			statement_list = append(statement_list, fromMC_token)
			tok_list = tok_list[4:]
		} else if tok_list[0].tok_type == "CallMethod" {
			methCall_token := AST{tok_list[0], make([]AST, 0)}

			var method_token AST
			if tok_list[1].tok_type == "Read" {
				method_token = AST{tok_list[1], make([]AST, 0)}
			} else if tok_list[1].tok_type == "params" {
				method_token = AST{tok_list[1], make([]AST, 0)}
				for _, p := range strings.Split(tok_list[1].value, " ") {
					method_token.children = append(method_token.children, AST{Token{"params", p, tok_list[1].line_no}, make([]AST, 0)})
				}
			} else {
				fmt.Println("ParseError: Expected an function call at line " + strconv.Itoa(tok_list[1].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}

			methCall_token.children = append(methCall_token.children, method_token)
			statement_list = append(statement_list, methCall_token)
			tok_list = tok_list[2:]
		} else if tok_list[0].tok_type == "Return" {
			return_token := AST{tok_list[0], make([]AST, 0)}

			if tok_list[1].tok_type != "params" {
				fmt.Println("ParseError: Expected a variable or integer at line " + strconv.Itoa(tok_list[1].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			return_token.children = append(return_token.children, AST{tok_list[1], make([]AST, 0)})

			statement_list = append(statement_list, return_token)
			tok_list = tok_list[2:]
		} else {
			fmt.Println("ParseError: Illegal token " + tok_list[0].tok_type + " at line " + strconv.Itoa(tok_list[0].line_no))
			fmt.Println(keywords["ParseError"])
			os.Exit(1)
		}
	}
	return statement_list, tok_list
}

func Parse (tok_list []Token) (AST, []Token) {

	if tok_list[0].tok_type == "BeginMain" {
		main_token := tok_list[0]
		ast := AST{main_token, make([]AST, 0)}
		var ast_list []AST
		ast_list, tok_list = ParseStatements(tok_list[1:])
		ast.children = ast_list

		if tok_list[0].tok_type != "EndMain" {
			fmt.Println("ParseError: Wrong token, expected and EndMain but got " + tok_list[0].tok_type + " at line " + strconv.Itoa(tok_list[0].line_no))
			fmt.Println(keywords["ParseError"])
			os.Exit(1)
		}
		// end_token := tok_list[0]
		return ast, tok_list[1:]
	} else if tok_list[0].tok_type == "DeclareMethod" {
		declare_method_token := AST{tok_list[0], make([]AST, 0)}

		if tok_list[1].tok_type != "params" {
			fmt.Println("ParseError: Expected a print value or string at line " + strconv.Itoa(tok_list[1].line_no))
			fmt.Println(keywords["ParseError"])
			os.Exit(1)
		}
		declare_method_token.children = append(declare_method_token.children, AST{tok_list[1], make([]AST, 0)})

		var method_params_token AST

		if tok_list[2].tok_type == "MethodArguments" {
			method_params_token = AST{tok_list[2], make([]AST, 0)}
			if tok_list[3].tok_type != "params" {
				fmt.Println("ParseError: Expected a variable at line " + strconv.Itoa(tok_list[3].line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
			method_params_token.children = append(method_params_token.children, AST{tok_list[3], make([]AST, 0)})
			tok_list = tok_list[4:]
			for tok_list[0].tok_type == "MethodArguments" {
				if tok_list[1].tok_type != "params" {
					fmt.Println("ParseError: Expected a variable at line " + strconv.Itoa(tok_list[1].line_no))
					fmt.Println(keywords["ParseError"])
					os.Exit(1)
				}
				method_params_token.children = append(method_params_token.children, AST{tok_list[1], make([]AST, 0)})
				tok_list = tok_list[2:]
			}
			declare_method_token.children = append(declare_method_token.children, method_params_token)
		} else {
			declare_method_token.children = append(declare_method_token.children, AST{Token{"NoMethodArguments", "", -1}, make([]AST, 0)})
		}		

		if tok_list[0].tok_type == "NonVoidMethod" {
			declare_method_token.children = append(declare_method_token.children, AST{tok_list[0], make([]AST, 0)})
			tok_list = tok_list[1:]
		} else {
			declare_method_token.children = append(declare_method_token.children, AST{Token{"VoidMethod", "", -1}, make([]AST, 0)})
		}

		var ast_list []AST
		ast_list, tok_list = ParseStatements(tok_list)
		declare_method_token.children = append(declare_method_token.children, ast_list...)

		if tok_list[0].tok_type != "EndMethodDeclaration" {
			fmt.Println("ParseError: Wrong token, expected and EndMethodDeclaration but got " + tok_list[0].tok_type + " at line " + strconv.Itoa(tok_list[0].line_no))
			fmt.Println(keywords["ParseError"])
			os.Exit(1)
		}
		// end_token := tok_list[0]
		return declare_method_token, tok_list[1:]
	} else {
		fmt.Println("ParseError: Wrong token, expected DeclareMethod or BeginMain but got " + tok_list[0].tok_type + " at line " + strconv.Itoa(tok_list[0].line_no))
		fmt.Println(keywords["ParseError"])
		os.Exit(1)
	}
	return AST{} , tok_list
}

func (state_map symbol_table) EvaluatePrint(ast AST) {
	if ast.children[0].ast_tok.value[0] == '"' {
		fmt.Println(ast.children[0].ast_tok.value[1:len(ast.children[0].ast_tok.value)-1])
	} else {
		if val, err := strconv.Atoi(ast.children[0].ast_tok.value); err == nil {
			fmt.Println(val)
		} else if val, ok := state_map[ast.children[0].ast_tok.value]; ok {
			fmt.Println(val)
		}
	}
}

func (state_map symbol_table) EvaluateDeclareInt(ast AST) {
	if _, ok := state_map[ast.children[0].ast_tok.value]; ok {
		fmt.Println("Redeclaring existing variable at line " + strconv.Itoa(ast.ast_tok.line_no) + ".")
		fmt.Println(keywords["ParseError"])
		os.Exit(1)
	}
	variable := ast.children[0].ast_tok.value
	value, ok := strconv.Atoi(ast.children[0].children[0].children[0].ast_tok.value)
	if ok == nil {
		state_map[variable] = value
	} else {
		fmt.Println("Expected an integer at line " + strconv.Itoa(ast.children[0].children[0].children[0].ast_tok.line_no) + ".")
		fmt.Println(keywords["ParseError"])
		os.Exit(1)
	}
}

func (state_map symbol_table) EvaluateExp(ast AST) int {
	if ast.ast_tok.tok_type == "DivisionOperator" {
		return state_map.EvaluateExp(ast.children[0]) / state_map.EvaluateExp(ast.children[1])
	} else if ast.ast_tok.tok_type == "MultiplicationOperator" {
		return state_map.EvaluateExp(ast.children[0]) * state_map.EvaluateExp(ast.children[1])
	} else if ast.ast_tok.tok_type == "MinusOperator" {
		return state_map.EvaluateExp(ast.children[0]) - state_map.EvaluateExp(ast.children[1])
	} else if ast.ast_tok.tok_type == "PlusOperator" {
		return state_map.EvaluateExp(ast.children[0]) + state_map.EvaluateExp(ast.children[1])
	} else if ast.ast_tok.tok_type == "Modulo" {
		return state_map.EvaluateExp(ast.children[0]) % state_map.EvaluateExp(ast.children[1])
	} else if ast.ast_tok.tok_type == "GreaterThan" {
		if state_map.EvaluateExp(ast.children[0]) < state_map.EvaluateExp(ast.children[1]) {
			return 1
		} else {
			return 0
		}
	} else if ast.ast_tok.tok_type == "EqualTo" {
		if state_map.EvaluateExp(ast.children[0]) == state_map.EvaluateExp(ast.children[1]) {
			return 1
		} else {
			return 0
		}
	} else if ast.ast_tok.tok_type == "And" {
		if state_map.EvaluateExp(ast.children[0]) == 1 && state_map.EvaluateExp(ast.children[1]) == 1 {
			return 1
		} else {
			return 0
		}
	} else if ast.ast_tok.tok_type == "Or" {
		if state_map.EvaluateExp(ast.children[0]) == 1 || state_map.EvaluateExp(ast.children[1]) == 1 {
			return 1
		} else {
			return 0
		}
	} else if ast.ast_tok.tok_type == "params" {
		if val, err := strconv.Atoi(ast.ast_tok.value); err == nil {
			return val
		} else if val, ok := state_map[ast.ast_tok.value]; ok {
			return val
		} else {
			fmt.Println("Expected a number or variable at line " + strconv.Itoa(ast.ast_tok.line_no) + ". Got " + ast.ast_tok.value)
			fmt.Println(keywords["ParseError"])
			os.Exit(1)
		}
	} else {
		fmt.Println("Unknown token at line " + strconv.Itoa(ast.ast_tok.line_no) + ". Got " + ast.ast_tok.value)
		os.Exit(1)
	}
	return 0
}

func (state_map symbol_table) EvaluateAssignVariable(ast AST) {
	if _, ok := state_map[ast.children[0].ast_tok.value]; !ok {
		fmt.Println("Settin undeclared variable at line " + strconv.Itoa(ast.ast_tok.line_no) + ".")
		fmt.Println(keywords["ParseError"])
		os.Exit(1)
	}
	variable := ast.children[0].ast_tok.value
	state_map[variable] = state_map.EvaluateExp(ast.children[0].children[0].children[0])
}

func (state_map symbol_table) EvaluateFuncCall(ast AST, FST map[string]AST, func_params []int) (int) {
	method_arguments := ast.children[1]

	if len(method_arguments.children) < len(func_params) {
		fmt.Println("Too many parameters at line " + strconv.Itoa(method_arguments.ast_tok.line_no))
	} else if len(method_arguments.children) > len(func_params) {
		fmt.Println("Too few parameters at line " + strconv.Itoa(method_arguments.ast_tok.line_no))
	}

	for i := 0; i < len(method_arguments.children); i++ {
		state_map[method_arguments.children[i].ast_tok.value] = func_params[i]
	}
	if ast.children[2].ast_tok.tok_type == "NonVoidMethod" {
		return state_map.Evaluate(ast, FST)
	} else {
		state_map.Evaluate(ast, FST)
	}
	return 0
}

func (state_map symbol_table) Evaluate(MST AST, FST map[string]AST) int {
	for _, statement := range MST.children {
		if statement.ast_tok.tok_type == "Print" {
			state_map.EvaluatePrint(statement)
		} else if statement.ast_tok.tok_type == "DeclareInt" {
			state_map.EvaluateDeclareInt(statement)
		} else if statement.ast_tok.tok_type == "AssignVariable" {
			state_map.EvaluateAssignVariable(statement)
		} else if statement.ast_tok.tok_type == "If" {
			if_var := statement.children[0]
			param := if_var.ast_tok.value
			if param[0] == '@' {
				if param[1:] == keywords["True"] {
					state_map.Evaluate(statement.children[1], FST)
				} else {
					if (len(statement.children) == 4) {
						state_map.Evaluate(statement.children[3], FST)
					}
				}
			} else if val , ok := state_map[param]; ok && val != 0 {
				state_map.Evaluate(statement.children[1], FST)
			} else {
				if (len(statement.children) == 4) {
					state_map.Evaluate(statement.children[3], FST)
				}
			}
		} else if statement.ast_tok.tok_type == "While" {
			while_var := statement.children[0]
			param := while_var.ast_tok.value
			if param[0] == '@' {
				if param[1:] == keywords["True"] {
					fmt.Println("Infinite loop starting at line " + strconv.Itoa(statement.ast_tok.line_no) + ". No way to break out. ")
					fmt.Println(keywords["ParseError"])
					os.Exit(1)
				}
			} else if _ , ok := state_map[param]; ok {
				for state_map[param] != 0 {
					state_map.Evaluate(statement.children[1], FST)
				}
			}
		} else if statement.ast_tok.tok_type == "AssignVariableFromMethodCall" {
			var_token := statement.children[0]

			if _, ok := state_map[var_token.ast_tok.value]; !ok {
				fmt.Println("Attempt to assign value to undeclared variable at line " + strconv.Itoa(var_token.ast_tok.line_no))
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}

			method_token := var_token.children[0].children[0]
			if method_token.ast_tok.tok_type == "Read" {
				var value int
				_, err := fmt.Scanf("%d", &value)
				if err != nil {
					fmt.Println("Could not parse input")
					os.Exit(1)
				}
				state_map[var_token.ast_tok.value] = value
				continue
			}
			func_AST := FST[method_token.children[0].ast_tok.value]
			func_params := make([]int, 0)
			for i := 1; i < len(method_token.children); i++  {
				if val, err := strconv.Atoi(method_token.children[i].ast_tok.value); err == nil {
					func_params = append(func_params, val)
				} else if val, ok := state_map[method_token.children[i].ast_tok.value]; ok {
					func_params = append(func_params, val)
				} else {
					fmt.Println("Expected a number or variable at line " + strconv.Itoa(method_token.ast_tok.line_no) + ". Got " + method_token.children[i].ast_tok.value)
					fmt.Println(keywords["ParseError"])
					os.Exit(1)
				}				
			}
			var func_state = make(symbol_table)

			state_map[var_token.ast_tok.value] = func_state.EvaluateFuncCall(func_AST, FST, func_params)
		} else if statement.ast_tok.tok_type == "CallMethod" {
			method_token := statement.children[0]
			if method_token.ast_tok.tok_type == "Read" {
				var value int
				_, err := fmt.Scan(&value)
				if err != nil {
					fmt.Println("Could not parse input")
					os.Exit(1)
				}
			}
			func_AST := FST[method_token.children[0].ast_tok.value]
			func_params := make([]int, 0)
			for i := 1; i < len(method_token.children); i++  {
				if val, err := strconv.Atoi(method_token.children[i].ast_tok.value); err == nil {
					func_params = append(func_params, val)
				} else if val, ok := state_map[method_token.children[i].ast_tok.value]; ok {
					func_params = append(func_params, val)
				} else {
					fmt.Println("Expected a number or variable at line " + strconv.Itoa(method_token.ast_tok.line_no) + ". Got " + method_token.children[i].ast_tok.value)
					fmt.Println(keywords["ParseError"])
					os.Exit(1)
				}				
			}
			var func_state = make(symbol_table)

			func_state.EvaluateFuncCall(func_AST, FST, func_params)
		} else if statement.ast_tok.tok_type == "Return" {
			if val, err := strconv.Atoi(statement.children[0].ast_tok.value); err == nil {
				return val
			} else if val, ok := state_map[statement.children[0].ast_tok.value]; ok {
				return val
			} else {
				fmt.Println("Expected a number or variable at line " + strconv.Itoa(statement.ast_tok.line_no) + ". Got " + statement.children[0].ast_tok.value)
				fmt.Println(keywords["ParseError"])
				os.Exit(1)
			}
		}
	}
	return 0;
}

func main() {

	var fileName = os.Args[1]
	var b, _ = ioutil.ReadFile(fileName)
	fileString = string(b)
	var tok_list = lexSource(fileString)

	var function_syntax_tree = make(map[string]AST)
	var main_syntax_tree AST

	for tok_list[0].tok_type != "END_PROGRAM" {
		var ST AST
		ST, tok_list = Parse(tok_list)
		if ST.ast_tok.tok_type == "BeginMain" {
			main_syntax_tree = ST
		} else {
			function_syntax_tree[ST.children[0].ast_tok.value] = ST
		}
	}

	// fmt.Println("Main method")
	// fmt.Println(main_syntax_tree)

	// for k,v := range function_syntax_tree {
	// 	fmt.Println("Function " + k)
	// 	fmt.Println(v)
	// }

	var main_state = make(symbol_table)
	main_state.Evaluate(main_syntax_tree, function_syntax_tree)
}