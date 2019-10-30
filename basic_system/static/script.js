let dmp = null;

document.addEventListener("DOMContentLoaded", async function() {
  dmp = new diff_match_patch();
  const doc = document.querySelector("#doc");

  let gotInitialDoc = false

  dmp = new diff_match_patch();

  ws = new WebSocket("ws://" + document.location.host + "/connect");

  ws.onopen = function(e) {
    console.log("OPEN WS");
  }
  ws.onclose = function(e) {
    console.log("CLOSE WS");
    ws = null;
  }
  ws.onmessage = function(e) {
    if (!gotInitialDoc) {
      doc.value = e.data;
      doc.prevValue = doc.value;
      gotInitialDoc = true;
    } else {
      getPatch(e);
    }
  }
  ws.onerror = function(e) {
    console.log("ERROR: " + e.data);
  }
  // On focus, save the value of the doc before any input
  doc.addEventListener("focus", () => doc.prevValue = doc.value)

  doc.addEventListener("input", (event) => {
    sendPatch(doc.prevValue, doc.value);
    doc.prevValue = doc.value;
  });
});


function getPatch(e) {
<<<<<<< HEAD
  let patches = dmp.patch_fromText(e.data);
  const cursorStart = doc.selectionStart; //the current cursor start position
  const cursorEnd = doc.selectionEnd; //the current cursor end position
  [doc.value] = dmp.patch_apply(patches, doc.value);
  if (cursorStart !== null) { // if there is a cursor position, update it.
    updateCursorPostition(doc.prevValue, doc.value, cursorStart, cursorEnd);
  }
  doc.prevValue = doc.value;
}

//sets the cursor position
function setCaretPosition(ctrl, start, end) {
  // Modern browsers
  if (ctrl.setSelectionRange) {
    ctrl.focus();
    ctrl.setSelectionRange(start, end);

  // IE8 and below
  } else if (ctrl.createTextRange) {
    var range = ctrl.createTextRange();
    range.collapse(true);
    range.moveEnd('character', start);
    range.moveStart('character', end);
    range.select();
  }
}

//updates the position of the cursor given the diff made by another user
function updateCursorPostition(prev, curr, cursorStart, cursorEnd) {
  let diff = dmp.diff_main(prev, curr);
  let start = cursorStart;
  let end = cursorEnd;
  let i = 0;
  let changePos = 0;
  for (; i < diff.length; i++) {
    let op = diff[i][0];
    let len = diff[i][1].length;
    if (op !== 0) {
      if (changePos < end) {
        end += op * (end - changePos < len ? end - changePos : len);
        if (changePos < start) {
          start += op * (start - changePos < len ? start - changePos : len);
        }
      }
      break;
    }
    changePos += len;
  }

  setCaretPosition(doc, start, end);
}

function sendPatch(prev, curr) {
  let patch = dmp.patch_toText(dmp.patch_make(prev, curr))
  ws.send(patch);
}
