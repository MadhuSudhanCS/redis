package store

import (
	"math"
	"strconv"
	"strings"
)

//Holds Score-Member pair
type ScoreMember struct {
	Member string `json:"member"`
	//String type is required to save the value in json
	Score string `json:"score"`
}

//Holds Score-Member along with the knowledge of their ranks and order
type SortedScores struct {
	Scores  map[string]ScoreMember `json:"scores"`
	Members []string               `json:"members"`
	Ranks   map[string]string      `json:"ranks"`
}

func NewSortedScores() SortedScores {
	s := SortedScores{
		Scores:  make(map[string]ScoreMember),
		Members: []string{},
		Ranks:   make(map[string]string),
	}

	return s
}

//Ensure both the map and array is consistent
func (s *SortedScores) Add(scoreMember ScoreMember) {
	_, ok := s.Scores[scoreMember.Member]
	if !ok {
		s.Members = append(s.Members, scoreMember.Member)
	}

	s.Scores[scoreMember.Member] = scoreMember
}

//Array [1, 2, 3] will have Rank{0, -3} = 1 Rank{1, -2} = 2 Rank{2, -1} = 3
//BuildRank will store the Ranks and their index to Members array for easy access
func (s *SortedScores) BuildRank() {
	s.Ranks = make(map[string]string)
	length := len(s.Members)
	for idx := range s.Members {
		value := strconv.Itoa(idx)
		s.Ranks[strconv.Itoa(idx)] = value
		s.Ranks[strconv.Itoa(idx-(length))] = value
	}
}

func (s *SortedScores) Get(member string) (ScoreMember, bool) {
	sM, ok := s.Scores[member]
	return sM, ok
}

//Return number of elements (min, max)
//min, max are inclusive
func (s *SortedScores) Count(min, max string) int {
	if strings.EqualFold(min, MINSCORE) && strings.EqualFold(max, MAXSCORE) {
		return s.Len()
	}

	minScore := math.SmallestNonzeroFloat64
	maxScore := math.MaxFloat64
	count := 0

	if max != MAXSCORE {
		maxScore = ToScoreF(max)
	}

	if min != MINSCORE {
		minScore = ToScoreF(min)
	}

	for _, member := range s.Members {
		score := ToScoreF(s.Scores[member].Score)

		if score >= minScore && score <= maxScore {
			count = count + 1
		}

		if score > maxScore {
			break
		}
	}

	return count
}

//Returns slice of ScoreMember within the range
func (s *SortedScores) Range(start, stop int) []ScoreMember {
	startIdx, ok := s.Ranks[strconv.Itoa(start)]
	if !ok {
		return []ScoreMember{}
	}

	start, err := strconv.Atoi(startIdx)
	if err != nil {
		return []ScoreMember{}
	}

	if stop >= len(s.Members) {
		stop = len(s.Members) - 1
	}

	stopIdx, ok := s.Ranks[strconv.Itoa(stop)]
	if !ok {
		return []ScoreMember{}
	}

	stop, err = strconv.Atoi(stopIdx)
	if err != nil {
		return []ScoreMember{}
	}

	if start > stop {
		return []ScoreMember{}
	}

	scores := []ScoreMember{}
	for i := start; i <= stop; i++ {
		member := s.Members[i]
		scores = append(scores, s.Scores[member])
	}

	return scores
}

func (s SortedScores) Len() int {
	return len(s.Members)
}

func (s SortedScores) Less(prev, next int) bool {
	member := s.Members[prev]
	nextMember := s.Members[next]
	score := s.Scores[member].Score
	nextScore := s.Scores[nextMember].Score

	value, err := strconv.ParseFloat(score, 64)
	if err != nil {
		doMemberComp(member, nextMember)
	}

	nextValue, err := strconv.ParseFloat(nextScore, 64)
	if err != nil {
		doMemberComp(member, nextMember)
	}

	if value < nextValue {
		return true
	}

	if nextValue < value {
		return false
	}

	return doMemberComp(member, nextMember)
}

func (s SortedScores) Swap(i, j int) {
	s.Members[i], s.Members[j] = s.Members[j], s.Members[i]
}

func doMemberComp(member, nextMember string) bool {
	if strings.Compare(member, nextMember) <= 0 {
		return true
	}

	return false
}

func ToScoreF(score string) float64 {
	value, err := strconv.ParseFloat(score, 64)
	if err != nil {
		return 0.0
	}

	return value
}

func ToScoreS(score float64) string {
	return strconv.FormatFloat(score, 'f', 6, 64)
}
