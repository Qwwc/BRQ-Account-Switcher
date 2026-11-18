package platform

// Tests run on real machines that may have either Epic launcher installed.
// Start every test from "neither installed" so Epic detection never changes
// what an unrelated test sees; Epic tests opt in with withEpicEnv.
func init() {
	epicEnv = func(string) string { return "" }
}
