<template>
  <div class="conversation">
    <div class="header">
      <h1>Conversation {{ $route.params.id }}</h1>
      <h2>{{ time }}</h2>
    </div>
    <div id="conversation-body" class="body" contenteditable="true"
    v-on:focus="storePrev" v-on:input="sendPatch">
    </div>
  </div>
</template>

<script>
import DiffMatchPatch from 'diff-match-patch';

const dmp = new DiffMatchPatch();

const data = () => ({
  time: 'Today at 9 P.M.',
  prevContent: '',
  ws: null,
  contentInitialized: false,
  convBody: null,
});

// sets the cursor position
function setCaretPosition(ctrl, start, end) {
  console.log(start, end);
  // Modern browsers
  if (ctrl.setSelectionRange) {
    ctrl.focus();
    ctrl.setSelectionRange(start, end);

  // IE8 and below
  } else if (ctrl.createTextRange) {
    const range = ctrl.createTextRange();
    range.collapse(true);
    range.moveEnd('character', start);
    range.moveStart('character', end);
    range.select();
  }
}

// updates the position of the cursor given the diff made by another user
function updateCursorPosition(prev, curr, cursorStart, cursorEnd) {
  const diff = dmp.diff_main(prev, curr);
  let start = cursorStart;
  let end = cursorEnd;
  let changePos = 0;
  for (let i = 0; i < diff.length; i += 1) {
    const op = diff[i][0];
    const len = diff[i][1].length;
    if (op !== 0) {
      if (changePos < end) {
        end += op * len;
        if (changePos < start) {
          start += op * len;
        }
      }
      break;
    }
    changePos += len;
  }

  setCaretPosition(this.convBody, start, end);
}

function getPatch(e) {
  this.convBody.focus();
  const patches = dmp.patch_fromText(e.data);
  const selectionRange = document.getSelection().getRangeAt(0);
  const cursorStart = selectionRange.startOffset; // the current cursor start position
  const cursorEnd = selectionRange.endOffset; // the current cursor end position
  [this.convBody.innerHTML] = dmp.patch_apply(patches, this.convBody.innerHTML);
  if (cursorStart !== null) { // if there is a cursor position, update it.
    this.updateCursorPosition(this.prevContent, this.convBody.innerHTML, cursorStart, cursorEnd);
  }
  this.prevContent = this.convBody.innerHTML;
}

function sendPatch() {
  const patch = dmp.patch_toText(dmp.patch_make(this.prevContent, this.convBody.innerHTML));
  this.ws.send(patch);
  this.prevContent = this.convBody.innerHTML;
}

function mounted() {
  this.$nextTick(() => {
    this.ws = new WebSocket('ws://localhost:8000/connect');
    this.convBody = this.$el.querySelector('#conversation-body');

    // Code that will run only after the
    // entire view has been rendered
    this.ws.onopen = () => {
      console.log('OPEN WS');
    };
    this.ws.onclose = () => {
      console.log('CLOSE WS');
      this.ws = null;
    };
    this.ws.onmessage = (e) => {
      if (!this.contentInitialized) {
        this.convBody.innerHTML = e.data;
        this.prevContent = this.convBody.innerHTML;
        this.contentInitialized = true;
      } else {
        this.getPatch(e);
      }
    };
    this.ws.onerror = (e) => {
      console.log(`ERROR: ${e.data}`);
    };
  });
}

function storePrev() {
  this.prevContent = this.convBody.innerHTML;
}

export default {
  name: 'Coversation',
  data,
  mounted,
  methods: {
    storePrev,
    sendPatch,
    getPatch,
    updateCursorPosition,
  },
};
</script>

<style scoped>
h1 {
  font-size: 12pt;
  font-weight: normal;
  font-style: normal;
  color: black;
  margin: 0px;
}

h2 {
  font-size: 8pt;
  font-weight: normal;
  font-style: normal;
  color: darkgrey;
  margin: 0px;
}

.conversation {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.header {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;
  height: 3em;
  background-color: white;
  border: solid lightgrey;
  border-width: 0px 0px thin;
  padding: 0px 1em;
}

.body {
  display: inline-block;
  color: black;
  font-size: 10pt;
  padding: 1em;
  height: 100%;
  overflow-y: scroll;
  word-wrap: break-word;
  word-break: break-all;
  white-space: pre-wrap;
}
</style>
