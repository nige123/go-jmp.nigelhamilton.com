package template

import "regexp"

var tokenPattern = regexp.MustCompile(`\[-([a-z-]+)-\]`)

func Render(input string, params map[string]string) string {
    return tokenPattern.ReplaceAllStringFunc(input, func(token string) string {
        match := tokenPattern.FindStringSubmatch(token)
        if len(match) != 2 {
            return token
        }
        if value, ok := params[match[1]]; ok {
            return value
        }
        return token
    })
}
