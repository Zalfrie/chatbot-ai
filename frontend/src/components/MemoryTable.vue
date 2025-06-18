<template>
  <v-container>
    <v-card>
      <v-card-title>Memory</v-card-title>
      <v-card-text>
        <table id="memory-table" class="display" style="width:100%"></table>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
import { onMounted } from 'vue';
import $ from 'jquery';
import 'datatables.net';
import axios from 'axios';
import store from '@/store';
import router from '@/router';

export default {
  setup() {
    const token = store.state.token;

    onMounted(() => {
      if (!token) router.push('/login');

      $('#memory-table').DataTable({
        ajax: {
          url: 'http://localhost:1323/api/memory',
          headers: { Authorization: `Bearer ${token}` },
          dataSrc: ''
        },
        columns: [
          { title: 'ID', data: 'id' },
          { title: 'User ID', data: 'user_id' },
          { title: 'Content', data: 'content' },
          { title: 'Created At', data: 'created_at' },
          { title: 'Actions', data: null, render: (d) => `<button class="delete" data-id="${d.id}">Delete</button>` }
        ],
        initComplete: function () {
          $('#memory-table').on('click', '.delete', async function () {
            const id = $(this).data('id');
            await axios.delete(`http://localhost:1323/api/memory/${id}`, { headers: { Authorization: `Bearer ${token}` } });
            $('#memory-table').DataTable().ajax.reload();
          });
        }
      });
    });
  }
};
</script>