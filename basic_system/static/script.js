let dmp = null;

document.addEventListener("DOMContentLoaded", async function() {
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

  doc.addEventListener("input", (event) => {
    sendPatch(doc.prevValue, doc.value);
    doc.prevValue = doc.value;
  });
});


function getPatch(e) {
  let patches = dmp.patch_fromText(e.data);
  [doc.value] = dmp.patch_apply(patches, doc.value);
  doc.prevValue = doc.value;
}

function sendPatch(prev, curr) {
  let patch = dmp.patch_toText(dmp.patch_make(prev, curr))
  ws.send(patch);
}
