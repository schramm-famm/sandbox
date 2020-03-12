package cursor

type Update struct {
	CursorStartDelta int
	CursorEndDelta   int
	DocDelta         int
}

func processAdd(
	receiverStart,
	receiverEnd,
	senderStart,
	senderEnd int,
	update Update,
) (int, int) {
	if senderEnd < receiverEnd && update.CursorEndDelta != 0 {
		receiverEnd += update.DocDelta
		if senderEnd <= receiverStart {
			receiverStart += update.DocDelta
		}
	}
	return receiverStart, receiverEnd
}

func processUpdate(
	receiverStart,
	receiverEnd,
	senderStart,
	senderEnd int,
	update Update,
) (int, int) {
	var rangeStart, rangeEnd int
	if senderStart != senderEnd {
		rangeStart = senderStart
		rangeEnd = senderEnd
	} else if update.DocDelta > 0 {
		rangeStart = senderStart
		rangeEnd = senderStart + update.CursorStartDelta
	} else if update.CursorStartDelta == 0 {
		rangeStart = senderStart
		rangeEnd = senderStart - update.DocDelta
	} else {
		rangeStart = senderStart + update.CursorStartDelta
		rangeEnd = senderStart
	}

	if update.DocDelta > 0 && senderStart == senderEnd {
		return processAdd(receiverStart, receiverEnd, senderStart, senderEnd, update)
	}

	delta := rangeStart - rangeEnd

	if rangeEnd < receiverEnd {
		receiverEnd += delta
		if rangeEnd <= receiverStart {
			receiverStart += delta
		} else {
			if rangeStart < receiverStart {
				receiverStart = rangeStart
			}
		}
	} else {
		if rangeEnd == receiverEnd {
			receiverStart = rangeEnd + delta
			receiverEnd = rangeEnd + delta
		} else if rangeStart < receiverEnd {
			receiverEnd = rangeStart
			if rangeStart < receiverStart {
				receiverStart = rangeStart
			}
		}
	}

	if update.CursorStartDelta > 0 {
		addUpdate := Update{
			CursorStartDelta: update.CursorStartDelta,
			CursorEndDelta:   update.CursorStartDelta,
			DocDelta:         update.CursorStartDelta,
		}

		return processAdd(receiverStart, receiverEnd, senderStart, senderStart, addUpdate)
	}
	return receiverStart, receiverEnd
}
