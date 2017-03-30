package main

import (
	"io"

	"github.com/tdewolff/parse/css"
)

func process(src io.Reader) ([]byte, error) {
	parser := css.NewParser(src, false)
	out := ""
	indentLevel := 0
processing:
	for {
		switch gt, _, data := parser.Next(); {
		case gt == css.ErrorGrammar, parser.Err() != nil:
			if err := parser.Err(); err != io.EOF {
				return nil, err
			}
			break processing
		case gt == css.AtRuleGrammar, gt == css.BeginAtRuleGrammar:
			for _, val := range parser.Values() {
				out += string(val.Data)
			}
			if gt == css.BeginAtRuleGrammar {
				indentLevel++
				out += " {\n"
			} else if gt == css.AtRuleGrammar {
				out += ";\n"
			}
		case gt == css.BeginRulesetGrammar:
			for i, val := range parser.Values() {
				switch val.TokenType {
				case css.LeftBracketToken, css.RightBracketToken, css.ColonToken, css.IdentToken:
					out += string(val.Data)
				case css.CommaToken:
					if i != 0 {
						out += ",\n"
						for i := 0; i < indentLevel; i++ {
							out += "\t"
						}
					}
				case css.DelimToken:
					if string(val.Data) == ">" {
						out += " " + string(val.Data) + " "
						continue
					}
					out += string(val.Data)
				case css.StringToken, css.FunctionToken, css.LeftParenthesisToken, css.RightParenthesisToken, css.HashToken:
					out += string(val.Data)
				case css.WhitespaceToken:
					out += " "
				default:
					out += "wtf:'" + string(val.Data) + "'"
					report("got a wtf token of type %q: %s. please open an issue.", val.String(), val.Data)
				}
				for i := 0; i < indentLevel; i++ {
					out += "\t"
				}
			}
			indentLevel++
			out += " {\n"

		case gt == css.DeclarationGrammar:
			for i := 0; i < indentLevel; i++ {
				out += "\t"
			}
			out += string(data)
			out += ": "
			for _, val := range parser.Values() {
				switch val.TokenType {
				case css.CommaToken:
					out += string(val.Data) + " "
				default:
					out += string(val.Data)
				}
			}
			out += ";\n"
		case gt == css.EndAtRuleGrammar, gt == css.EndRulesetGrammar:
			indentLevel--
			out += "}\n\n"
		default:
			out += string(data) + "\n"
		}
	}

	return []byte(out), nil
}
