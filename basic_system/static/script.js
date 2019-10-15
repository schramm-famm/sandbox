let dmp = null;
let id = '';

document.addEventListener("DOMContentLoaded", async function() {
  dmp = new diff_match_patch();
  const doc = document.querySelector("#doc");

  doc.value = await getState().catch(err => console.log(err));
  doc.prevValue = doc.value;
  // On focus, save the value of the doc before any input
  doc.addEventListener("focus", () => doc.prevValue = doc.value)

  doc.addEventListener("input", (event) => {
    sendPatch(doc.prevValue, doc.value);
    // Update previous value of the doc
    doc.prevValue = doc.value;
  });

  let client = new EventSource("/subscribe");
  client.onmessage = getPatch;
  client.onerror = (err) => {
    console.error("EventSource failed: ", err);
  }
});

// getPatch handles EventSource events and applies the patch to the doc
function getPatch(e) {
  let body = JSON.parse(e.data);
  // Patch doesn't get applied if it's the same as the last patch that the
  // client sent
  if (body.patch) {
    let patches = dmp.patch_fromText(body.patch);
    const prevPos = doc.selectionStart; //the current cursor position
    [doc.value] = dmp.patch_apply(patches, doc.value);
    console.log(prevPos);
    if(prevPos !== null){ // if there is a cursor position, update it.
      updateCursorPostition(doc.prevValue, doc.value, prevPos);
    }
    doc.prevValue = doc.value;
  } else if (body.id) {
    id = body.id;
  }
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

// sendPatch sends the patches between the prev text and curr text to the server
// to be applied to the global doc
function sendPatch(prev, curr) {
  let patch = dmp.patch_toText(dmp.patch_make(prev, curr))
  let body = { patch, id };

  const req = new XMLHttpRequest();
  const url = "/patch/";
  req.open("POST", url);
  req.setRequestHeader("Content-Type", "application/json");
  req.send(JSON.stringify(body));

  req.onreadystatechange = (event) => {
    if (req.readState === 4 && req.status !== 200) {
      console.log(req.response);
    }
  };
}

async function getState() {
  const req = new XMLHttpRequest();;
  const url = "/state/";
  req.open("GET", url);

  req.send();

  return new Promise((resolve, reject) => {
    req.onreadystatechange = (event) => {
      if (req.readyState === 4) {
        if (req.status === 200) {
          resolve(req.responseText);
        } else {
          reject(req.responseText);
        }
      }
    };
  });
}
