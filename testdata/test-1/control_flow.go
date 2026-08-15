package main

func controlFlow() int {
	sum := 0

outer:
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j == 3 {
				continue outer
			}
			if i == 4 {
				break outer
			}
			sum += j
		}
	}

	for range 3 {
		sum++
	}

	for _, n := range []int{1, 2, 3} {
		sum += n
	}

	switch x := sum % 3; x {
	case 0:
		sum++
		fallthrough
	case 1:
		sum += 2
	default:
		sum += 3
	}

	var payload any = sum
	switch v := payload.(type) {
	case int:
		sum = v
	default:
		sum = 0
	}

	if sum < 0 {
		goto end
	}
	sum++
end:
	return sum
}
