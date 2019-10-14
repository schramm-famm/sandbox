document.addEventListener("DOMContentLoaded", async function() {
  const doc = document.querySelector("#doc");

  let gotInitialDoc = false

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
      gotInitialDoc = true;
    } else {
      doc.value += e.data
    }
  }
  ws.onerror = function(e) {
    console.log("ERROR: " + e.data);
  }

  doc.addEventListener("input", (event) => {
    if (event.inputType === "insertText") {
      ws.send(event.data)
    } else if (event.inputType === "insertLineBreak") {
      ws.send("\n")
    }
  });
});
