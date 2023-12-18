/**
 * @Author: steven
 * @Description:
 * @File: similarity_test
 * @Version: 1.0.0
 * @Date: 22/05/23 08.44
 */

package algo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"testing"
)

type similarityTestSuite struct {
	suite.Suite
	m SimilarityManager
}

type similarityItemScoringConfig struct {
	Key      string       `yaml:"key"`
	Weight   float64      `yaml:"weight"`
	Matchers [][2]float64 `yaml:"matchers"`
}

const (
	cfg = `
simple_date:
  - key: "dob"
    weight: 0.25
    matchers:
      - [1, 0.4]
      - [2, 0.6]
      - [3, 1.0]
jaro_winkler:
  - key: "full_name"
    weight: 0.35
levenshtein:
  - key: "ktp"
    weight: 0.4
    matchers:
      - [0,1.0]
      - [1,0.8]
      - [2,0.6]
      - [3,0.4]
`
)

func (ts *similarityTestSuite) SetupTest() {
	var similarityConfig map[string][]similarityItemScoringConfig
	err := yaml.Unmarshal([]byte(cfg), &similarityConfig)
	fmt.Println(err)
	opts := make([]SimilarityOption, 0)
	weights := make(map[string]float64, 0)
	for algorithm, configs := range similarityConfig {
		for _, item := range configs {
			weights[item.Key] = item.Weight
			if item.Matchers == nil || len(item.Matchers) < 1 {
				continue
			}
			scoringMap := func(items [][2]float64) (rs map[int]float64) {
				rs = make(map[int]float64, 0)
				for _, v := range items {
					rs[int(v[0])] = v[1]
				}
				return
			}
			switch algorithm {
			case "levenshtein":
				opts = append(opts, WithSimilarityLevenshteinScoringMap(scoringMap(item.Matchers)))
			case "simple_date":
				opts = append(opts, WithSimilaritySimpleDateScoringMap(scoringMap(item.Matchers)))
			}
		}
	}
	if len(weights) > 0 {
		opts = append(opts, WithSimilarityWeightMap(weights))
	}
	ts.m = NewSimilarity(opts...)
}

func (ts *similarityTestSuite) TestSimpleDate() {
	type args struct {
		value string // input value
		ref   string // reference value
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Given valid string value then return valid simple date result",
			args: args{value: "1997-06-06T00:00:00", ref: "1997-06-06"},
			want: 1.0,
		},
	}
	for _, tt := range tests {
		ts.Run(tt.name, func() {
			actual := ts.m.SimpleDate(tt.args.value, tt.args.ref)
			assert.Equal(ts.T(), tt.want, actual)
		})
	}
}

func TestSimilaritySuite(t *testing.T) {
	suite.Run(t, new(similarityTestSuite))
}
