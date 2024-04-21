package _95midianFinder

// IsValid2 最小内存的方式
func IsValid2(s string) bool {
	m := map[byte]byte{
		'(': ')',
		'[': ']',
		'{': '}',
	}
	l := len(s)
	stack := make([]byte, l)
	top := 0
	for i := 0; i < l; i++ {
		switch s[i] {
		case '(', '[', '{':
			stack[top] = m[s[i]]
			top++
		case ')', ']', '}':
			if top > 0 && stack[top-1] == s[i] {
				top--
			} else {
				return false
			}
		}
	}
	return top == 0
}
