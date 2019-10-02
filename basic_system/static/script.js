document.addEventListener("DOMContentLoaded", async function() {
  const doc = document.querySelector("#doc");

  updateText(doc);

  doc.addEventListener("input", (event) => {
    console.log(event);

    if (event.inputType === "insertText") {
      sendSnippet(event.data);
    } else if (event.inputType === "insertLineBreak") {
      sendSnippet("\n");
    }
  });

  console.log("Test");
});

function sendSnippet(snippet) {
  let body = { snippet };
  let blob = new Blob([JSON.stringify(body)], { type: "application/json" });

  const req = new XMLHttpRequest();
  const url = "/snippet/";
  req.open("POST", url);

  req.send(blob);

  req.onreadystatechange = (event) => {
      console.log(req.responseText);
  };
}

async function getState() {
  const req = new XMLHttpRequest();;
  const url = "/state/";
  req.open("GET", url);

  req.send();

  return new Promise((resolve, reject) => {
    req.onreadystatechange = (event) => {
      if(req.readyState === 4 && req.status === 200) {
        console.log(req.responseText);
        resolve(req.responseText);
      }
    };
  });
}

async function updateText(doc) {
  doc.value = await getState();
  setTimeout(() => updateText(doc), 500)
}
