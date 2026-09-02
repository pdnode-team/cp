<script lang="ts">
	import pb from '$lib/pocketbase';
	import { goto } from '$app/navigation';
	import CharacterCreateModal from '$lib/components/CharacterCreateModal.svelte';

	type ViewMode = 'cps' | 'characters';
	let mode = $state<ViewMode>('cps');
	let items = $state<any[]>([]);
	let isLoading = $state(true);
	let deleteModal = $state() as HTMLDialogElement;
	let itemToDelete = $state<null | any>(null);
	let openCreateCharDialog = $state(false);

	const handleDelete = (item: any) => {
		itemToDelete = item;
		deleteModal.showModal();
	};

	const confirmDelete = async () => {
		try {
			await pb.collection(itemToDelete.collectionName).delete(itemToDelete.id);
			deleteModal.close();
			fetchMyData(mode);
		} catch (err) {
			alert('Delete failed. Please try again.');
			console.error('Delete error:', err);
		}
	};

	async function fetchMyData(currentMode: ViewMode) {
		// 安全守卫：如果没有登录，直接跳走
		if (!pb.authStore.isValid || !pb.authStore.record) {
			goto('/login');
			return;
		}

		isLoading = true;
		const userId = pb.authStore.record.id;

		try {
			const collection = currentMode === 'cps' ? 'cps' : 'characters';
			const options = {
				// 关键点：只获取 owner 是当前用户的记录
				filter: `owner = "${userId}"`,
				sort: '-created',
				expand: currentMode === 'cps' ? 'characters' : ''
			};

			items = await pb.collection(collection).getFullList(options);
		} catch (err) {
			console.error('Fetch error:', err);
		} finally {
			isLoading = false;
		}
	}

	// 监听 mode 变化自动刷新
	$effect(() => {
		fetchMyData(mode);
	});
</script>

<div class="mx-auto max-w-6xl px-4 py-8 sm:px-6 sm:py-12 md:py-16">
	<div class="mb-8 sm:mb-12">
		<h1 class="text-3xl font-bold sm:text-4xl md:text-5xl">My Stuff</h1>
		<p class="mt-2 text-base-content/60">Manage the CPs and Characters you've created.</p>

		<div role="tablist" class="tabs-boxed mt-6 tabs w-fit">
			<button
				role="tab"
				class="tab {mode === 'cps' ? 'tab-active' : ''}"
				onclick={() => (mode = 'cps')}>My CPs</button
			>
			<button
				role="tab"
				class="tab {mode === 'characters' ? 'tab-active' : ''}"
				onclick={() => (mode = 'characters')}>My Characters</button
			>
		</div>
	</div>

	{#if isLoading}
		<div class="flex justify-center py-20">
			<span class="loading loading-lg loading-spinner text-primary"></span>
		</div>
	{:else if items.length === 0}
		<div class="card border-2 border-dashed border-base-300 bg-base-200 py-16 text-center sm:py-20">
			<div class="card-body items-center">
				<h2 class="card-title text-xl opacity-40 sm:text-2xl">
					You haven't created any {mode === 'cps' ? 'CP' : 'characters'} yet.
				</h2>
				<div class="mt-4 card-actions">
					{#if mode === 'cps'}
						<a href="/create" class="btn btn-primary">Create Your First CP</a>
					{:else}
						<button class="btn btn-primary" onclick={() => (openCreateCharDialog = true)}
							>Create Your First Character</button
						>
					{/if}
				</div>
			</div>
		</div>
	{:else}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each items as item}
				<div class="card border border-base-200 bg-base-100 shadow-xl">
					<div class="card-body p-5">
						<h2 class="card-title">{item.name}</h2>
						<div class="mt-4 card-actions flex flex-wrap justify-end gap-2">
							<button class="btn btn-soft btn-sm btn-error flex-1 sm:flex-initial" onclick={() => handleDelete(item)}>
								Delete
							</button>

							<a href="/edit/{mode}/{item.id}" class="btn btn-ghost btn-sm flex-1 sm:flex-initial"> Edit </a>

							<a href="/{mode}/{item.id}" class="btn btn-sm btn-primary flex-1 sm:flex-initial"> View </a>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<dialog bind:this={deleteModal} class="modal modal-bottom sm:modal-middle">
	<div class="modal-box border border-red-200 bg-base-100">
		<h3 class="text-lg font-bold text-base-content">Confirm Deletion</h3>
		<p class="py-4 text-base-content">
			Are you sure you want to delete this? This action cannot be undone.
		</p>
		{#if itemToDelete?.collectionName! == 'cps'}
			<p class="py-4 text-base-content">
				CP cannot be restored after deletion; the character still exists.
			</p>
		{:else}
			<p class="py-4 text-base-content">
				Once a character is deleted, it cannot be recovered, and the associated CP will also be
				deleted.
			</p>
		{/if}

		<div class="modal-action">
			<form method="dialog">
				<button class="btn btn-ghost">Cancel</button>
			</form>
			<button class="btn text-white btn-soft btn-error" onclick={confirmDelete}> Delete </button>
		</div>
	</div>

	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
<CharacterCreateModal bind:open={openCreateCharDialog} afterCreate={() => fetchMyData(mode)} />
