package tess

type listIter struct {
	start, cur *listItem
}

func (li *listIter) Next() bool {
	var next *listItem
	if li.cur == nil {
		next = li.start.next
	} else {
		next = li.cur.next
	}

	if next == li.start {
		return false
	}
	li.cur = next
	return true
}

func (li *listIter) Item() *listItem {
	return li.cur
}

type listItem struct {
	*ecItem
	prev, next *listItem
}

type list struct {
	head *listItem
}

func newList() *list {
	head := &listItem{}
	head.prev = head
	head.next = head
	return &list{head}
}

func (l *list) Iter() *listIter {
	return &listIter{start: l.head}
}

func (l *list) Append(item *listItem) {
	tail := l.head.prev
	tail.next = item
	item.prev = tail

	l.head.prev = item
	item.next = l.head
}

func (l *list) Delete(item *listItem) {
	if l.head == item {
		if item.next == item {
			// last item
			l.head = nil
		} else {
			l.head = l.head.next
		}
	}

	item.prev.next = item.next
	item.next.prev = item.prev
	//item.next = nil
	//item.prev = nil
}

func (l *list) Prev(item *listItem) *listItem {
	// TODO still breaks if list is empty or one,
	//      since prev will point to itself
	if item.prev == l.head {
		return item.prev.prev
	}
	return item.prev
}

func (l *list) Next(item *listItem) *listItem {
	// TODO still breaks if list is empty or one,
	//      since next will point to itself
	if item.next == l.head {
		return item.next.next
	}
	return item.next
}

func (l *list) InsertBefore(item, index *listItem) {
	if index.prev == l.head {
		index = l.head
	}

	item.prev = index.prev
	index.prev.next = item
	item.next = index
	index.prev = item
}

func (l *list) InsertAfter(item, index *listItem) {
	item.next = index.next
	index.next.prev = item
	item.prev = index
	index.next = item
}
