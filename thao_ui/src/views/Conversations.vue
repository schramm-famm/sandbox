<template>
  <div class="conversations">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'Home',
};
</script>

<style>
.conversations {
  flex: 1;
  display: flex;
  flex-direction: column;
}
</style>
