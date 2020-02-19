<template>
  <div class="conversation">
    <div class="header">
      <h1>Conversation {{ $route.params.id }}</h1>
      <h2>{{ time }}</h2>
    </div>
    <div class="body" contenteditable="true">
    </div>
  </div>
</template>

<script>
export default {
  name: 'Coversation',
  data: () => ({
    time: 'Today at 9 P.M.',
  }),
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
