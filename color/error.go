package color

type Error string

func (e Error) Error() string { return string(e) }

const ExpectedHash = Error("Parse Error: expected # as first character for color reference")
const InvalidLength = Error("Parse Error: color reference does not have a valid length")
const InvalidPattern = Error("Parse Error: color reference does not match pattern #rrggbbaa or #rrggbb")
