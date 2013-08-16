package mqtt

func validWillMessage(s String) bool {
    for _, r := range s {
        if r > 0x7f {
            return false
        }
    }
    return true
}
