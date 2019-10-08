let dmp = null;
let prevPatch = "";

document.addEventListener("DOMContentLoaded", async function() {
  dmp = new diff_match_patch();
  const doc = document.querySelector("#doc");

  doc.value = await getState().catch(err => console.log(err));

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
  if (body.patch && body.patch != prevPatch) {
    let patches = dmp.patch_fromText(body.patch);
    [doc.value] = dmp.patch_apply(patches, doc.value);
    doc.prevValue = doc.value;
    prevPatch = "";
  }
}

// sendPatch sends the patches between the prev text and curr text to the server
// to be applied to the global doc
function sendPatch(prev, curr) {
  let patch = dmp.patch_toText(dmp.patch_make(prev, curr))
  let body = { patch };

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

  prevPatch = patch;
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
