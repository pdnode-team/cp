<script lang="ts">
	import { goto } from '$app/navigation';
	import pb from '$lib/pocketbase';

	type ViewMode = 'cps' | 'characters';
	let mode = $state<ViewMode>('cps');
	let items = $state<any[]>([]);
	let isLoading = $state(true);

	let likedTargetIds = $state<Set<string>>(new Set());


	const reloadLikes = async () => {

		if (!pb.authStore.record) return;

		const likes = await pb.collection('likes').getFullList({
			filter: `user = "${pb.authStore.record?.id}"`
		})

		likedTargetIds = new Set(likes.map(l => l.target_id));
	};

	reloadLikes();

	async function toggleLike(item: any) {
		if (!pb.authStore.record) { alert("Please login first"); return}
		const result = await fetch(`/api/${item.collectionName}/${item.id}/toggle-like`, {
			method: 'POST',
			headers: {
				Authorization: 'Bearer ' + pb.authStore.token
			}
		});
		if (!result.ok) {
			alert('Something went wrong');
			return;
		}

		await result.json();
		reloadLikes();
	}

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

<div class="mx-auto max-w-6xl px-4 py-8 sm:px-6 sm:py-12 md:py-16">
	<div class="mb-8 flex flex-col justify-between gap-4 sm:gap-6 md:mb-12 md:flex-row md:items-center">
		<div>
			<h1 class="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">Explore</h1>
			<div role="tablist" class="tabs-boxed mt-4 tabs w-fit">
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
		<a href="/create" class="btn btn-primary self-start md:self-auto">Create New</a>
	</div>

	{#if isLoading}{:else}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each items as item}
				<div
					role="link"
					tabindex="0"
					onkeydown={(e) => (e.key === 'Enter' || e.key === ' ') && goto(`/${mode}/${item.id}`)}
					aria-label={`View ${item.name}'s Details'`}
					class="card cursor-pointer border border-base-200 bg-base-100 shadow-xl transition-all hover:bg-base-300 active:scale-95"
					onclick={() => goto(`/${mode}/${item.id}`)}
				>
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
					<div class="card-body flex-row items-center justify-between p-5">
						<div class="flex-1">
							<h2 class="card-title">
								{item.name}
								<div class="badge {mode === 'cps' ? 'badge-secondary' : 'badge-accent'}">
									{mode === 'cps' ? 'CP' : 'Character'}
								</div>
							</h2>
						</div>

						<button
							onclick={(e) => {
								e.stopPropagation();
								toggleLike(item);
							}}
							aria-label="like"
							title="like"
							class="btn btn-circle btn-ghost"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-6 w-6 {likedTargetIds.has(item.id) ? 'fill-current' : ''} text-error"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
								/>
							</svg>
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
