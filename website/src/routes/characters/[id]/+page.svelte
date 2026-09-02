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
<div class="mx-auto max-w-6xl px-4 py-6 sm:py-10">
  {#if isLoading}
    <div class="flex h-96 items-center justify-center">
      <span class="loading loading-spinner loading-lg text-primary"></span>
    </div>
  {:else if error}
    <div class="alert alert-error">{error}</div>
  {:else}
    <div class="flex flex-col gap-6 md:gap-10 lg:flex-row">
      <div class="lg:w-1/3">
        <div class="overflow-hidden rounded-2xl border border-base-200 bg-base-100 shadow-xl sm:rounded-3xl sm:shadow-2xl lg:sticky lg:top-24">
          <figure class="max-h-80 sm:max-h-96 lg:max-h-none aspect-square lg:aspect-3/4 bg-base-300">
            <img 
              src={pb.files.getURL(char, char.images?.[0])} 
              alt={char.name} 
              class="h-full w-full object-cover object-top"
            />
          </figure>
          <div class="p-4 sm:p-6">
            <h1 class="break-words text-2xl font-black sm:text-3xl">{char.name}</h1>
            {#if char.origin}
              <p class="mt-1 break-all text-sm opacity-60">Origin: {char.origin}</p>
            {/if}
          </div>
        </div>
      </div>

      <div class="flex-1">
        <div class="stats stats-horizontal w-full border border-base-200 bg-base-100 shadow-md sm:shadow-lg">
          <div class="stat p-4 sm:p-6">
            <div class="stat-title text-xs font-bold uppercase tracking-wider">Creator</div>
            <div class="stat-value truncate text-base sm:text-xl">{char.expand?.owner?.name ?? 'Unknown'}</div>
          </div>
          <div class="stat p-4 sm:p-6">
            <div class="stat-title text-xs font-bold uppercase tracking-wider">Created At</div>
            <div class="stat-value text-base sm:text-xl">{new Date(char.created).toLocaleDateString()}</div>
          </div>
        </div>

        <div class="mt-6 rounded-2xl border border-base-200 bg-base-100 p-5 sm:mt-8 sm:rounded-3xl sm:p-8">
          <h2 class="mb-3 border-b border-base-200 pb-2 text-lg font-bold sm:mb-4 sm:text-xl">Description</h2>
          <div class="prose max-w-none break-words text-sm leading-relaxed text-base-content/80 sm:text-base">
            {char.description}
          </div>
        </div>

        {#if char.tag_names}
          <div class="mt-6 flex flex-wrap gap-2">
            {#each [].concat(char.tag_names || []) as tag}
              <span class="badge badge-outline">{tag}</span>
            {/each}
          </div>
        {/if}

        <div class="mt-8 sm:mt-12">
          <h3 class="mb-4 text-lg font-bold sm:mb-6 sm:text-xl">Participated CPs</h3>
          <p class="text-sm italic opacity-50">Feature coming soon...</p>
        </div>
      </div>
    </div>
  {/if}
</div>