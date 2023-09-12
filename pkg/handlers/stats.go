package handlers

import "math"

func getStats(notSeen, inProgress, cardsLearned int) (float64, float64, float64) {
	total := notSeen + inProgress + cardsLearned

	return percentage(notSeen, total), percentage(inProgress, total), percentage(cardsLearned, total)
}

func roundFloat(num float64, precision int) float64 {
	round := func(num float64) int {
		return int(num + math.Copysign(0.5, num))
	}
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func percentage(part int, total int) float64 {
	x := float64(part) / float64(total)
	return 100 * roundFloat(x, 4)
}
