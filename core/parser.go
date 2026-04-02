package core

import (
	"bufio"
	"os"
	"strings"
)

func ParseFunctions(filepath string) []Function {
	file, err := os.Open(filepath)
	if err != nil {
		return []Function{}
	}
	defer file.Close()

	ext := strings.ToLower(filepath[strings.LastIndex(filepath, ".")+1:])

	var functions []Function
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		name, signature, ok := extractFunction(line, ext)
		if ok {
			functions = append(functions, Function{
				Name:      name,
				Signature: signature,
				Notes:     "",
				Author:    "",
				UpdatedAt: "",
			})
		}
	}

	return functions
}

func extractFunction(line string, ext string) (name string, signature string, ok bool) {
	switch ext {
	case "go":
		return extractGo(line)
	case "java":
		return extractJava(line)
	case "js", "ts":
		return extractJS(line)
	case "py":
		return extractPython(line)
	default:
		return "", "", false
	}
}

func extractGo(line string) (string, string, bool) {
	if !strings.HasPrefix(line, "func ") {
		return "", "", false
	}
	end := strings.Index(line, "{")
	if end == -1 {
		end = len(line)
	}
	signature := strings.TrimSpace(line[:end])
	name := signature
	parenStart := strings.Index(signature, "(")
	if parenStart != -1 {
		parts := strings.Fields(signature[:parenStart])
		if len(parts) >= 2 {
			name = parts[len(parts)-1]
		}
	}
	return name, signature, true
}

func extractJava(line string) (string, string, bool) {
	keywords := []string{"public ", "private ", "protected ", "static "}
	isMethod := false
	for _, kw := range keywords {
		if strings.Contains(line, kw) {
			isMethod = true
			break
		}
	}
	if !isMethod || !strings.Contains(line, "(") {
		return "", "", false
	}
	end := strings.Index(line, "{")
	if end == -1 {
		end = len(line)
	}
	signature := strings.TrimSpace(line[:end])
	parenStart := strings.Index(signature, "(")
	parts := strings.Fields(signature[:parenStart])
	if len(parts) == 0 {
		return "", "", false
	}
	name := parts[len(parts)-1]
	return name, signature, true
}

func extractJS(line string) (string, string, bool) {
	if strings.HasPrefix(line, "function ") {
		end := strings.Index(line, "{")
		if end == -1 {
			end = len(line)
		}
		signature := strings.TrimSpace(line[:end])
		parts := strings.Fields(signature)
		if len(parts) >= 2 {
			name := strings.Split(parts[1], "(")[0]
			return name, signature, true
		}
	}
	if strings.Contains(line, "= (") || strings.Contains(line, "= async (") || strings.Contains(line, "=> {") {
		end := strings.Index(line, "=>")
		if end == -1 {
			end = len(line)
		}
		signature := strings.TrimSpace(line[:end])
		parts := strings.Fields(signature)
		if len(parts) >= 1 {
			name := strings.TrimRight(parts[0], "=: ")
			if name != "const" && name != "let" && name != "var" {
				return name, signature, true
			}
			if len(parts) >= 2 {
				return parts[1], signature, true
			}
		}
	}
	return "", "", false
}

func extractPython(line string) (string, string, bool) {
	if !strings.HasPrefix(line, "def ") {
		return "", "", false
	}
	end := strings.Index(line, ":")
	if end == -1 {
		end = len(line)
	}
	signature := strings.TrimSpace(line[:end])
	parts := strings.Fields(signature)
	if len(parts) >= 2 {
		name := strings.Split(parts[1], "(")[0]
		return name, signature, true
	}
	return "", "", false
}
