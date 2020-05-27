<template>
<div id="app">
  <div>
    <router-view />
  </div>
</div>
</template>

<script>
export default {
    name: 'App'
}
</script>

<style lang="less">

@import url("pages/common.less");

#app {
    margin: 0;
    padding: 0;
}
</style>
