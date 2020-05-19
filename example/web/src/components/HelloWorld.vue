<template>
  <div class="hello">
    <p>BUI port: {{port}}</p>
    <p>sum: 5+3={{ sum }}</p>
    <p>state: {{ state.State }}, updated at {{ state.Time }}</p>
    <button @click="minimizeToTray">Minimize to tray</button>
  </div>
</template>

<script>
const WebSocket = require('rpc-websockets').Client
const ws = new WebSocket(`ws://127.0.0.1:${window.location.port}/rpc`)

export default {
  name: 'HelloWorld',
  data() {
    return {
      port: window.location.port,
      sum: '?',
      state: 'idle',
    }
  },
  async mounted() {
    ws.on('open', () => {
      ws.call('sum', [5, 3]).then(result => {
        this.sum = result
      })

      ws.notify('open_url', {url: "https://www.google.com"})

      ws.subscribe('state_changed')

      ws.on('state_changed', state => {
        this.state = {
          State: state.State,
          Time: state.Time,
        }
        console.log(`state_changed ${state.State} ${state.Time}`)
      })

      setTimeout(() => {
        ws.unsubscribe('state_changed')
      }, 20000);
    })
  },
  methods: {
    minimizeToTray() {
      ws.notify('minimize_to_tray')
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
