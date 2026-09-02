<script lang="ts">
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	// 1. 从 URL 获取 ID
	const cpId = $derived(page.params.id);

	let cp = $state<any>(null);
	let isLoading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			// 展开 characters 和 owner
			cp = await pb.collection('cps').getOne(cpId!, {
				expand: 'characters,owner'
			});
		} catch (err: any) {
			console.error(err);
			error = 'CP not found or connection failed.';
		} finally {
			isLoading = false;
		}
	});

	let selectedImg = $state('');
	let galleryModal = $state() as HTMLDialogElement;

	function openImage(img: string) {
		selectedImg = pb.files.getURL(cp, img);
		galleryModal.showModal();
	}

	function isSafeHttpUrl(url?: string): boolean {
		if (!url) return false;
		return /^https?:\/\//i.test(url.trim());
	}
</script>

<div class="mx-auto max-w-5xl px-4 py-6 sm:py-10">
	{#if isLoading}
		<div class="flex h-64 items-center justify-center">
			<span class="loading loading-lg loading-spinner text-primary"></span>
		</div>
	{:else if error}
		<div class="alert alert-error">
			<span>{error}</span>
			<a href="/explore" class="btn btn-sm">Back to Explore</a>
		</div>
	{:else}
		<div class="mb-8 text-center sm:mb-10">
			<h1 class="break-words text-3xl font-black italic tracking-tighter text-primary sm:text-4xl md:text-5xl">
				{cp.name}
			</h1>
			<div class="mt-4 flex flex-wrap justify-center gap-2">
				{#each [].concat(cp.tag_names || []) as tag}
					{#if tag}
						<span class="badge badge-outline">{tag}</span>
					{/if}
				{/each}
			</div>
		</div>

		<div class="flex flex-col items-center justify-center gap-6 md:flex-row md:gap-8">
			{#each cp.expand?.characters || [] as char}
				<div class="card w-full max-w-sm border border-base-200 bg-base-100 shadow-xl transition-all hover:bg-base-200/50">
					<a href="/characters/{char.id}" class="block overflow-hidden rounded-t-2xl">
						<figure class="h-48 bg-base-300 sm:h-64">
							{#if char.images?.[0]}
								<img
									src={pb.files.getURL(char, char.images?.[0])}
									alt={char.name}
									class="h-full w-full object-cover object-top transition-transform hover:scale-105"
								/>
							{:else}
								<p>{char.name}</p>
							{/if}
						</figure>
					</a>
					<div class="card-body p-4 sm:p-6">
						<h3 class="card-title text-xl sm:text-2xl break-words">
							<a href="/characters/{char.id}" class="transition-colors hover:text-primary">{char.name}</a>
						</h3>
						{#if char.origin}
							{#if isSafeHttpUrl(char.origin)}
								<a class="text-sm opacity-70 link break-all" href="{char.origin}" target="_blank" rel="noopener noreferrer">{char.origin}</a>
							{:else}
								<span class="text-sm opacity-70 break-all">{char.origin}</span>
							{/if}
						{/if}
					</div>
				</div>
			{/each}
		</div>

		<div class="mt-10 rounded-2xl bg-base-200 p-5 shadow-inner sm:p-8 md:mt-16 md:rounded-3xl">
			<h2 class="mb-4 text-xl font-bold sm:mb-6 sm:text-2xl">Description</h2>
			<div class="prose max-w-none break-words text-base leading-relaxed sm:text-lg">
				{cp.description}
			</div>

			<div class="divider my-4"></div>

			<p class="text-sm text-base-content/60">
				Author:
				<span class="font-medium text-base-content">
					{cp.expand?.owner?.name ?? 'Anonymous'}
				</span>
			</p>
		</div>

		<div class="mt-10 md:mt-16">
			<h3 class="mb-4 text-xl font-bold sm:mb-6">Gallery</h3>
			<div class="grid grid-cols-2 gap-3 sm:gap-4 md:grid-cols-4">
				{#each cp.images as img}
					<button
						type="button"
						onclick={() => openImage(img)}
						class="group relative aspect-square w-full cursor-zoom-in overflow-hidden rounded-xl border-none bg-base-300 p-0"
						aria-label="View larger image"
					>
						<img
							src={pb.files.getURL(cp, img, { thumb: '400x400' })}
							class="h-full w-full object-cover transition-transform group-hover:scale-110"
							alt="Gallery Preview"
						/>
					</button>
				{/each}
			</div>
		</div>
	{/if}
</div>

<dialog bind:this={galleryModal} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box max-w-5xl overflow-hidden bg-transparent p-0 shadow-none">
		<form method="dialog">
			<button class="btn absolute top-2 right-2 z-10 btn-circle bg-base-100/50 btn-ghost btn-sm"
				>✕</button
			>
		</form>
		{#if selectedImg}
			<img
				src={selectedImg}
				class="mx-auto h-auto max-h-[90vh] w-full object-contain"
				alt="Full size"
			/>
		{/if}
	</div>
	<form method="dialog" class="modal-backdrop bg-black/80">
		<button>close</button>
	</form>
</dialog>
