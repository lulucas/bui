<template>
  <div class="hello">
    <p>BUI port: {{port}}</p>
    <p>sum: 5+3={{ sum }}</p>
    <p>notify: {{ notify.State }}, updated at {{ notify.Time }}</p>
    <button @click="minimizeToTray">Minimize to tray</button>
    <button @click="showModal">Show Modal</button>
  </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
  name: 'Home',
  computed: {
    ...mapState(['notify'])
  },
  data() {
    return {
      port: window.location.port,
      sum: '?',
    }
  },
  mounted() {
    setTimeout(() => {
      this.$store.dispatch('sum', [5, 3])
    }, 1000);
  },
  methods: {
    sendNotify() {
      this.$store.dispatch('open_url', {url: "https://www.google.com"})
    },
    minimizeToTray() {
      this.$store.dispatch('minimize_to_tray')
    },
    showModal() {
      this.$store.dispatch('show_modal', {
        Width: 400,
        Height: 400,
        Url: `http://${window.location.host}/#/modal`,
      })
    }
  }
}
</script>

<style scoped>
</style>
