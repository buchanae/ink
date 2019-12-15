package tess

import (
	"log"
	"testing"
)

func checkItems(t *testing.T, l *list, expect ...*listItem) {

	var actual []*listItem
	it := l.Iter()
	for it.Next() {
		item := it.Item()
		actual = append(actual, item)
	}

	if len(actual) != len(expect) {
		t.Log(actual, expect)
		t.Errorf("expect %d items, got %d", len(expect), len(actual))
		return
	}

	for i := range actual {
		if actual[i] != expect[i] {
			t.Errorf("expected %p at position %d, got %p",
				expect[i], i, actual[i])
		}
	}
}

func TestList(t *testing.T) {
	l := newList()

	var items []*listItem
	for i := 0; i < 5; i++ {
		items = append(items, &listItem{})
	}

	t.Log("append")
	l.Append(items[0])

	checkItems(t, l, items[0])

	t.Log("append")
	l.Append(items[1])
	checkItems(t, l, items[0], items[1])

	t.Log("insert before")
	l.InsertBefore(items[2], items[0])
	checkItems(t, l, items[0], items[1], items[2])

	t.Log("insert before")
	l.InsertBefore(items[3], items[1])
	checkItems(t, l, items[0], items[3], items[1], items[2])

	t.Log("delete")
	l.Delete(items[0])
	checkItems(t, l, items[3], items[1], items[2])

	t.Log("insert after")
	l.InsertAfter(items[4], items[3])
	checkItems(t, l, items[3], items[4], items[1], items[2])

	t.Log("delete 4")
	l.Delete(items[3])
	l.Delete(items[4])
	l.Delete(items[1])
	l.Delete(items[2])

	checkItems(t, l)
}

func TestListIter(t *testing.T) {
	l := newList()

	var items []*listItem
	for i := 0; i < 5; i++ {
		item := &listItem{}
		l.Append(item)
		items = append(items, item)
	}

	count := 0
	it := l.Iter()
	for it.Next() {
		count++
		log.Println(it.Item())
		//l.Delete(it.Item())
	}

	if count != 5 {
		t.Errorf("expected 5 iterations, got %d", count)
	}
}

func TestListIterDelete(t *testing.T) {
	l := newList()

	var items []*listItem
	for i := 0; i < 5; i++ {
		item := &listItem{}
		l.Append(item)
		items = append(items, item)
	}

	count := 0
	it := l.Iter()
	for it.Next() {
		expect := items[count]
		if it.Item() != expect {
			t.Errorf("unexpected item at position %d", count)
		}

		count++
		log.Println(it.Item())
		l.Delete(it.Item())
	}

	if count != 5 {
		t.Errorf("expected 5 iterations, got %d", count)
	}
}

func TestIterEmptyList(t *testing.T) {
	l := newList()
	it := l.Iter()
	count := 0
	for it.Next() {
		count++
		log.Println(it.Item())
	}

	if count != 0 {
		t.Errorf("expected zero iterations, got %d", count)
	}
}

func TestIterOneList(t *testing.T) {
	l := newList()
	item := &listItem{}
	l.Append(item)

	it := l.Iter()
	count := 0
	for it.Next() {
		count++
		log.Println(it.Item())
	}

	if count != 1 {
		t.Errorf("expected one iteration, got %d", count)
	}
}
