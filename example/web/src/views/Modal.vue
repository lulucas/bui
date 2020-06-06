<template>
  <div class="modal">
    <h1>This is modal window</h1>
    <button @click="onClose">close</button>
  </div>
</template>

<script>
export default {
  methods: {
    onClose() {
      this.$store.dispatch('close_modal')
    }
  }
}
</script>

<style scoped>
.modal {
  width: 400px;
  height: 300px;
  color: #000;
  background-color: cornflowerblue;
}
</style>
