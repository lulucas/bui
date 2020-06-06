import Vue from 'vue'
import Vuex from 'vuex'

const WebSocket = require('rpc-websockets').Client
const ws = new WebSocket(`ws://127.0.0.1:${window.location.port}/rpc`)

Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    notify: {
      State: '',
      Time: '',
    },
  },
  mutations: {
    setNotify(state, notify) {
      state.notify = notify
    }
  },
  actions: {
    sum(context, params) {
      return ws.call('sum', params)    
    },
    minimize_to_tray() {
      ws.notify('minimize_to_tray')
    },
    show_modal(context, params) {
      ws.notify('show_modal', params)
    },
    close_modal() {
      ws.notify('close_modal')
    }
  },
  modules: {
  }
})

ws.on('open', () => {
  ws.subscribe('state_changed')
  ws.on('state_changed', event => {
    store.commit('setNotify', event)
  })

  setTimeout(() => {
    ws.unsubscribe('state_changed')
  }, 20000);
})

export default store
