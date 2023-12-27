<script lang="ts" setup>
import { onBeforeMount, onMounted, ref } from 'vue';
import { nanoid } from 'nanoid';
import { CheckAdmin, GetHosts } from '../wailsjs/go/main/App';
import { Host } from '@/interfaces/interfaces';

const hosts = ref<Host[]>([]);

const refreshHosts = async () => {
  try {
    const gottenHosts = await GetHosts();
    if (!gottenHosts) return;
    console.log(gottenHosts);
    hosts.value = gottenHosts.map((h) => ({
      id: nanoid(),
      host: h.host,
      ip: h.ip,
      comment: h.comment,
    }));
  } catch (e) {
    console.error(e);
  }
};

const copyClipboard = (text: string) => {
  navigator.clipboard.writeText(text);
};

const toggleEditMode = (host: Host) => {
  host.editMode = !host.editMode;

  if (!host.editMode) {
    host.host = host.host.trim();
    host.ip = host.ip.trim();
    host.comment = host.comment?.trim();

    if (!host.host || !host.ip) {
      deleteHost(host.id);
    }
  }
};

const deleteHost = (host: string) => {
  hosts.value = hosts.value.filter((h) => h.id !== host);
};

const addHost = () => {
  hosts.value.push({
    id: nanoid(),
    host: 'localhost',
    ip: '127.0.0.1',
    editMode: true,
  });
};

onBeforeMount(async () => {
  try {
    await CheckAdmin();
  } catch (e) {
    console.error(e);
    alert('You need to run this app as admin!');
  }
});

onMounted(async () => {
  await refreshHosts();
});
</script>

<template>
  <div class="w-full h-full flex flex-col text-white p-[1rem]">
    <div class="grid select-none">
      <h1 class="text-2xl place-self-center font-semibold">Simple Hosts</h1>
      <button @click="addHost" class="place-self-end">Add</button>
    </div>
    <div>
      <table class="w-full table-fixed text-center text-white">
        <thead>
          <tr class="select-none">
            <th>Host</th>
            <th>IP</th>
            <th>Comment</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="host in hosts"
            :key="host.id"
            class="transition-all hover:bg-neutral-800 table-row justify-evenly w-full"
          >
            <td>
              <input
                v-if="host.editMode"
                v-model="host.host"
                class="w-full bg-neutral-900 text-center"
              />
              <p
                v-else
                @click="copyClipboard(host.host)"
                class="w-fit cursor-pointer self-center m-auto"
              >
                {{ host.host }}
              </p>
            </td>
            <td>
              <input
                v-if="host.editMode"
                v-model="host.ip"
                class="w-full bg-neutral-900 text-center"
              />
              <p
                v-else
                @click="copyClipboard(host.ip)"
                class="w-fit cursor-pointer self-center m-auto"
              >
                {{ host.ip }}
              </p>
            </td>
            <td>
              <input
                v-if="host.editMode"
                v-model="host.comment"
                class="w-full bg-neutral-900 text-center"
              />
              <p v-else class="w-fit self-center m-auto">
                {{ host.comment }}
              </p>
            </td>
            <td class="flex justify-center gap-x-2 select-none">
              <button @click="toggleEditMode(host)">
                {{ host.editMode ? 'Save' : 'Edit' }}
              </button>
              <button @click="deleteHost(host.id)">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.grid {
  @apply justify-center;
  grid-template-areas: '. text button';
  grid-template-columns: 1fr auto 1fr;

  h1 {
    grid-area: text;
  }

  button {
    grid-area: button;
  }
}
</style>
