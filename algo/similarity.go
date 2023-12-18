/**
 * @Author: steven
 * @Description:
 * @File: similarity
 * @Date: 29/09/23 10.50
 */

package algo

import (
	"github.com/hbollon/go-edlib"
	"regexp"
	"time"
)

type SimilarityManager interface {
	SimpleDate(value, ref string) (score float64)
	Levenshtein(value, ref string) (score float64)
	JaroWinkler(value, ref string) (score float64)
	Weight(key string) float64
}

type SimilarityOption interface {
	apply(m *similarity)
}

type similarity struct {
	weights     map[string]float64
	levenshtein map[int]float64
	simpleDate  map[int]float64
}

type similarityOption func(m *similarity)

func (o similarityOption) apply(m *similarity) {
	o(m)
}

func WithSimilarityWeightMap(weights map[string]float64) SimilarityOption {
	return similarityOption(func(m *similarity) {
		m.weights = weights
	})
}

func WithSimilarityLevenshteinScoringMap(scoreMap map[int]float64) SimilarityOption {
	return similarityOption(func(m *similarity) {
		m.levenshtein = scoreMap
	})
}

func WithSimilaritySimpleDateScoringMap(scoreMap map[int]float64) SimilarityOption {
	return similarityOption(func(m *similarity) {
		m.simpleDate = scoreMap
	})
}

// cleanup unnecessary characters from date value
var rxDateValue = regexp.MustCompile("([0-9]{4}-[0-9]{2}-[0-9]{2})(.*)")

// both value & ref should be in yyyy-mm-dd format
func (s *similarity) simpleDateSimilarity(value, ref string) int {
	// parse and clean value/ref
	valueRsArr := rxDateValue.FindStringSubmatch(value)
	refRsArr := rxDateValue.FindStringSubmatch(ref)
	if len(valueRsArr) < 2 || len(refRsArr) < 2 {
		return 0
	}
	value = valueRsArr[1]
	ref = refRsArr[1]
	valueDate, errV := time.Parse("2006-01-02", value)
	if errV != nil {
		return 0
	}
	refDate, errR := time.Parse("2006-01-02", ref)
	if errR != nil {
		return 0
	}
	similarPart := 0
	if valueDate.Year() == refDate.Year() {
		similarPart++
	}
	if valueDate.Month() == refDate.Month() {
		similarPart++
	}
	if valueDate.Day() == refDate.Day() {
		similarPart++
	}
	return similarPart
}

func (s *similarity) Weight(key string) float64 {
	if v, ok := s.weights[key]; ok {
		return v
	}
	return 0.0
}

func (s *similarity) SimpleDate(value, ref string) (score float64) {
	diffCount := s.simpleDateSimilarity(value, ref)
	if v, ok := s.simpleDate[diffCount]; ok {
		return v
	}
	return 0.0
}

func (s *similarity) Levenshtein(value, ref string) (score float64) {
	replacementCount := edlib.LevenshteinDistance(value, ref)
	if v, ok := s.levenshtein[replacementCount]; ok {
		return v
	}
	return 0.0
}

func (s *similarity) JaroWinkler(value, ref string) (score float64) {
	result := edlib.JaroWinklerSimilarity(value, ref)
	return float64(result)
}

func NewSimilarity(opts ...SimilarityOption) SimilarityManager {
	m := &similarity{
		levenshtein: make(map[int]float64, 0),
	}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}
