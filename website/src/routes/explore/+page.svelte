<script lang="ts">
	import pb from '$lib/pocketbase';

	type ViewMode = 'cps' | 'characters';
	let mode = $state<ViewMode>('cps');
	let items = $state<any[]>([]);
	let isLoading = $state(true);

	// 封装请求逻辑
	async function fetchData(currentMode: ViewMode) {
		isLoading = true;
		try {
			const collection = currentMode === 'cps' ? 'cps' : 'characters';
			const options = currentMode === 'cps' ? { expand: 'characters' } : {};

			items = await pb.collection(collection).getFullList({
				sort: '-created',
				...options
			});
		} catch (err) {
			console.error(err);
		} finally {
			isLoading = false;
		}
	}

	// 监听 mode 变化（Svelte 5 可以在监听到变量修改后触发函数）
	$effect(() => {
		fetchData(mode);
	});
</script>

<div class="mx-auto max-w-6xl px-6 py-16">
	<div class="mb-12 flex flex-col justify-between gap-6 md:flex-row md:items-center">
		<div>
			<h1 class="text-4xl font-bold tracking-tight md:text-5xl">Explore</h1>
			<div role="tablist" class="tabs-boxed mt-4 tabs">
				<button
					role="tab"
					class="tab {mode === 'cps' ? 'tab-active' : ''}"
					onclick={() => (mode = 'cps')}>CPs</button
				>
				<button
					role="tab"
					class="tab {mode === 'characters' ? 'tab-active' : ''}"
					onclick={() => (mode = 'characters')}>Characters</button
				>
			</div>
		</div>
		<a href="/create" class="btn btn-primary">Create New</a>
	</div>

	{#if isLoading}{:else}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each items as item}
				<div class="card border border-base-200 bg-base-100 shadow-xl">
					<figure class="h-48 bg-base-300">
						{#if item.images[0]}
							<img
							src={pb.files.getURL(item, item.images[0], { thumb: '400x300' })}
							alt={item.name}
							class="h-full w-full object-cover"
						/>
						{:else}
						<p>{item.name}</p>
						{/if}
						
					</figure>
					<div class="card-body p-5">
						<h2 class="card-title">
							{item.name}
							<div class="badge {mode === 'cps' ? 'badge-secondary' : 'badge-accent'}">
								{mode === 'cps' ? 'CP' : 'Character'}
							</div>
						</h2>
						{#if mode === 'cps'}
							<div class="flex gap-1">
								{#each item.expand?.characters || [] as c}
									<span class="badge badge-ghost text-xs">{c.name}</span>
								{/each}
							</div>
						{/if}
						<p class="line-clamp-2 text-sm opacity-70">{item.description}</p>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
