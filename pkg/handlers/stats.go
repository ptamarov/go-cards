package handlers

import "math"

func getStats(notSeen, inProgress, cardsLearned int) (float64, float64, float64) {
	total := notSeen + inProgress + cardsLearned
	return toPercentage(notSeen, total), toPercentage(inProgress, total), toPercentage(cardsLearned, total)
}

func toPercentage(part int, total int) float64 {
	x := float64(part) / float64(total)
	return math.Round(10000*x) / 100
}
