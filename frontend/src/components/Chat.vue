<template>
  <v-container>
    <v-card>
      <v-card-title>Chat Room</v-card-title>
      <v-divider />
      <v-card-text>
        <div v-for="msg in messages" :key="msg.id">
          <strong>{{ msg.user_id }}:</strong> {{ msg.content }}
        </div>
      </v-card-text>
      <v-text-field
        v-model="input"
        @keyup.enter="send"
        label="Type a message"
      />
    </v-card>
  </v-container>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import store from '@/store';
import router from '@/router';

export default {
  name: 'ChatRoom',
  setup() {
    const messages = ref([]);
    const input = ref('');
    const token = store.state.token;

    // Establish a native WebSocket connection to Echo backend
    let socket;
    onMounted(async () => {
      // Verify authentication
      try {
        await axios.get('http://localhost:1323/api/memory', {
          headers: { Authorization: `Bearer ${token}` }
        });
      } catch {
        router.push('/login');
        return;
      }

      // Connect to /api/ws endpoint
      socket = new WebSocket(`ws://localhost:1323/api/ws`);

      socket.addEventListener('open', () => {
        console.log('WebSocket connected');
      });

      socket.addEventListener('message', (event) => {
        const msg = JSON.parse(event.data);
        messages.value.push(msg);
      });

      socket.addEventListener('close', () => {
        console.log('WebSocket disconnected');
      });
    });

    function send() {
      if (!input.value || socket.readyState !== WebSocket.OPEN) return;
      const payload = { content: input.value, private: false };
      socket.send(JSON.stringify(payload));
      input.value = '';
    }

    return { messages, input, send };
  }
};
</script>