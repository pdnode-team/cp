<script lang="ts">
  import pb from '$lib/pocketbase';
  import { onMount } from 'svelte';
  import { page } from '$app/state';

  const charId = $derived(page.params.id);

  let char = $state<any>(null);
  let isLoading = $state(true);
  let error = $state("");

  onMount(async () => {
    try {
      char = await pb.collection('characters').getOne(charId!,{
        expand: 'owner'
      });
    } catch (err: any) {
      error = "Character not found.";
    } finally {
      isLoading = false;
    }
  });
</script>
<div class="mx-auto max-w-6xl px-4 py-10">
  {#if isLoading}
    <div class="flex h-96 items-center justify-center">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>
  {:else if error}
    <div class="alert alert-error">{error}</div>
  {:else}
    <div class="flex flex-col gap-10 lg:flex-row">
      <div class="lg:w-1/3">
        <div class="sticky top-24 overflow-hidden rounded-3xl border border-base-200 bg-base-100 shadow-2xl">
          <figure class="aspect-3/4">
            <img 
              src={pb.files.getURL(char, char.images?.[0])} 
              alt={char.name} 
              class="h-full w-full object-cover object-top"
            />
          </figure>
          <div class="p-6">
            <h1 class="text-3xl font-black">{char.name}</h1>
            <p class="text-sm opacity-60 mt-1">Origin: {char.origin}</p>
          </div>
        </div>
      </div>

      <div class="flex-1">
        <div class="stats stats-vertical w-full bg-base-100 shadow-lg md:stats-horizontal border border-base-200">
          <div class="stat">
            <div class="stat-title text-xs uppercase font-bold tracking-widest">Creator</div>
            <div class="stat-value text-xl">{char.expand?.owner?.name ?? 'Unknown'}</div>
          </div>
          <div class="stat">
            <div class="stat-title text-xs uppercase font-bold tracking-widest">Created At</div>
            <div class="stat-value text-xl">{new Date(char.created).toLocaleDateString()}</div>
          </div>
        </div>

        <div class="mt-8 rounded-3xl bg-base-100 p-8 border border-base-200">
          <h2 class="mb-4 text-xl font-bold border-b pb-2">Description</h2>
          <div class="prose max-w-none text-base-content/80 leading-relaxed">
            {char.description}
          </div>
        </div>

        {#if char.tags_names}
          <div class="mt-6 flex flex-wrap gap-2">
            {#each [].concat(char.tags_names || []) as tag}
              <span class="badge badge-lg badge-outline">{tag}</span>
            {/each}
          </div>
        {/if}

        <div class="mt-12">
          <h3 class="mb-6 text-xl font-bold">Participated CPs</h3>
          <p class="text-sm opacity-50 italic">Feature coming soon...</p>
        </div>
      </div>
    </div>
  {/if}
</div>