package main

func getDoc(usg, exp, opt string) string {
	opts := ""
	if len(opt) > 0 {
		opts = "\nOptions:\n" + opt + "\n"
	}

	return "\n" +
		"Usage:  " + usg + "\n" +
		"\n" +
		"        " + exp + "\n" + opts
}
