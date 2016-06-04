package store

import (
	"sort"
	"testing"
)

var (
	score1 = ScoreMember{
		Member: "GO",
		Score:  "1.5",
	}

	score2 = ScoreMember{
		Member: "C",
		Score:  "2",
	}

	score3 = ScoreMember{
		Member: "Java",
		Score:  "2.1",
	}

	score4 = ScoreMember{
		Member: "Python",
		Score:  "2.1",
	}
)

func TestSortedScore(t *testing.T) {
	expectedScores := []ScoreMember{score1, score2, score3, score4}

	s := SortedScores{
		Scores:  make(map[string]ScoreMember),
		Members: []string{},
	}

	s.Add(score1)
	s.Add(score4)
	s.Add(score2)
	s.Add(score3)
	s.Add(score3)
	s.Add(score3)

	sort.Sort(s)

	for _, sM := range expectedScores {
		scoreMem, ok := s.Get(sM.Member)
		if !ok {
			t.Fatalf("expected member: %s to be present", sM.Member)
		}

		if scoreMem != sM {
			t.Fatalf("Expected: %v found: %v", sM, scoreMem)
		}
	}
}

func TestScoreCount(t *testing.T) {
	s := SortedScores{
		Scores:  make(map[string]ScoreMember),
		Members: []string{},
	}

	s.Add(score1)
	s.Add(score4)
	s.Add(score2)
	s.Add(score3)

	sort.Sort(s)

	if count := s.Count(MINSCORE, MAXSCORE); count != 4 {
		t.Fatalf("Expected: %d found: %d", 4, count)
	}

	if count := s.Count(MINSCORE, "2.0"); count != 2 {
		t.Fatalf("Expected: %d found: %d", 2, count)
	}

	if count := s.Count("-1.0", "2.0"); count != 2 {
		t.Fatalf("Expected: %d found: %d", 2, count)
	}

	if count := s.Count("-1.0", MAXSCORE); count != 4 {
		t.Fatalf("Expected: %d found: %d", 4, count)
	}

	if count := s.Count("2.0", "2.2"); count != 3 {
		t.Fatalf("Expected: %d found: %d", 3, count)
	}

	if count := s.Count("2", "-2"); count != 0 {
		t.Fatalf("Expected: %d found: %d", 0, count)
	}
}

func TestScoreRange(t *testing.T) {
	s := SortedScores{
		Scores:  make(map[string]ScoreMember),
		Members: []string{},
	}

	s.Add(score1)
	s.Add(score4)
	s.Add(score2)
	s.Add(score3)

	sort.Sort(s)
	s.BuildRank()

	if values := s.Range(-2, -5); len(values) != 0 {
		t.Fatalf("Expected: %d found: %v", 4, values)
	}

	if values := s.Range(0, -1); len(values) != 4 {
		t.Fatalf("Expected: %d found: %v", 4, values)
	}

	if values := s.Range(0, -2); len(values) != 3 {
		t.Fatalf("Expected: %d found: %v", 4, values)
	}

	if values := s.Range(-2, -1); len(values) != 2 {
		t.Fatalf("Expected: %d found: %v", 4, values)
	}

	if values := s.Range(0, -100); len(values) != 0 {
		t.Fatalf("Expected: %d found: %v", 4, values)
	}

	if values := s.Range(-2, -3); len(values) != 0 {
		t.Fatalf("Expected: %d found: %v", 4, values)
	}

	if values := s.Range(0, 3); len(values) != 4 {
		t.Fatalf("Expected: %d found: %v", 4, values)
	}

	if values := s.Range(0, 5); len(values) != 4 {
		t.Fatalf("Expected: %d found: %v", 4, values)
	}

	if values := s.Range(100, 5); len(values) != 0 {
		t.Fatalf("Expected: %d found: %v", 4, values)
	}
}
