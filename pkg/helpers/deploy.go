package helpers

func DeployStatusRussianName(status bool) string {
	if status {
		return "открыт"
	} else {
		return "закрыт"
	}
}
