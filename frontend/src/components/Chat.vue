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
import io from 'socket.io-client';
import axios from 'axios';
import store from '@/store';
import router from '@/router';

export default {
  setup() {
    const messages = ref([]);
    const input = ref('');
    const token = store.state.token;
    const socket = io('http://localhost:1323/ws', { auth: { token } });

    socket.on('message', (m) => messages.value.push(m));

    function send() {
      if (!input.value) return;
      socket.emit('message', { content: input.value, private: false });
      input.value = '';
    }

    onMounted(async () => {
      try {
        const res = await axios.get('http://localhost:1323/api/memory', {
          headers: { Authorization: `Bearer ${token}` }
        });
        messages.value = res.data;
      } catch {
        router.push('/login');
      }
    });

    return { messages, input, send };
  }
};