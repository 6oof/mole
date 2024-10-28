package helpers

func ServiceNameModifier(service, pname string) string {
	return "mole-" + pname + "-" + service
}
