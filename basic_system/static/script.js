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
  let patches = dmp.patch_fromText(e.data);
  const prevPos = doc.selectionStart; //the current cursor position
  [doc.value] = dmp.patch_apply(patches, doc.value);
  if (prevPos !== null) { // if there is a cursor position, update it.
    updateCursorPostition(doc.prevValue, doc.value, prevPos);
  }
  doc.prevValue = doc.value;
}

//sets the cursor position
function setCaretPosition(ctrl, pos) {
  // Modern browsers
  if (ctrl.setSelectionRange) {
    ctrl.focus();
    ctrl.setSelectionRange(pos, pos);

  // IE8 and below
  } else if (ctrl.createTextRange) {
    var range = ctrl.createTextRange();
    range.collapse(true);
    range.moveEnd('character', pos);
    range.moveStart('character', pos);
    range.select();
  }
}

//updates the position of the cursor given the diff made by another user
function updateCursorPostition(prev, curr, prevPos) {
  let diff = dmp.diff_main(prev, curr);
  let pos = prevPos;
  if ((((diff[0][1]).length) <= prevPos)
      || (prevPos === 0 && diff[0][0] === 1)){
    pos++;
  }
  setCaretPosition(doc, pos);
}

function sendPatch(prev, curr) {
  let patch = dmp.patch_toText(dmp.patch_make(prev, curr))
  ws.send(patch);
}
