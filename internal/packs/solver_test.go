package packs

import "testing"

var defaultSizes = []int{250, 500, 1000, 2000, 5000}

func TestSolve_Examples(t *testing.T) {
	tests := []struct {
		name     string
		items    int
		shipped  int
		expected map[int]int
	}{
		{"1 item", 1, 250, map[int]int{250: 1}},
		{"250 items", 250, 250, map[int]int{250: 1}},
		{"251 items", 251, 500, map[int]int{500: 1}},
		{"501 items", 501, 750, map[int]int{500: 1, 250: 1}},
		{"12001 items", 12001, 12250, map[int]int{5000: 2, 2000: 1, 250: 1}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Solve(tc.items, defaultSizes)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.ItemsShipped != tc.shipped {
				t.Fatalf("expected shipped=%d got=%d", tc.shipped, res.ItemsShipped)
			}

			// Check expected packs
			for size, wantCount := range tc.expected {
				got := res.Packs[size]
				if got != wantCount {
					t.Fatalf("expected pack %d count=%d got=%d", size, wantCount, got)
				}
			}

			// Ensure no extra unexpected packs
			for size := range res.Packs {
				if _, ok := tc.expected[size]; !ok {
					t.Fatalf("unexpected pack size in result: %d", size)
				}
			}
		})
	}
}

func TestSolve_Boundary4999(t *testing.T) {
	res, err := Solve(4999, defaultSizes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ItemsShipped != 5000 {
		t.Fatalf("expected shipped=5000 got=%d", res.ItemsShipped)
	}
	if res.Packs[5000] != 1 || len(res.Packs) != 1 {
		t.Fatalf("expected packs {5000:1} got=%v", res.Packs)
	}
}

func TestSolve_InvalidInput(t *testing.T) {
	if _, err := Solve(0, defaultSizes); err == nil {
		t.Fatal("expected error for items=0")
	}
	if _, err := Solve(10, []int{}); err == nil {
		t.Fatal("expected error for empty pack sizes")
	}
	if _, err := Solve(10, []int{-1, 250}); err == nil {
		t.Fatal("expected error for negative pack size")
	}
}
