package nb_challenges

import "node-backend/util/auth"

// GenerationFunction is a function that generates a new challenge
// the JS file contains a complete function that will be called by the client
var GenerationFunction func() (tk string, result string, js string, attach interface{}) = func() (tk string, result string, js string, attach interface{}) {
	tk = auth.GenerateToken(100)

	// Create complete function
	js = "function complete() {"
	js += "return '" + tk + "';"
	js += "}"

	return tk, tk, js, ""
}
